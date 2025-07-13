package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func main() {
	ctx := context.Background()
	var (
		err error
	)
	// db, err = pgxpool.New(ctx, cfg.DATA_BASE_URI)
	db, err = pgxpool.New(ctx, "postgres://alex01:pwd1234@localhost:5432/blog")
	if err != nil {
		return
	}
	r := gin.Default()
	r.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"hello": "hello world",
		})
	})
	r.POST("/create", createPost)
	r.GET("/get/:id", getPostbyId)
	r.Run(":8080")
}

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  int       `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

func createPost(c *gin.Context) {
	var p Post
	if err := c.ShouldBindBodyWithJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input incorrect"})
		return
	}
	err := db.QueryRow(context.Background(), "insert into posts (tittle, content, author_id) values($1,$2,$3)returning id, created_at",
		p.Title, p.Content, p.AuthorID).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}
func getPostbyId(c *gin.Context) {
	id := c.Param("id")
	var p Post
	err := db.QueryRow(context.Background(), "select * from posts where id=$1", id).Scan(
		&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, p)
}
