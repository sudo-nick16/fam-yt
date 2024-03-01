package repository

import (
	"context"
	"errors"
	"time"

	"github.com/sudo-nick16/fam-yt/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// COMMENT: Redis could be used to store search queries
type SearchRepository struct {
	coll *mongo.Collection
}

func NewSearchRepository(db *mongo.Database, collectionName string) *SearchRepository {
	coll := db.Collection(collectionName)
	return &SearchRepository{
		coll: coll,
	}
}

func (p *SearchRepository) FindAll() ([]types.SearchQuery, error) {
	cursor, err := p.coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	queries := []types.SearchQuery{}
	err = cursor.All(context.TODO(), &queries)
	if err != nil {
		return nil, err
	}
	return queries, nil
}

func (p *SearchRepository) FindByQuery(query string) (*types.SearchQuery, error) {
	filter := bson.D{{
		"query", query,
	}}
	result := &types.SearchQuery{}
	err := p.coll.FindOne(context.TODO(), filter).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *SearchRepository) FindById(id primitive.ObjectID) (*types.SearchQuery, error) {
	filter := bson.D{{
		"_id", id,
	}}
	result := &types.SearchQuery{}
	err := p.coll.FindOne(context.TODO(), filter).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *SearchRepository) UpdateLatest(queryId primitive.ObjectID,
	time primitive.DateTime) error {
	filter := bson.D{{
		"_id", queryId,
	}}
	update := bson.D{{
		"$set", bson.D{{
			"latestPublishedAt", time,
		}},
	}}
	_, err := p.coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (p *SearchRepository) Create(query string) (*types.SearchQuery, error) {
	// COMMENT: This is not needed because of the unique index
	// _, err := p.FindByQuery(query)
	// if err == nil {
	// 	return nil, errors.New("Query already exists")
	// }
	sq := types.SearchQuery{
		Query:             query,
		LatestPublishedAt: primitive.NewDateTimeFromTime(time.Now().AddDate(-10, 0, 0)),
	}
	res, err := p.coll.InsertOne(context.TODO(), sq)
	if err != nil {
		return nil, err
	}
	if res.InsertedID == nil {
		return nil, errors.New("Could not insert query")
	}
	resId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("Could not convert InsertedID to ObjectID")
	}
	sq.Id = resId
	return &sq, nil
}

// Creating index to avoid duplicate queries
func (p *SearchRepository) CreateSimpleIndex() error {
	model := mongo.IndexModel{
		Keys: bson.D{{
			"query", 1,
		}},
	}
	_, err := p.coll.Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		return err
	}
	return nil
}
