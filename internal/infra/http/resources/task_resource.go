package resources

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type TaskDto struct {
	Id          uint64            `json:"id"`
	UserId      uint64            `json:"userId"`
	Title       string            `json:"title"`
	Description *string           `json:"description,omitempty"`
	Status      domain.TaskStatus `json:"status"`
	Deadline    *time.Time        `json:"deadline,omitempty"`
	CreatedDate time.Time         `json:"createdDate"`
	UpdatedDate time.Time         `json:"updatedDate"`
}

func (d TaskDto) DomainToDto(t domain.Task) TaskDto {
	return TaskDto{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Deadline:    t.Deadline,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
	}
}
