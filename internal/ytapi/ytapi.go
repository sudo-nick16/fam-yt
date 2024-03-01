package ytapi

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sudo-nick16/fam-yt/internal/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

type YtApi struct {
	currKeyIndex       int
	apiKeys            []string
	service            *youtube.Service
	maxResults         int
	quotaExceededCount int
}

func NewYtApi(apiKeys []string, maxResults int) (*YtApi, error) {
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
	log.Println("[INFO] API key rotated.")
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

func (yt *YtApi) GetLatestVideos(query *types.SearchQuery) ([]types.Video, error) {
	sl := yt.service.Search.List([]string{"id", "snippet"}).
		Q(query.Query).
		Order("date").
		Type("video").
		PublishedAfter(query.LatestPublishedAt.Time().Format(time.RFC3339)).
		MaxResults(int64(yt.maxResults))

	resp, err := sl.Do()
	if err != nil {
		if strings.Contains(err.Error(), "quotaExceeded") {
			log.Printf("[ERROR] Quota exceeded: %v", err)
			yt.quotaExceededCount++
			if yt.quotaExceededCount >= len(yt.apiKeys) {
				yt.quotaExceededCount = 0
				return nil, errors.New("All API keys have exceeded their quota")
			}
			yt.RotateApiKey()
			return yt.GetLatestVideos(query)
		}
		return nil, err
	}
	videos := make([]types.Video, 0, yt.maxResults)
	for _, item := range resp.Items {
		switch item.Id.Kind {
		case "youtube#video":
			{
				pubAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
				if err != nil {
					log.Printf("[ERROR] Could not parse publishedAt: %v", err)
					continue
				}
				if pubAt.Before(query.LatestPublishedAt.Time()) || pubAt.Equal(query.LatestPublishedAt.Time()) {
					continue
				}
				videos = append(videos, types.Video{
					VideoId:     item.Id.VideoId,
					Title:       item.Snippet.Title,
					Description: item.Snippet.Description,
					PublishedAt: primitive.NewDateTimeFromTime(pubAt),
					Thumbnail:   item.Snippet.Thumbnails.Default.Url,
					SearchQuery: query.Query,
					SearchId:    query.Id,
				})
			}
		}
	}
	return videos, nil
}
