package repository

import (
	"context"
	"log"
	"time"

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

func (yt *VideoRepository) Find(query string,
	limit int64, page int64, order string) ([]types.Video, error) {
	regex := `(.*)` + query + `(.*)`
	filter := bson.M{
		"$text": bson.M{
			"$search": regex,
		},
	}
	sortOrder := -1
	if so, ok := orderMap[order]; ok {
		sortOrder = so
	}
	sort := bson.M{"publishedAt": sortOrder}
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

func (yt *VideoRepository) CreateTextIndex() error {
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

// COMMENT: this is atlas deployment specific, so the index needs to be created
// either using atlas cli, or the atlas dashboard
func (yt *VideoRepository) FindInAtlasSearch(query string,
	limit int64, page int64, order string) ([]types.Video, error) {
	regex := `(.*)` + query + `(.*)`
	filterStage := bson.D{
		{"$search", bson.D{
			{"index", "text"},
			{"regex", bson.D{
				{"path", "searchQuery"},
				{"query", regex},
				{"allowAnalyzedField", true},
			},
			}},
		},
	}
	sortStage := bson.D{{"$sort", bson.D{{"publishedAt", orderMap[order]}}}}
	skipStage := bson.D{{"$skip", page * limit}}
	limitStage := bson.D{{"$limit", limit}}
	projectStage := bson.D{{"$project", bson.D{
		{"_id", 1},
		{"videoId", 1},
		{"title", 1},
		{"description", 1},
		{"publishedAt", 1},
		{"thumbnail", 1},
		{"searchQuery", 1},
		{"searchId", 1},
	}}}
	opts := options.Aggregate().SetMaxTime(5 * time.Second)
	cursor, err := yt.coll.Aggregate(context.TODO(), mongo.Pipeline{
		filterStage, sortStage, skipStage, limitStage, projectStage}, opts)
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
