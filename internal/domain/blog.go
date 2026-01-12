package domain

import "time"

type BlogPost struct {
	ID        uint
	AuthorID  uint
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBlogPost(authorId uint, title, content string) *BlogPost {
	return &BlogPost{
		AuthorID:  authorId,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (b *BlogPost) Update(title, content string) {
	b.Title = title
	b.Content = content
	b.UpdatedAt = time.Now()
}
