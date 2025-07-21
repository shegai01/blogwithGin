package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

type Config struct {
	DataBase_URI string `json:"database_uri"`
	BindAddr     string `json:"bind_ddr"`
}

const (
	getPostPagineted = `
	select id, title, content, author_id, created_at from posts where (created_at < $1) or
	(created_at =$1 and id<$2) order by created_at desc, id desc limit $3;
	`
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var (
		err error
	)
	cfg := Config{
		DataBase_URI: os.Getenv("APP_DATABASE_URI"),
		BindAddr:     os.Getenv("APP_BIND_ADDR"),
	}
	db, err = pgxpool.New(ctx, cfg.DataBase_URI)
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
	r.GET("/posts", GetPostPagineted)
	r.Run(cfg.BindAddr)
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
	err := db.QueryRow(context.Background(), "insert into posts (title, content, author_id) values($1,$2,$3)returning id, created_at",
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
func GetPostPagineted(ctx *gin.Context) {
	const limit = 20
	createdAtParam := ctx.Query("create_at")
	idParams := ctx.Query("id")
	var cursorTime time.Time = time.Now()
	var cursorId int64 = 1 << 62
	if createdAtParam != "" {
		t, err := time.Parse(time.RFC3339, createdAtParam)
		if err == nil {
			cursorTime = t
		}
	}
	if idParams != "" {
		id, err := strconv.ParseInt(idParams, 10, 64)
		if err == nil {
			cursorId = id
		}
	}
	rows, err := db.Query(context.Background(), getPostPagineted, cursorTime, cursorId, limit)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()
	var posts []Post

	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.AuthorID, &p.Content, &p.CreatedAt); err != nil {
			continue
		}
		posts = append(posts, p)
	}
	var nextCursor string
	if len(nextCursor) > 0 {
		last := posts[len(posts)-1]
		nextCursor = "?create_at" + last.CreatedAt.Format(time.RFC3339) + "&id=" + strconv.FormatInt(int64(last.ID), 10)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"posts":       posts,
		"next_cursor": nextCursor,
	})
}
