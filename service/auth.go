package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"testTask/entity"
	"testTask/repository"
	"time"
)

var (
	jwtKey          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
	Ip     string `json:"ip"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	jwtKey = viper.GetString("jwt.key")
	accessTokenTTL = viper.GetDuration("jwt.access_token_ttl")
	refreshTokenTTL = viper.GetDuration("jwt.refresh_token_ttl")
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user entity.User) (string, error) {
	hashedPassword, err := hash(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = hashedPassword
	user.Id = uuid.New()
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateAccessToken(id, ip string) (string, error) {
	user, err := s.repo.GetUser(id)
	if err != nil {
		return "", err
	}

	currentTime := time.Now()
	expTime := time.Now().Add(accessTokenTTL)
	token := GenerateJWT(currentTime, expTime, user.Id.String(), ip)

	return token.SignedString([]byte(jwtKey))
}

func (s *AuthService) GenerateRefreshToken(id, ip string) (string, error) {
	user, err := s.repo.GetUser(id)
	if err != nil {
		return "", err
	}

	currentTime := time.Now()
	expTime := time.Now().Add(refreshTokenTTL)
	token := GenerateJWT(currentTime, expTime, user.Id.String(), ip)

	tokenSigned, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	hashedToken, err := hash(sha256Hash(tokenSigned))
	if err != nil {
		return "", errors.New("error while hashing refresh token")
	}

	// удаление прошлого рефреш токена из базы
	err = s.repo.DeleteRefreshTokenById(id)

	err = s.repo.CreateRefreshToken(
		entity.RefreshToken{
			TokenHash: hashedToken,
			IpAddress: ip,
			UserId:    id,
			CreatedAt: currentTime,
			ExpiresAt: expTime})

	if err != nil {
		return "", errors.New("error while saving refresh token")
	}

	return tokenSigned, nil
}

func GenerateJWT(currentTime, expTime time.Time, userId, ip string) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
			IssuedAt:  currentTime.Unix(),
		},
		userId,
		ip,
	})
	return token
}

func (s *AuthService) ParseToken(inputToken string) (string, error) {
	token, err := jwt.ParseWithClaims(inputToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(jwtKey), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New(fmt.Sprintf("invalid or expired token: %s", err))
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func (s *AuthService) RefreshTokens(refreshToken, ip string) (string, string, error) {
	userId, err := s.ParseToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// Получаем данные рефреш токена из базы
	tokenData, err := s.repo.GetRefreshTokenById(userId)
	if err != nil {
		return "", "", errors.New("refresh token not found")
	}

	// проверяем совпадение хэшей
	err = bcrypt.CompareHashAndPassword([]byte(tokenData.TokenHash), []byte(sha256Hash(refreshToken)))
	if err != nil {
		return "", "", errors.New("refresh token is invalid")
	}

	// Проверяем совпадение IP-адресов
	if tokenData.IpAddress != ip {
		// TODO: отправка уведомления на почту
		logrus.Println("IP address has changed, sending warning email")
	}

	// Проверяем истечение срока действия Refresh токена
	if time.Now().After(tokenData.ExpiresAt) {
		return "", "", errors.New("refresh token has expired")
	}

	// Генерация нового Access токена
	newAccessToken, err := s.GenerateAccessToken(tokenData.UserId, ip)
	if err != nil {
		return "", "", err
	}

	// Генерация нового Refresh токена
	newRefreshToken, err := s.GenerateRefreshToken(tokenData.UserId, ip)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func hash(input string) (string, error) {
	hashedInput, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedInput), nil
}

// refreshToken оказался слишком большим для хэширования с bcrypt напрямую
// поэтому хэшируем его с помощью sha256 перед bcrypt
func sha256Hash(token string) string {
	hash := sha256.New()
	hash.Write([]byte(token))
	return hex.EncodeToString(hash.Sum(nil))
}
