package service

import (
	"errors"
	"project-bioskop/config"
	"project-bioskop/models"
	"project-bioskop/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(input models.User) error
	Login(email, password string) (string, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Register(input models.User) error {
	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	input.Password = string(hashedPassword)

	// Default Role = user
	input.Role = "user"

	return s.repo.Create(&input)
}

func (s *authService) Login(email, password string) (string, error) {
	// Cari user berdasarkan email
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("email not found")
	}

	// Cek password cocok gak?
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("wrong password")
	}

	// Bikin JWT Token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.AppConfig.JWT.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}