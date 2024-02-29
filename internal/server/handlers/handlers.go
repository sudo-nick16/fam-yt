package handlers

import (
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sudo-nick16/fam-yt/internal/repository"
)

func GetVideo(vidRepo *repository.VideoRepository) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		query := ctx.QueryParam("q")
		if query == "" {
			return echo.NewHTTPError(400, "Query parameter 'q' is required")
		}
		p := ctx.QueryParam("page")
		if p == "" {
			p = "0"
		}
		page, err := strconv.Atoi(p)
		if err != nil {
			return echo.NewHTTPError(400, "Page must be a valid integer")
		}
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
		videos, err := vidRepo.Find(query, int64(limit), int64(page))
		log.Println(videos, err)
		if err != nil {
			return echo.NewHTTPError(500, "Could not fetch videos.")
		}
		return ctx.JSON(200, videos)
	}
}
