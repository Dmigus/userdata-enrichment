package create

import (
	"context"
	"enrichstorage/pkg/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	request struct {
		Name       string `json:"name"  validate:"required"`
		Surname    string `json:"surname"  validate:"required"`
		Patronymic string `json:"patronymic"  validate:"optional"`
	}
	CreatorService interface {
		Create(ctx context.Context, fio types.FIO) error
	}
	Handler struct {
		creator CreatorService
	}
)

func NewHandler(creator CreatorService) *Handler {
	return &Handler{creator: creator}
}

func (ch *Handler) Handle(c *gin.Context) {
	rec := request{}
	if err := c.ShouldBindJSON(&rec); err != nil {
		c.JSON(http.StatusBadRequest, c.Error(err))
		return
	}
	fio, err := recToFio(rec)
	if err != nil {
		err = fmt.Errorf("incorrect request: %w", err)
		c.JSON(http.StatusBadRequest, c.Error(err))
		return
	}
	err = ch.creator.Create(c, fio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, c.Error(err))
		return
	}
	c.Status(http.StatusOK)
}

func recToFio(rec request) (types.FIO, error) {
	return types.NewFIO(rec.Name, rec.Surname, rec.Patronymic)
}
