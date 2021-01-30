package model

import (
	"errors"
	"html"
	"strings"

	"gorm.io/gorm"
)

// Post contains the blog post details
type Post struct {
	gorm.Model
	Title    string `gorm:"size:100;not null;unique;" json:"title"`
	Content  string `gorm:"size:255;not null;" json:"content"`
	AuthorID uint   `json:"author_id"`
	Author   User   `json:"author"`
}

// Prepare Escapes and Trims title and content
func (p *Post) Prepare() {
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
}

// Validate checks required fields
func (p *Post) Validate() error {
	if p.Title == "" {
		return errors.New("Required: Title")
	}
	if p.Content == "" {
		return errors.New("Required: Content")
	}
	if p.AuthorID < 1 {
		return errors.New("Required: Author")
	}
	return nil
}

// CreatePost Inserts new post row in the Post Table
func (p *Post) CreatePost(db *gorm.DB) (*Post, error) {
	if err := db.Create(&p).Error; err != nil {
		return &Post{}, err
	}
	// Todo: if p.id != 0 return fresh pull, check for errors
	return p, nil
}

// ReadAllPosts returns all records from the Post Table
func (p *Post) ReadAllPosts(db *gorm.DB) (*[]Post, error) {
	var posts []Post
	if err := db.Find(&posts).Error; err != nil {
		return &[]Post{}, err
	}

	// Assembles the Author
	if len(posts) > 0 {
		for i, post := range posts {
			if err := db.Model(&User{}).Where("id = ?", post.AuthorID).Take(&posts[i].Author).Error; err != nil {
				return &[]Post{}, err
			}
		}
	}
	return &posts, nil
}

// ReadPostByID queries the Post table by supplied ID returns match
func (p *Post) ReadPostByID(db *gorm.DB, id uint) (*Post, error) {
	var err error
	if err = db.Take(&p, id).Error; err != nil {
		return &Post{}, err
	}
	// Post not found Error, otherwise error is....?
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Post{}, errors.New("Post Not Found")
	}

	// Assembles the Author
	if p.ID != 0 {
		if err = db.Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error; err != nil {
			return &Post{}, err
		}
	}
	return p, err
}

// UpdatePost saves columns
func (p *Post) UpdatePost(db *gorm.DB) (*Post, error) {
	res := db.Model(&p).Updates(map[string]interface{}{
		"title":     p.Title,
		"content":   p.Content,
		"author":    p.Author,
		"author_id": p.AuthorID,
	})

	var err error
	if err = res.Error; err != nil {
		return &Post{}, err
	}

	// Fresh pull and assemble author
	postUpdated := &Post{}
	if err = db.Take(&postUpdated, p.ID).Error; err != nil {
		return &Post{}, err
	}
	if postUpdated.ID != 0 {
		if err = db.Model(&User{}).Where("id = ?", postUpdated.ID).Take(&postUpdated.Author).Error; err != nil {
			return &Post{}, err
		}
	}

	return postUpdated, err
}

// DeletePost sets post row inactive (will need to purge), the return value is used for Testing suite to check isDeleted = 1)
func (p *Post) DeletePost(db *gorm.DB, id uint) (int64, error) {
	res := db.Delete(&p, id)
	if err := res.Error; err != nil {
		// check for error type recordNotFound, respond accordingly
		return res.RowsAffected, err
	}
	return res.RowsAffected, nil
}
