package api

import (
	"fmt"
	"strings"
	"time"

	"forum/database"
)



type Post struct {
    Id         int
    UserId     int
    Created_at time.Time
    Title      string
    Text       string
    Comments   []Comment
    Categories []string // Add this field to store category names
}

func GetPosts() ([]Post, error) {
    var posts []Post
    
    // Modified query to include user_id since we need it for categories
    rows, err := database.Database.Query(
        "SELECT id, user_id, created_at, title, text FROM posts ORDER BY created_at DESC",
    )
    if err != nil {
        return nil, fmt.Errorf("getPosts error: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var p Post
        if err := rows.Scan(&p.Id, &p.UserId, &p.Created_at, &p.Title, &p.Text); err != nil {
            return nil, fmt.Errorf("getPosts error: %v", err)
        }

        // Get comments for the post
        if p.Comments, err = GetCommentsByPost(p.Id); err != nil {
            return nil, err
        }

        // Get categories for the post
        categories, err := GetCategoriesByPost(p.Id)
        if err != nil {
            return nil, err
        }
        p.Categories = categories

        posts = append(posts, p)
    }

    // Check for any errors that occurred during iteration
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("getPosts error: %v", err)
    }

    return posts, nil
}

// Helper function to get categories for a specific post
func GetCategoriesByPost(postId int) ([]string, error) {
    var categories []string
    
    rows, err := database.Database.Query(`
        SELECT c.name 
        FROM category c
        JOIN post_category pc ON c.id = pc.category_id
        WHERE pc.post_id = ?
        ORDER BY c.name
    `, postId)
    if err != nil {
        return nil, fmt.Errorf("GetCategoriesByPost error: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var category string
        if err := rows.Scan(&category); err != nil {
            return nil, fmt.Errorf("GetCategoriesByPost error: %v", err)
        }
        categories = append(categories, category)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("GetCategoriesByPost error: %v", err)
    }

    return categories, nil
}

// Alternative optimized version that gets all posts with their categories in a single query
func GetPostsOptimized() ([]Post, error) {
    var posts []Post
    
    // This query gets all posts with their categories aggregated as a JSON array or comma-separated string
    // Since SQLite doesn't have native JSON functions in older versions, we'll use GROUP_CONCAT
    rows, err := database.Database.Query(`
        SELECT 
            p.id, 
            p.user_id, 
            p.created_at, 
            p.title, 
            p.text,
            COALESCE(GROUP_CONCAT(c.name, ','), '') as categories
        FROM posts p
        LEFT JOIN post_category pc ON p.id = pc.post_id
        LEFT JOIN category c ON pc.category_id = c.id
        GROUP BY p.id
        ORDER BY p.created_at DESC
    `)
    if err != nil {
        return nil, fmt.Errorf("GetPostsOptimized error: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var p Post
        var categoriesStr string
        
        if err := rows.Scan(&p.Id, &p.UserId, &p.Created_at, &p.Title, &p.Text, &categoriesStr); err != nil {
            return nil, fmt.Errorf("GetPostsOptimized error: %v", err)
        }

        // Parse comma-separated categories
        if categoriesStr != "" {
            p.Categories = strings.Split(categoriesStr, ",")
        } else {
            p.Categories = []string{}
        }

        // Get comments for the post
        if p.Comments, err = GetCommentsByPost(p.Id); err != nil {
            return nil, err
        }

        posts = append(posts, p)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("GetPostsOptimized error: %v", err)
    }

    return posts, nil
}