package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Video struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	VideoId     string             `bson:"videoId" json:"videoId"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	PublishedAt primitive.DateTime `bson:"publishedAt,timestamp" json:"publishedAt"`
	Thumbnail   string             `bson:"thumbnail" json:"thumbnail"`
	SearchQuery string             `bson:"searchQuery,omitempty"`
	SearchId    primitive.ObjectID `bson:"searchId,omitempty"`
}

type SearchQuery struct {
	Id    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Query string             `bson:"query" json:"query"`
	// LatestPublishedAt is the latest published time of the videos fetched for this query
	LatestPublishedAt primitive.DateTime `bson:"latestPublishedAt,timestamp" json:"latestPublishedAt"`
}
