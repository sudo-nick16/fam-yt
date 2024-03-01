package server

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sudo-nick16/fam-yt/internal/config"
	"github.com/sudo-nick16/fam-yt/internal/repository"
	"github.com/sudo-nick16/fam-yt/internal/server/handlers"
	"github.com/sudo-nick16/fam-yt/internal/ytapi"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := err.Error()
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message.(string)
	}
	// log.Printf("[ERROR] %v", msg)
	c.JSON(code, map[string]interface{}{
		"error": msg,
	})
}

func Start() {
	config := config.GetConfig()
	opts := options.Client().ApplyURI(config.MongoUri)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Panicln("[ERROR] Could not connect to the database:", err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Panicln("[ERROR] Could not ping the database:", err)
	}
	log.Println("[INFO] Connected to the database")

	db := client.Database(config.DbName)
	vidRepo := repository.NewVideoRepository(db, "videos")
	err = vidRepo.CreateTextIndex()
	if err != nil {
		log.Println("[INFO] Could not ensure text index for videos")
	} else {
		log.Println("[INFO] Ensured text index for videos")
	}
	searchRepo := repository.NewSearchRepository(db, "search-queries")
	err = searchRepo.CreateSimpleIndex()
	if err != nil {
		log.Println("[INFO] Could not ensure simple index for search queries:", err)
	} else {
		log.Println("[INFO] Ensured simple index for search queries")
	}

	ytApi, err := ytapi.NewYtApi(config.YtApiKeys, config.MaxResults)
	if err != nil {
		log.Panicln("[ERROR] Could not create ytapi:", err)
	}

	e := echo.New()
	e.Use(middleware.CORS())
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.GET("/api/videos", handlers.GetVideos(vidRepo))

	e.POST("/api/queries", handlers.CreateQuery(searchRepo, vidRepo, ytApi))

	e.GET("/api/queries", handlers.GetQueries(searchRepo))

	e.GET("/api/info", handlers.GetInfo(config))

	if err := e.StartTLS(config.Port, "cert.pem", "key.pem"); err != http.ErrServerClosed {
		log.Panicln("[ERROR] Could not start the server:", err)
	}
}
