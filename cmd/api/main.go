package main

import (
	"database/sql"
	"github.com/dnurmeden/flow-engine/internal/api"
	"github.com/dnurmeden/flow-engine/internal/repo"
	"github.com/dnurmeden/flow-engine/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load(".env")
	dsn := os.Getenv("DATABASE_URL")
	if !strings.Contains(dsn, "sslmode=") {
		if strings.Contains(dsn, "?") {
			dsn += "&sslmode=disable"
		} else {
			dsn += "?sslmode=disable"
		}
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("db ping failed: %v (dsn=%q)", err, dsn)
	}

	defer db.Close()
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	defRepo := repo.NewDefinitionRepo(db)
	instRepo := repo.NewInstanceRepo(db)
	taskRepo := repo.NewTaskRepo(db)
	procService := service.NewProcessService(defRepo, instRepo, taskRepo)
	handler := api.NewHandler(procService)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	wf := r.Group("/wf")
	{
		wf.POST("/instances", handler.StartProcess)   // POST /wf/instances
		wf.GET("/instances/:id", handler.GetInstance) // GET  /wf/instances/:id

		wf.POST("/tasks/:id/claim", handler.ClaimTask)
		wf.POST("/tasks/:id/complete", handler.CompleteTask)
	}

	addr := ":" + getenvDefault("APP_PORT", "8080")
	log.Println("API on", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func getenvDefault(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
