package fetcher

import (
	"context"
	"log"
	"time"

	"github.com/sudo-nick16/fam-yt/internal/config"
	"github.com/sudo-nick16/fam-yt/internal/fetcher/tasks"
	"github.com/sudo-nick16/fam-yt/internal/repository"
	"github.com/sudo-nick16/fam-yt/internal/workerpool"
	"github.com/sudo-nick16/fam-yt/internal/ytapi"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Description: This service fetches videos from the YouTube API for pre-defined queries and caches the results.

func StartFetching() {
	config := config.GetConfig()
	opts := options.Client().ApplyURI(config.MongoUri)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Panicln("[ERROR] Could not connect to the database:", err)
	}
	defer func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Panicln("[ERROR] Could not disconnect:", err)
		}
	}()
	db := client.Database(config.DbName)
	sqRepo := repository.NewSearchRepository(db, "search-queries")
	vidRepo := repository.NewVideoRepository(db, "videos")
	log.Println("[INFO] Fetching videos for pre-defined queries", *sqRepo)
	log.Println("[INFO] Fetching videos for pre-defined queries", *vidRepo)

	// ytApi, err := ytapi.New([]string{config.YtApiKey}, 5)
	// if err != nil {
	// 	log.Panicln("[ERROR] Could not disconnect:", err)
	// }

	pool := workerpool.NewWorkerPool(10)
	pool.Start()

	ticker := time.NewTicker(10 * time.Second)

	ytApi, err := ytapi.NewYtApi([]string{config.YtApiKey}, 2)
	if err != nil {
		log.Panicln("[ERROR] Could not create YouTube API client:", err)
	}

	for {
		select {
		case <-ticker.C:
			queries, err := sqRepo.FindAll()
			if err != nil {
				log.Println("[ERROR] Could not fetch queries:", err)
				continue
			}
			log.Println("[INFO] Fetched queries:", queries)
			for _, query := range queries {
				task := tasks.NewFetchQueryTask(ytApi, sqRepo, vidRepo, &query, 1)
				pool.AddTask(task)
			}
		}
	}
}
