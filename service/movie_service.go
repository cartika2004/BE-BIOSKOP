package service

import (
	"context"
	"encoding/json"
	"project-bioskop/models"
	"project-bioskop/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type MovieService interface {
	GetMovies() ([]models.Movie, error)
}

type movieService struct {
	repo  repository.MovieRepository
	redis *redis.Client
}

func NewMovieService(repo repository.MovieRepository, redis *redis.Client) MovieService {
	return &movieService{repo, redis}
}

func (s *movieService) GetMovies() ([]models.Movie, error) {
	ctx := context.Background()
	cacheKey := "list_movies"

	val, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var movies []models.Movie
		json.Unmarshal([]byte(val), &movies)
		return movies, nil
	}

	// Kalau KOSONG (Miss), ambil dari SQL Server
	movies, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	dataJSON, _ := json.Marshal(movies)
	s.redis.Set(ctx, cacheKey, dataJSON, 60*time.Second)

	return movies, nil
}