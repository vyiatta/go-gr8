package requests

import (
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type CodeRequest struct {
	Id   uint64 `json:"id" validate:"required,numeric,min=1"`
	Code uint16 `json:"code"  validate:"required,min=1000,max=9999"`
}

type UpdateUserRequest struct {
	FirstName  string `json:"firstName" validate:"required,gte=1,max=40"`
	SecondName string `json:"secondName" validate:"required,gte=1,max=40"`
	Email      string `json:"email" validate:"required,email"`
}

func (r UpdateUserRequest) ToDomainModel() (interface{}, error) {
	return domain.User{
		FirstName:  r.FirstName,
		SecondName: r.SecondName,
		Email:      r.Email,
	}, nil
}
