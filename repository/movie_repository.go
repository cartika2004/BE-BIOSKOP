package repository

import (
	"project-bioskop/models"
	"gorm.io/gorm"
)

type MovieRepository interface {
	FindAll() ([]models.Movie, error)
}

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db}
}

func (r *movieRepository) FindAll() ([]models.Movie, error) {
	var movies []models.Movie
	err := r.db.Find(&movies).Error
	return movies, err
}