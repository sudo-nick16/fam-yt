package tasks

import (
	"log"

	"github.com/sudo-nick16/fam-yt/internal/repository"
	"github.com/sudo-nick16/fam-yt/internal/types"
	"github.com/sudo-nick16/fam-yt/internal/ytapi"
)

type FetchQueryTask struct {
	query      *types.SearchQuery
	ytApi      *ytapi.YtApi
	searchRepo *repository.SearchRepository
	videoRepo  *repository.VideoRepository
}

func NewFetchQueryTask(ytApi *ytapi.YtApi,
	searchRepo *repository.SearchRepository,
	videoRepo *repository.VideoRepository,
	query *types.SearchQuery) *FetchQueryTask {
	return &FetchQueryTask{
		ytApi:      ytApi,
		searchRepo: searchRepo,
		videoRepo:  videoRepo,
		query:      query,
	}
}

func (ft *FetchQueryTask) Execute() error {
	videos, err := ft.ytApi.GetLatestVideos(ft.query)
	log.Printf("Got %d videos for query: %s , err: %v\n", len(videos), ft.query.Query, err)
	if err != nil {
		return err
	}
	if len(videos) == 0 {
		log.Println("No videos found for query after")
		return nil
	}

	latest := videos[0].PublishedAt
	log.Printf("Latest video published at: %+v\n", latest)

	// COMMENT: This is not needed because the API already returns the
	// latest videos after the latestPublishedAt time
	// if latest.Time().Equal(ft.query.LatestPublishedAt.Time()) ||
	// 	latest.Time().Before(ft.query.LatestPublishedAt.Time()) {
	// 	return nil
	// }

	for i := range videos {
		videos[i].SearchQuery = ft.query.Query
		videos[i].SearchId = ft.query.Id
	}
	err = ft.videoRepo.InsertMany(videos)
	if err != nil {
		return err
	}
	return ft.searchRepo.UpdateLatest(ft.query.Id, latest)
}

func (ft *FetchQueryTask) Failed() {
	// TODO: handle failure :)
	// log.Println("Failed to fetch query.")
}

func (ft *FetchQueryTask) Success() {
	// TODO: handle success :)
	// log.Println("Successfully fetched query.")
}
