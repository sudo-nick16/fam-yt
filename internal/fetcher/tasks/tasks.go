package tasks

import (
	"log"
	"time"

	"github.com/sudo-nick16/fam-yt/internal/repository"
	"github.com/sudo-nick16/fam-yt/internal/types"
	"github.com/sudo-nick16/fam-yt/internal/ytapi"
)

type FetchQueryTask struct {
	retryCount int
	maxRetries int
	query      *types.SearchQuery
	ytApi      *ytapi.YtApi
	searchRepo *repository.SearchRepository
	videoRepo  *repository.VideoRepository
}

func NewFetchQueryTask(ytApi *ytapi.YtApi,
	searchRepo *repository.SearchRepository,
	videoRepo *repository.VideoRepository,
	query *types.SearchQuery,
	maxRetries int) *FetchQueryTask {
	return &FetchQueryTask{
		ytApi:      ytApi,
		searchRepo: searchRepo,
		videoRepo:  videoRepo,
		query:      query,
		retryCount: 0,
		maxRetries: maxRetries,
	}
}

func (ft *FetchQueryTask) Execute() error {
	ft.retryCount++
	videos, err := ft.ytApi.GetLatestVideos(ft.query.Query,
		ft.query.LatestPublishedAt.Time().Format(time.RFC3339))
	if err != nil {
		return err
	}
	if len(videos) == 0 {
		log.Println("No videos found for query after")
		return nil
	}

	latest := videos[0].PublishedAt

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
	if ft.retryCount < ft.maxRetries {
		ft.Execute()
	}
}

func (ft *FetchQueryTask) Success() {
	// TODO: handle success :)
}
