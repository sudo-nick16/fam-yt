package handlers

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sudo-nick16/fam-yt/internal/config"
	"github.com/sudo-nick16/fam-yt/internal/repository"
	"github.com/sudo-nick16/fam-yt/internal/types"
	"github.com/sudo-nick16/fam-yt/internal/ytapi"
)

func GetVideos(vidRepo *repository.VideoRepository) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		query := ctx.QueryParam("query")
		if query == "" {
			return echo.NewHTTPError(400, "Query parameter 'query' is required")
		}
		p := ctx.QueryParam("pageno")
		if p == "" {
			p = "1"
		}
		page, err := strconv.Atoi(p)
		if err != nil {
			return echo.NewHTTPError(400, "Page must be a valid integer")
		}
		// 0-based page number
		page -= 1
		if page < 0 {
			return echo.NewHTTPError(400, "Page must be a positive integer")
		}
		l := ctx.QueryParam("limit")
		if l == "" {
			l = "10"
		}
		limit, err := strconv.Atoi(l)
		if err != nil {
			return echo.NewHTTPError(400, "Limit must be a valid integer")
		}
		if limit < 1 {
			return echo.NewHTTPError(400, "Limit must be a positive integer")
		}
		order := ctx.QueryParam("order")
		if order == "" {
			order = "desc"
		}
		if order != "asc" && order != "desc" {
			return echo.NewHTTPError(400, "Order must be 'asc' or 'desc'")
		}
		videos, err := vidRepo.Find(query, int64(limit), int64(page), order)
		if err != nil {
			log.Println("[ERROR] Could not fetch videos:", err)
			return echo.NewHTTPError(500, "Could not fetch videos.")
		}
		return ctx.JSON(200, map[string]interface{}{
			"videos": videos,
		})
	}
}

func CreateQuery(searchRepo *repository.SearchRepository,
	vidRepo *repository.VideoRepository, ytApi *ytapi.YtApi) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sq := &types.SearchQuery{}
		err := json.NewDecoder(ctx.Request().Body).Decode(sq)
		if err != nil {
			return echo.NewHTTPError(500, "Could not parse body.")
		}
		if sq.Query == "" {
			return echo.NewHTTPError(400, "Received empty query.")
		}
		sq, err = searchRepo.Create(sq.Query)
		if err != nil {
			return echo.NewHTTPError(500, "Could not create query.")
		}
		// COMMENT: Initially, we would wait for the fetcher to do this, but
		// for better ux we'll fetch few results immediately
		// task := tasks.NewFetchQueryTask(ytApi, searchRepo, vidRepo, sq)
		// COMMENT: We don't need to wait for the task to complete or fail,
		// because the fetcher will do the same task in the background in the
		// next interval
		// go task.Execute()
		return ctx.JSON(200, map[string]interface{}{
			"msg":   "Query created successfully.",
			"query": *sq,
		})
	}
}

func GetQueries(searchRepo *repository.SearchRepository) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		queries, err := searchRepo.FindAll()
		if err != nil {
			return echo.NewHTTPError(500, "Could not fetch queries.")
		}
		return ctx.JSON(200, map[string]interface{}{
			"queries": queries,
		})
	}
}

func GetInfo(config *config.Config) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.JSON(200, map[string]interface{}{
			"pollInterval":    config.PollInterval,
			"ytApiMaxResults": config.MaxResults,
		})
	}
}
