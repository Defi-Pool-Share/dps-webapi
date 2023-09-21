package net

import (
	"net/http"
	"sort"

	"github.com/defi-pool-share/dps-webapi/blockchain/contractEntity"
	"github.com/defi-pool-share/dps-webapi/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func InitAPI() {
	// Echo instance
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/loan", fetchLoans)
	e.GET("/loan/available", fetchAvailableLoans)

	// Start server
	e.Logger.Fatal(e.Start(":1333"))
}

func fetchLoans(c echo.Context) error {
	loans, err := storage.FetchAllLoans()
	if err != nil {
		log.Errorf("%v", err)
		return c.String(http.StatusInternalServerError, "Service unavailable at the moment")
	}

	// Sort loans by CreationTime in descending order
	sort.Slice(loans, func(i, j int) bool {
		return loans[i].CreationTime > loans[j].CreationTime
	})
	return c.JSON(http.StatusOK, loans)
}

func fetchAvailableLoans(c echo.Context) error {
	loans, err := storage.FetchAllLoans()
	if err != nil {
		log.Errorf("%v", err)
		return c.String(http.StatusInternalServerError, "Service unavailable at the moment")
	}

	// Filter only IsActive loans
	activeLoans := make([]*contractEntity.Loan, 0)
	for _, loan := range loans {
		if loan.IsActive {
			activeLoans = append(activeLoans, loan)
		}
	}
	loans = activeLoans

	// Sort loans by CreationTime in descending order
	sort.Slice(loans, func(i, j int) bool {
		return loans[i].CreationTime > loans[j].CreationTime
	})
	return c.JSON(http.StatusOK, loans)
}
