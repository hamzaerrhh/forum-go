package database

import (
	"fmt"
	"time"
)

type Comment struct {
	Created_at time.Time
	Text       string
}

func GetCommentsByPost(postId int) ([]Comment, error) {
	var comments []Comment
	rows, err := Database.Query(
		"SELECT c.created_at, c.text FROM Comments c INNER JOIN Posts p ON c.post_id = p.id WHERE p.id = ?",
		postId,
	)
	defer rows.Close() // release database resources
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.Created_at, &c.Text); err != nil {
			return nil, fmt.Errorf("getCommentsByPost error: %v", err)
		}
		comments = append(comments, c)
	}
	// Important: Check for any errors that occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getCommentsByPost error: %v", err)
	}
	return comments, err
}
