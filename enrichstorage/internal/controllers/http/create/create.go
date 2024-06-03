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
		Name       string `json:"name"`
		Surname    string `json:"surname"`
		Patronymic string `json:"patronymic"`
	}
	creator interface {
		Create(ctx context.Context, fio types.FIO) error
	}
	CreateHandler struct {
		creator creator
	}
)

func (ch *CreateHandler) Handle(c *gin.Context) {
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
	}
	c.Status(http.StatusOK)
}

func recToFio(rec request) (types.FIO, error) {
	return types.NewFIO(rec.Name, rec.Surname, rec.Patronymic)
}
