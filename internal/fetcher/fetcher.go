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

	pool := workerpool.NewWorkerPool(10)
	pool.Start()

	ticker := time.NewTicker(time.Duration(config.PollInterval) * time.Second)

	ytApi, err := ytapi.NewYtApi(config.YtApiKeys, config.MaxResults)
	if err != nil {
		log.Panicln("[ERROR] Could not create YouTube API client:", err)
	}

	for {
		select {
		case <-ticker.C:
			log.Println("[INFO] tick")
			queries, err := sqRepo.FindAll()
			if err != nil {
				log.Println("[ERROR] Could not fetch queries:", err)
				continue
			}
			log.Println("[INFO] fetched queries:", queries)
			for _, query := range queries {
				log.Println("[INFO] current query:", query)
				task := tasks.NewFetchQueryTask(ytApi, sqRepo, vidRepo, query)
				pool.AddTask(task)
			}
		}
	}
}
