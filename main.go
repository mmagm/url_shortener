package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mmagm/url_shortener/config"
	"github.com/mmagm/url_shortener/db"
)

var store *db.Store
var configuration config.Configuration

func AuthByApiToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token != configuration.ApiToken {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

func setupRouter() *gin.Engine {
	if configuration.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.LoadHTMLGlob("templates/*")

	r.GET("/r/:token", RedirectPage)

	api := r.Group("/api", AuthByApiToken())

	setupV1Routes(api)

	return r
}

func setupV1Routes(apiGroup *gin.RouterGroup) {
	v1 := apiGroup.Group("v1")

	v1.POST("/links", RegisterURL)
	v1.GET("/links/:token", RetrieveURL)
}

type linkPayload struct {
	URL string `json:"url" binding:"required"`
}

func RegisterURL(c *gin.Context) {
	payload := linkPayload{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}

	link, err := store.RegisterURL(payload.URL)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"url": linkShortUrl(link)})
}

func linkShortUrl(link *db.Link) string {
	return strings.Join([]string{configuration.BaseUrl, "r", link.Token}, "/")
}

func RetrieveURL(c *gin.Context) {
	token := c.Param("token")

	link, err := store.RetrieveURL(token)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": linkShortUrl(link)})
}

func RedirectPage(c *gin.Context) {
	token := c.Param("token")
	link, err := store.RetrieveURL(token)

	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{"url": link.URL})
}

func main() {
	var err error

	configuration, err = config.Load()
	if err != nil {
		panic(err)
	}
	log.Printf("Start service with configuration: %+v", configuration)

	store, err = db.NewStore()
	if err != nil {
		panic(err)
	}
	defer store.Close()

	r := setupRouter()
	r.Use(cors.Default())

	// Listen and Server in 0.0.0.0:8080
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", configuration.Server.ListenPort),
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}
	log.Println("Server exiting")
}
