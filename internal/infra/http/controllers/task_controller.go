package controllers

import (
	"log"
	"net/http"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
)

type TaskController struct {
	taskService app.TaskService
}

func NewTaskController(ts app.TaskService) TaskController {
	return TaskController{
		taskService: ts,
	}
}

func (c TaskController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController -> Save: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		task.UserId = user.Id
		task.Status = domain.NewTaskStatus
		task, err = c.taskService.Save(task)
		if err != nil {
			log.Printf("TaskController -> Save: %s", err)
			InternalServerError(w, err)
			return
		}

		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(task)
		Created(w, tDto)
	}
}

func (c TaskController) FindByUserId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		tasks, err := c.taskService.FindByUserId(user.Id)
		if err != nil {
			log.Printf("TaskController -> FindByUserId: %s", err)
			InternalServerError(w, err)
			return
		}

		var tsDto resources.TasksDto
		tsDto = tsDto.DomainToDtoCollection(tasks)
		Success(w, tsDto)
	}
}

func (c TaskController) FindByTaskId(taskId uint64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := c.taskService.FindByTaskId(taskId)
		if err != nil {
			log.Printf("TaskController -> FindByTaskId: %s", err)
			InternalServerError(w, err)
			return
		}

		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(task)
		Success(w, tDto)
	}
}

func (c TaskController) Update(taskId uint64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController -> Update: %s", err)
			BadRequest(w, err)
			return
		}

		task.Id = taskId
		updatedTask, err := c.taskService.Update(task)
		if err != nil {
			log.Printf("TaskController -> Update: %s", err)
			InternalServerError(w, err)
			return
		}

		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(updatedTask)
		Success(w, tDto)
	}
}

func (c TaskController) Delete(taskId uint64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := c.taskService.DeleteByTaskId(taskId)
		if err != nil {
			log.Printf("TaskController -> Delete: %s", err)
			InternalServerError(w, err)
			return
		}
	}
}
