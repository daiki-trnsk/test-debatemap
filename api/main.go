package main

import (
	"log"
	"os"

	"github.com/daiki-trnsk/go-debatemap/api/controllers"
	customMiddleware "github.com/daiki-trnsk/go-debatemap/api/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"http://localhost:3000"},
        AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
    }))

	e.POST("/signup", controllers.SignUp)
	e.POST("/login", controllers.Login)

	auth := e.Group("")
	auth.Use(customMiddleware.Authentication)
	auth.GET("/debatemaps_list", controllers.GetDebateMaps)
	auth.GET("/debatemap/:id", controllers.GetDebateMap)
	auth.PUT("/update/:id", controllers.UpdateDebateMap)
	auth.POST("/adddebatemap", controllers.AddDebateMap)
	auth.GET("/search", controllers.SearchDebateMapByQuery)

	log.Fatal(e.Start(":" + port))
}
