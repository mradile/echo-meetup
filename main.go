package main

import (
	"github.com/labstack/echo/v4"

	"log"
	"net/http"
	"time"
)

type Report struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type reportsStorage map[string]Report

var reports = make(reportsStorage)

func main() {
	e := echo.New()

	httpServer := &http.Server{
		Addr:              ":3131",
		ReadTimeout:       60 * time.Second,  // time to read request from client
		ReadHeaderTimeout: 10 * time.Second,  // time to read header
		WriteTimeout:      20 * time.Second,  // time to write to the client
		IdleTimeout:       120 * time.Second, // time between keep-alives requests before connection is closed
	}

	reportRoutes := e.Group("/api/v1")
	reportRoutes.POST("/report", Add)
	reportRoutes.GET("/report/:id", Get)
	reportRoutes.GET("/report", All)
	reportRoutes.PUT("/report/:id", Update)
	reportRoutes.DELETE("/report/:id", Delete)

	e.Static("/", "public")

	//e.Use(middleware.Logger())
	e.Use(LoggerMiddleware())

	e.Logger.Fatal(e.StartServer(httpServer))
}

func Add(c echo.Context) error {
	var r Report
	if err := c.Bind(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "request body invalid")
	}

	reports[r.ID] = r

	return c.JSONPretty(http.StatusOK, r, "  ")
}

func Get(c echo.Context) error {
	id := c.Param("id")

	r, ok := reports[id]
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSONPretty(http.StatusOK, r, "  ")
}

func All(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, reports, "  ")
}

func Update(c echo.Context) error {
	var r Report
	if err := c.Bind(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "request body invalid")
	}

	_, ok := reports[r.ID]
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	reports[r.ID] = r

	return c.JSONPretty(http.StatusOK, r, "  ")
}

func Delete(c echo.Context) error {
	id := c.Param("id")

	r, ok := reports[id]
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	delete(reports, r.ID)

	return c.NoContent(http.StatusNoContent)
}

func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)

			path := c.Request().URL.Path
			code := c.Response().Status
			duration := time.Now().Sub(start)

			log.Printf("request %s %d %s", path, code, duration)

			return err
		}
	}
}
