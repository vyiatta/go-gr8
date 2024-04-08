package app

import (
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"log"
	"time"
)

type AuthService interface {
	Register(user domain.User) (domain.User, error)
	Login(user domain.User) (domain.User, error)
	Logout(sess domain.Session) error
	Check(sess domain.Session) error
	GenerateJwt(user domain.User) (string, error)
}

type authService struct {
	authRepo    database.SessionRepository
	userService UserService
	tokenAuth   *jwtauth.JWTAuth
	jwtTTL      time.Duration
}

func NewAuthService(ar database.SessionRepository, us UserService, ta *jwtauth.JWTAuth, jwtTtl time.Duration) AuthService {
	return authService{
		authRepo:    ar,
		userService: us,
		tokenAuth:   ta,
		jwtTTL:      jwtTtl,
	}
}

func (s authService) Register(user domain.User) (domain.User, error) {
	user, err := s.userService.Save(user)
	if err != nil {
		log.Print(err)
		return domain.User{}, err
	}
	return user, err
}

func (s authService) Login(user domain.User) (domain.User, error) {
	user, err := s.userService.Update(user)
	if err != nil {
		log.Print(err)
		return domain.User{}, err
	}
	return user, err
}

func (s authService) CheckCode(user domain.User) (domain.User, error) {
	user, err := s.userService.Save(user)
	if err != nil {
		log.Print(err)
		return domain.User{}, err
	}
	return user, err
}

func (s authService) Logout(sess domain.Session) error {
	return s.authRepo.Delete(sess)
}

func (s authService) GenerateJwt(user domain.User) (string, error) {
	sess := domain.Session{UserId: user.Id, UUID: uuid.New()}
	err := s.authRepo.Save(sess)
	if err != nil {
		log.Printf("AuthService: failed to save session %s", err)
		return "", err
	}

	claims := map[string]interface{}{
		"user_id": sess.UserId,
		"uuid":    sess.UUID,
	}
	jwtauth.SetExpiryIn(claims, s.jwtTTL)
	_, tokenString, err := s.tokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s authService) Check(sess domain.Session) error {
	return s.authRepo.Exists(sess)
}
