package main

import (
	"log"
	"os"

	"github.com/daiki-trnsk/test-debatemap/api/controllers"
	"github.com/daiki-trnsk/test-debatemap/api/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load("/app/.env")
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

	
    database.Client = database.DBSet()
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
    }))

	auth := e.Group("")
	auth.GET("/debatemaps_list", controllers.GetDebateMaps)
	auth.GET("/debatemap/:id", controllers.GetDebateMap)
	auth.PUT("/update/:id", controllers.UpdateDebateMap)
	auth.POST("/adddebatemap", controllers.AddDebateMap)
	auth.GET("/search", controllers.SearchDebateMapByQuery)

	log.Fatal(e.Start(":" + port))
}
