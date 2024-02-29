package repository

import (
	"context"
	"log"

	"github.com/sudo-nick16/fam-yt/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderMap = map[string]int{
	"asc":  1,
	"desc": -1,
}

type VideoRepository struct {
	coll *mongo.Collection
}

func NewVideoRepository(db *mongo.Database, collectionName string) *VideoRepository {
	coll := db.Collection(collectionName)
	return &VideoRepository{
		coll: coll,
	}
}

func (yt *VideoRepository) InsertMany(videos []types.Video) error {
	docs := make([]interface{}, len(videos))
	for i := range docs {
		docs[i] = videos[i]
	}
	res, err := yt.coll.InsertMany(context.TODO(), docs)
	if err != nil {
		return err
	}
	if len(res.InsertedIDs) != len(docs) {
		return mongo.ErrUnacknowledgedWrite
	}
	return nil
}

func (yt *VideoRepository) Find(query string, limit int64, page int64, order string) ([]types.Video, error) {
	filter := bson.D{{
		"$text", bson.D{{
			"$search", query,
		}},
	}}
	sort := bson.D{
		{
			"publishedAt", orderMap[order],
		},
	}
	opts := options.Find().SetLimit(limit).SetSkip(page * limit).SetSort(sort)
	cursor, err := yt.coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	videos := []types.Video{}
	err = cursor.All(context.TODO(), &videos)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (yt *VideoRepository) CreateIndex() error {
	model := mongo.IndexModel{
		Keys: bson.D{{
			"searchQuery", "text",
		}},
	}
	name, err := yt.coll.Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created index: %s", name)
	return nil
}
