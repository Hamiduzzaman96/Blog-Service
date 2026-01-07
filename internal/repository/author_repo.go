package repository

import (
	"github.com/Hamiduzzaman96/Blog-Service/internal/domain"
	"gorm.io/gorm"
)

type AuthorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

func (r *AuthorRepository) Migrate() error {
	return r.db.AutoMigrate(&AuthorModel{})
}

func authorModelToDomain(m *AuthorModel) *domain.Author {
	return &domain.Author{
		ID:     m.ID,
		UserID: m.UserID,
	}
}

func authorDomainToModel(a *domain.Author) *AuthorModel {
	return &AuthorModel{
		ID:     a.ID,
		UserID: a.UserID,
	}
}

func (r *AuthorRepository) Create(a *domain.Author) (*domain.Author, error) {
	m := authorDomainToModel(a)

	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}

	a.ID = m.ID
	return a, nil
}

func (r *AuthorRepository) FindByUserID(userID uint) (*domain.Author, error) {
	var m AuthorModel

	if err := r.db.Where("user_id = ?", userID).First(&m).Error; err != nil { //First matched row নেয়
		return nil, err //SELECT * FROM authors WHERE user_id = 10 LIMIT 1;
	}
	return authorModelToDomain(&m), nil
}
