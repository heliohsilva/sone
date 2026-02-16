package models

import (
	"errors"
	"strings"
	"time"
)

type Post struct {
	ID        int64     `json:"id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	UserID    uint      `json:"userID,omitempty"`
	UserNick  string    `json:"userNick,omitempty"`
	Likes     uint      `json:"likes"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

func (post *Post) Prepare() error {
	if err := post.Validate(); err != nil {
		return err
	}

	post.Format()
	return nil
}

func (post *Post) Validate() error {
	if post.Title == "" {
		return errors.New("the title cannot be empty ")
	}

	if post.Content == "" {
		return errors.New("Content cannot be empty")
	}

	return nil
}

func (post *Post) Format() {
	post.Title = strings.TrimSpace(post.Title)
	post.Content = strings.TrimSpace(post.Content)
}
