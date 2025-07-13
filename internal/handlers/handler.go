package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/blog/internal/config"
	"github.com/gin-gonic/gin"
)

const (
	getOffset = `select id, tittle, content, author_id, created_at
	from posts order by created_at desc limit = $1 offset = $2;`
)

func GetPostOffset(c *gin.Context) {
	limitstr := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	limit, err := strconv.Atoi(limitstr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "incorrect input",
		})
		return
	}
	db := config.NewConfigStorage().DB
	rows, err := db.Query(context.Background(), getOffset, limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	defer rows.Close()
	var posts []map[string]interface{}
	for rows.Next() {
		var (
			id        int
			title     string
			content   string
			authorID  string
			createdAt string
		)
		err = rows.Scan(&id, &title, &content, &authorID, &createdAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		posts = append(posts, map[string]interface{}{
			"id":         id,
			"title":      title,
			"content":    content,
			"author_id":  authorID,
			"created_at": createdAt,
		})
	}
	c.JSON(http.StatusOK, posts)
}
