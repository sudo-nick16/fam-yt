package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sudo-nick16/fam-yt/internal/repository"
	"github.com/sudo-nick16/fam-yt/internal/types"
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
			return echo.NewHTTPError(500, "Could not fetch videos.")
		}
		return ctx.JSON(200, map[string]interface{}{
			"videos": videos,
		})
	}
}

func CreateQuery(searchRepo *repository.SearchRepository) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sq := types.SearchQuery{}
		err := json.NewDecoder(ctx.Request().Body).Decode(&sq)
		if err != nil {
			return echo.NewHTTPError(500, "Could not parse body.")
		}
		if sq.Query == "" {
			return echo.NewHTTPError(400, "Received empty query.")
		}
		err = searchRepo.Create(sq.Query)
		if err != nil {
			return echo.NewHTTPError(500, "Could not create query.")
		}
		return ctx.JSON(200, map[string]string{
			"msg": "Query created successfully.",
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
