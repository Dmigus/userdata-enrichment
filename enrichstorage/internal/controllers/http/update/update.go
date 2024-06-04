package create

import (
	"context"
	"enrichstorage/internal/service/enrichstorage/update"
	"enrichstorage/pkg/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	request struct {
		Name        string  `json:"name" binding:"required"`
		Surname     string  `json:"surname" binding:"required"`
		Patronymic  string  `json:"patronymic"`
		Age         *int    `json:"age" validate:"gte=0,lte=130"`
		Sex         *string `json:"sex"`
		Nationality *string `json:"nationality"`
	}
	UpdateService interface {
		Update(ctx context.Context, rec update.Request) error
	}
	Handler struct {
		updater UpdateService
	}
)

func NewHandler(updater UpdateService) *Handler {
	return &Handler{updater: updater}
}

func (ch *Handler) Handle(c *gin.Context) {
	rec := request{}
	if err := c.ShouldBindJSON(&rec); err != nil {
		c.JSON(http.StatusBadRequest, c.Error(err))
		return
	}
	updRec, err := recToUpdateRequest(rec)
	if err != nil {
		err = fmt.Errorf("incorrect request: %w", err)
		c.JSON(http.StatusBadRequest, c.Error(err))
		return
	}
	err = ch.updater.Update(c, updRec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, c.Error(err))
	}
	c.Status(http.StatusOK)
}

func recToUpdateRequest(rec request) (update.Request, error) {
	result := update.Request{}
	var err error
	result.Fio, err = types.NewFIO(rec.Name, rec.Surname, rec.Patronymic)
	if err != nil {
		return result, err
	}
	if rec.Age != nil {
		result.AgePresents = true
		result.NewAge = *rec.Age
	}
	if rec.Sex != nil {
		result.SexPresents = true
		result.NewSex = *rec.Sex
	}
	if rec.Nationality != nil {
		result.NationalityPresents = true
		result.NewNat = *rec.Nationality
	}
	return result, nil
}
