package controllers

import (
	"errors"
	"fmt"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
	"github.com/upper/db/v4"
	"log"
	"net/http"
)

type AuthController struct {
	authService app.AuthService
	userService app.UserService
}

func NewAuthController(as app.AuthService, us app.UserService) AuthController {
	return AuthController{
		authService: as,
		userService: us,
	}
}

func (c AuthController) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.AuthRequest{}, domain.User{})
		if err != nil {
			log.Printf("AuthController: %s", err)
			BadRequest(w, err)
			return
		}

		userExist, err := c.userService.FindByPhone(user.Phone)
		if err != nil {
			log.Printf("AuthController: %s", err)
			BadRequest(w, err)
			return
		}

		var u domain.User
		u, err = c.authService.Login(userExist)
		if err != nil {
			log.Printf("AuthController: %s", err)
			InternalServerError(w, err)
			return
		}

		var uDto resources.SimpleUserDto
		Success(w, uDto.DomainToDto(u))
	}
}

func (c AuthController) CheckNewPhoneCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.CodeNewPhoneRequest{}, domain.User{})
		if err != nil {
			log.Printf("AuthController: %s", err)
			BadRequest(w, err)
			return
		}

		u := r.Context().Value(UserKey).(domain.User)
		var userExists domain.User
		userExists, err = c.userService.FindByPhone(user.Phone)
		if err != nil && !errors.Is(err, db.ErrNoMoreRows) {
			log.Printf("AuthController: %s", err)
			BadRequest(w, err)
			return
		}
		if userExists.Id != 0 {
			err = fmt.Errorf("%s phone number belongs to another user", user.Phone)
			log.Printf("AuthController: %s", err)
			BadRequest(w, err)
			return
		}

		u.Phone = user.Phone
		u, err = c.userService.Update(u)
		if err != nil {
			log.Printf("UserController: %s", err)
			InternalServerError(w, err)
			return
		}

		sess := r.Context().Value(SessKey).(domain.Session)
		err = c.authService.Logout(sess)
		if err != nil {
			log.Printf("AuthController: %s", err)
			InternalServerError(w, err)
			return
		}

		var token string
		token, err = c.authService.GenerateJwt(u)
		if err != nil {
			log.Printf("AuthController: %s", err)
			BadRequest(w, err)
			return
		}

		var authDto resources.AuthDto
		Success(w, authDto.DomainToDto(token, u))
	}
}

func (c AuthController) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := r.Context().Value(SessKey).(domain.Session)
		err := c.authService.Logout(sess)
		if err != nil {
			log.Printf("AuthController: %s", err)
			InternalServerError(w, err)
			return
		}

		noContent(w)
	}
}
