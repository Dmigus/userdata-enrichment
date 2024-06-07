package get

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"enrichstorage/internal/service/enrichstorage/get"
	"enrichstorage/pkg/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	GetterService interface {
		GetWithPaging(ctx context.Context, req get.Request) (get.Result, error)
	}
	Handler struct {
		getter GetterService
	}
)

func NewHandler(getter GetterService) *Handler {
	return &Handler{getter: getter}
}

func (h *Handler) Handle(c *gin.Context) {
	var req request
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, c.Error(err))
		return
	}
	usecaseReq, err := req.ToUsecaseRequest()
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Error(err))
		return
	}
	results, err := h.getter.GetWithPaging(c, *usecaseReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Error(err))
		return
	}
	resp := newResponse(c, req, results)
	c.JSON(http.StatusOK, resp)
}

func unmarshalFIO(str string) (*types.FIO, error) {
	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	var fioPlain struct {
		Surname    string `json:"surname"`
		Name       string `json:"name"`
		Patronymic string `json:"patronymic"`
	}
	if err = json.Unmarshal(bytes, &fioPlain); err != nil {
		return nil, err
	}
	fio, err := types.NewFIO(fioPlain.Name, fioPlain.Surname, fioPlain.Patronymic)
	if err != nil {
		return nil, err
	}
	return &fio, nil
}

func marshallFIO(fio types.FIO) string {
	fioPlain := struct {
		Surname    string `json:"surname"`
		Name       string `json:"name"`
		Patronymic string `json:"patronymic"`
	}{
		fio.Surname(), fio.Name(), fio.Patronymic(),
	}
	bytes, _ := json.Marshal(fioPlain)
	encoded := base64.StdEncoding.EncodeToString(bytes)
	return encoded
}
