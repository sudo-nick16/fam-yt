package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sudo-nick16/fam-yt/internal/config"
	"github.com/sudo-nick16/fam-yt/internal/repository"
	"github.com/sudo-nick16/fam-yt/internal/server/handlers"
	"github.com/sudo-nick16/fam-yt/internal/types"
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
	db := client.Database(config.DbName)
	vidRepo := repository.NewVideoRepository(db, "videos")
	vidRepo.CreateIndex()
	searchRepo := repository.NewSearchRepository(db, "search-queries")

	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.GET("/api/search", handlers.GetVideo(vidRepo))

	e.POST("/api/queries", func(c echo.Context) error {
		sq := types.SearchQuery{}
		json.NewDecoder(c.Request().Body).Decode(&sq)
		err := searchRepo.Create(sq.Query)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, sq)
	})

	e.Start(config.Port)
}
