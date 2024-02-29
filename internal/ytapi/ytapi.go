package ytapi

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/sudo-nick16/fam-yt/internal/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

type YtApi struct {
	currKeyIndex int
	apiKeys      []string
	service      *youtube.Service
	maxResults   uint64
}

func NewYtApi(apiKeys []string, maxResults uint64) (*YtApi, error) {
	if len(apiKeys) == 0 {
		return nil, errors.New("[ERROR] No API keys provided")
	}
	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKeys[0]},
	}
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("[ERROR] Could not create a new YouTube client: %v", err)
	}
	return &YtApi{
		// TODO: Is it okay to start with the first API key?
		currKeyIndex: 0,
		apiKeys:      apiKeys,
		service:      service,
		maxResults:   maxResults,
	}, nil
}

func (yt *YtApi) RotateApiKey() {
	yt.currKeyIndex++
	if yt.currKeyIndex >= len(yt.apiKeys) {
		yt.currKeyIndex %= len(yt.apiKeys)
	}
	yt.UpdateYtClient()
}

func (yt *YtApi) UpdateYtClient() {
	client := &http.Client{
		Transport: &transport.APIKey{Key: yt.apiKeys[yt.currKeyIndex]},
	}
	service, err := youtube.New(client)
	if err != nil {
		log.Printf("[ERROR] Could not update the YouTube client: %v", err)
		return
	}
	yt.service = service
}

func (yt *YtApi) GetLatestVideos(query string, publishedAfter string) ([]types.Video, error) {
	sl := yt.service.Search.List([]string{"id", "snippet"}).
		Q(query).
		Order("date").
		PublishedAfter(publishedAfter).
		MaxResults(int64(yt.maxResults))
	resp, err := sl.Do()
	// TODO: Check if error is caused due to API key limit,
	// 	if yes then rotate the API key
	if err != nil {
		return nil, err
	}
	videos := []types.Video{}
	for _, item := range resp.Items {
		switch item.Id.Kind {
		case "youtube#video":
			{
				pubAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
				if err != nil {
					log.Printf("[ERROR] Could not parse publishedAt: %v", err)
					continue
				}
				videos = append(videos, types.Video{
					VideoId:     item.Id.VideoId,
					Title:       item.Snippet.Title,
					Description: item.Snippet.Description,
					PublishedAt: primitive.NewDateTimeFromTime(pubAt),
					Thumbnail:   item.Snippet.Thumbnails.Default.Url,
				})
			}
		}
	}
	return videos, nil
}
