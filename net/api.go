package net

import (
	"net/http"

	"github.com/defi-pool-share/dps-webapi/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func InitAPI() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/loan", fetchLoans)

	// Start server
	e.Logger.Fatal(e.Start(":1333"))
}

func fetchLoans(c echo.Context) error {
	loans, err := storage.FetchAllLoans()
	if err != nil {
		log.Errorf("%v", err)
		return c.String(http.StatusInternalServerError, "Service unavailable at the moment")
	}
	/**
	loanJSON, err := json.Marshal(loans)
	if err != nil {
		log.Errorf("%v", err)
		return c.String(http.StatusInternalServerError, "Service unavailable at the moment")
	}
	*/
	return c.JSON(http.StatusOK, loans)
}
