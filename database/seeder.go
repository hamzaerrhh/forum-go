package database

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var (
	firstNames = []string{"Alice", "Bob", "Charlie", "Diana", "Eve", "Frank", "Grace", "Henry", "Iris", "Jack", "Kate", "Liam", "Mia", "Noah", "Olivia"}
	lastNames  = []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez"}

	postTitles = []string{
		"Tips for staying healthy during winter",
		"Best places to visit in Europe",
		"How to learn programming in 2024",
		"My journey to financial freedom",
		"Amazing hiking trails near the city",
		"Beginner's guide to meditation",
		"Why remote work is the future",
		"Top 10 productivity hacks",
		"Travel on a budget: The ultimate guide",
		"How I improved my fitness in 3 months",
		"The art of minimalist living",
		"Starting your own business: First steps",
		"Sustainable living practices",
		"Best books I read this year",
		"How to ace your job interview",
	}

	postTexts = []string{
		"This is an amazing topic that I've been thinking about lately. There are so many great insights and perspectives to explore here.",
		"I've been doing research on this subject for quite some time now, and I've discovered some really interesting findings that I'd love to share with everyone.",
		"Has anyone else experienced something similar? I'm really curious to hear what others have to say about this.",
		"Just finished reading up on this and wanted to share my thoughts. The information out there is quite helpful.",
		"This has been a game-changer for me! I can't recommend it enough to anyone interested in this area.",
		"The more I learn about this topic, the more fascinated I become. There's always something new to discover.",
		"I tried this approach and was amazed by the results. It's definitely worth giving a shot.",
		"After months of research and experimentation, here's what I've learned...",
		"This is something that many people struggle with, so I wanted to create a comprehensive guide.",
		"The key takeaway here is that consistency and patience are essential for success in this area.",
	}

	placeholderImages = []string{
		"https://images.unsplash.com/photo-1496442226666-8d4d0e62e6e9?w=400",
		"https://images.unsplash.com/photo-1517694712202-14dd9538aa97?w=400",
		"https://images.unsplash.com/photo-1552664730-d307ca884978?w=400",
		"https://images.unsplash.com/photo-1492684223066-81342ee5ff30?w=400",
		"https://images.unsplash.com/photo-1516534775068-bb557ca6b371?w=400",
		"https://images.unsplash.com/photo-1522202176988-41dc08e5ddca?w=400",
		"https://images.unsplash.com/photo-1552668743-c1ad1d66fcab?w=400",
		"https://images.unsplash.com/photo-1540575467063-178f50002cbc?w=400",
		"https://images.unsplash.com/photo-1513635269975-59663e0ac1ad?w=400",
		"https://images.unsplash.com/photo-1549887534-b9b5e2f6f0b6?w=400",
	}

	commentTexts = []string{
		"Great post! Really helpful information.",
		"I completely agree with you on this.",
		"Thanks for sharing! This is exactly what I needed.",
		"Interesting perspective. I hadn't thought about it that way.",
		"This is very informative. Thanks for taking the time to write this.",
		"I've had a similar experience. Glad I'm not the only one!",
		"Well written and easy to understand. Love it!",
		"This really resonates with me. Thanks for the insights.",
		"Excellent breakdown of the topic. Much appreciated!",
		"I learned something new today. Thank you for this!",
	}

	categories = []string{"General", "Lifestyle", "Health & Fitness", "Travel", "Food & Cooking", "Education", "Business", "Finance", "Entertainment", "Sports", "Personal Dev", "Culture", "News"}
)

// SeedDatabase generates fake data for testing and development
func SeedDatabase() error {
	// Check if data already exists
	var count int
	err := Database.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing users: %v", err)
	}

	if count > 0 {
		fmt.Println("Database already seeded, skipping...")
		return nil
	}

	fmt.Println("Seeding database with fake data...")

	// Create users
	userIDs, err := createFakeUsers(15)
	if err != nil {
		return fmt.Errorf("failed to create fake users: %v", err)
	}

	// Create posts
	postIDs, err := createFakePosts(userIDs, 30)
	if err != nil {
		return fmt.Errorf("failed to create fake posts: %v", err)
	}

	// Create comments
	err = createFakeComments(userIDs, postIDs, 100)
	if err != nil {
		return fmt.Errorf("failed to create fake comments: %v", err)
	}

	// Create reactions
	err = createFakeReactions(userIDs, postIDs)
	if err != nil {
		return fmt.Errorf("failed to create fake reactions: %v", err)
	}

	fmt.Println("Database seeding completed successfully!")
	return nil
}

func createFakeUsers(count int) ([]int, error) {
	var userIDs []int

	for i := 0; i < count; i++ {
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		username := strings.ToLower(firstName + lastName + fmt.Sprintf("%d", rand.Intn(100)))
		email := strings.ToLower(firstName) + "." + strings.ToLower(lastName) + fmt.Sprintf("%d@example.com", rand.Intn(1000))

		result, err := Database.Exec(
			"INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
			username,
			email,
			"$2a$10$dummyhashedpassword", // dummy hashed password
		)
		if err != nil {
			continue // Skip if user already exists
		}

		id, err := result.LastInsertId()
		if err != nil {
			continue
		}

		userIDs = append(userIDs, int(id))
	}

	fmt.Printf("Created %d fake users\n", len(userIDs))
	return userIDs, nil
}

func createFakePosts(userIDs []int, count int) ([]int, error) {
	var postIDs []int

	for i := 0; i < count; i++ {
		userID := userIDs[rand.Intn(len(userIDs))]
		title := postTitles[rand.Intn(len(postTitles))]
		text := postTexts[rand.Intn(len(postTexts))]
		image := placeholderImages[rand.Intn(len(placeholderImages))]
		createdAt := time.Now().Add(-time.Duration(rand.Intn(7*24)) * time.Hour)

		result, err := Database.Exec(
			"INSERT INTO posts (user_id, created_at, title, text, image) VALUES (?, ?, ?, ?, ?)",
			userID,
			createdAt,
			title,
			text,
			image,
		)
		if err != nil {
			fmt.Printf("Error creating post: %v\n", err)
			continue
		}

		postID, err := result.LastInsertId()
		if err != nil {
			continue
		}

		postIDs = append(postIDs, int(postID))

		// Assign random categories to post (1-3 categories per post)
		numCategories := rand.Intn(3) + 1
		usedCategories := make(map[int]bool)

		for j := 0; j < numCategories; j++ {
			categoryIdx := rand.Intn(len(categories))
			if !usedCategories[categoryIdx] {
				usedCategories[categoryIdx] = true
				categoryName := categories[categoryIdx]

				// Get category ID
				var categoryID int
				err := Database.QueryRow("SELECT id FROM category WHERE name = ?", categoryName).Scan(&categoryID)
				if err == nil {
					Database.Exec(
						"INSERT OR IGNORE INTO post_category (post_id, category_id) VALUES (?, ?)",
						postID,
						categoryID,
					)
				}
			}
		}
	}

	fmt.Printf("Created %d fake posts\n", len(postIDs))
	return postIDs, nil
}

func createFakeComments(userIDs []int, postIDs []int, count int) error {
	createdCount := 0

	for i := 0; i < count; i++ {
		if len(userIDs) == 0 || len(postIDs) == 0 {
			break
		}

		userID := userIDs[rand.Intn(len(userIDs))]
		postID := postIDs[rand.Intn(len(postIDs))]
		commentText := commentTexts[rand.Intn(len(commentTexts))]
		createdAt := time.Now().Add(-time.Duration(rand.Intn(7*24)) * time.Hour)

		_, err := Database.Exec(
			"INSERT INTO comments (user_id, post_id, created_at, text) VALUES (?, ?, ?, ?)",
			userID,
			postID,
			createdAt,
			commentText,
		)
		if err != nil {
			continue
		}

		createdCount++
	}

	fmt.Printf("Created %d fake comments\n", createdCount)
	return nil
}

func createFakeReactions(userIDs []int, postIDs []int) error {
	postReactionsCount := 0
	commentReactionsCount := 0

	// Create post reactions
	for _, postID := range postIDs {
		// Generate 3-15 reactions per post
		numReactions := rand.Intn(13) + 3
		usedUsers := make(map[int]bool)

		for j := 0; j < numReactions && len(usedUsers) < len(userIDs); j++ {
			userID := userIDs[rand.Intn(len(userIDs))]

			if !usedUsers[userID] {
				usedUsers[userID] = true
				isLike := rand.Intn(10) // 70% likes, 30% dislikes
				if isLike < 7 {
					isLike = 1
				} else {
					isLike = -1
				}

				_, err := Database.Exec(
					"INSERT OR IGNORE INTO post_reactions (user_id, post_id, is_like) VALUES (?, ?, ?)",
					userID,
					postID,
					isLike,
				)
				if err == nil {
					postReactionsCount++
				}
			}
		}
	}

	// Create comment reactions
	rows, err := Database.Query("SELECT id FROM comments")
	if err == nil {
		defer rows.Close()
		var commentIDs []int

		for rows.Next() {
			var commentID int
			if err := rows.Scan(&commentID); err == nil {
				commentIDs = append(commentIDs, commentID)
			}
		}

		for _, commentID := range commentIDs {
			// Generate 0-10 reactions per comment
			numReactions := rand.Intn(11)
			usedUsers := make(map[int]bool)

			for j := 0; j < numReactions && len(usedUsers) < len(userIDs); j++ {
				userID := userIDs[rand.Intn(len(userIDs))]

				if !usedUsers[userID] {
					usedUsers[userID] = true
					isLike := rand.Intn(10)
					if isLike < 7 {
						isLike = 1
					} else {
						isLike = -1
					}

					_, err := Database.Exec(
						"INSERT OR IGNORE INTO comment_reactions (user_id, comment_id, is_like) VALUES (?, ?, ?)",
						userID,
						commentID,
						isLike,
					)
					if err == nil {
						commentReactionsCount++
					}
				}
			}
		}
	}

	fmt.Printf("Created %d post reactions and %d comment reactions\n", postReactionsCount, commentReactionsCount)
	return nil
}
