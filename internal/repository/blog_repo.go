package repository

import (
	"github.com/Hamiduzzaman96/Blog-Service/internal/domain"
	"gorm.io/gorm"
)

type BlogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

func (r *BlogRepository) Migrate() error { //Database schema ensure
	return r.db.AutoMigrate(&BlogModel{}) /*GORM এর AutoMigrate:
	যদি table না থাকে → create করে
	যদি column না থাকে → add করে
	যদি column type change করা safe হয় → update করে*/
}

// MAPPERS
// func blogModelToDomain(m *BlogModel) *domain.BlogPost {
// 	return &domain.BlogPost{
// 		ID:        m.ID,
// 		AuthorID:  m.AuthorID,
// 		Title:     m.Title,
// 		Content:   m.Content,
// 		CreatedAt: m.CreatedAt,
// 		UpdatedAt: m.UpdatedAt,
// 	}
// }

func blogDomainToModel(b *domain.BlogPost) *BlogModel {
	return &BlogModel{
		ID:       b.ID,
		AuthorID: b.AuthorID,
		Title:    b.Title,
		Content:  b.Content,
	}
}

// CRUD

func (r *BlogRepository) Create(b *domain.BlogPost) (*domain.BlogPost, error) {
	m := blogDomainToModel(b)

	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	b.ID = m.ID
	return b, nil
}
