package get

import (
	"enrichstorage/internal/service/enrichstorage/get"
	"enrichstorage/pkg/types"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type (
	responseRecord struct {
		Surname     string `json:"surname"`
		Name        string `json:"name"`
		Patronymic  string `json:"patronymic"`
		Sex         string `json:"sex"`
		Nationality string `json:"nationality"`
		Age         int    `json:"age"`
	}

	response struct {
		Data   []responseRecord `json:"data"`
		Paging struct {
			Previous *string `json:"previous,omitempty"`
			Next     *string `json:"next,omitempty"`
		} `json:"paging"`
	}
)

func newResponse(c *gin.Context, req request, res get.Result) response {
	resp := response{}
	data := lo.Map(res.Records, func(resultRecord types.EnrichedRecord, _ int) responseRecord {
		return newResponseRecord(resultRecord)
	})
	resp.Data = data
	resp.Paging.Previous = compPrevURL(c, req, res)
	resp.Paging.Next = compNextURL(c, req, res)
	return resp
}

func compPrevURL(c *gin.Context, req request, res get.Result) *string {
	if res.PrevPage == nil {
		return nil
	}
	url := c.Request.URL
	query := url.Query()
	fio := types.FIO(*res.PrevPage)
	query.Del("after")
	query.Set("before", marshallFIO(fio))
	url.RawQuery = query.Encode()
	url.Host = c.Request.Host
	url.Scheme = "http"
	if c.Request.TLS != nil {
		url.Scheme = "https"
	}
	result := url.String()
	return &result
}

func compNextURL(c *gin.Context, req request, res get.Result) *string {
	if res.NextPage == nil {
		return nil
	}
	url := *c.Request.URL
	query := url.Query()
	fio := types.FIO(*res.NextPage)
	query.Del("before")
	query.Set("after", marshallFIO(fio))
	url.RawQuery = query.Encode()
	url.Host = c.Request.Host
	url.Scheme = "http"
	if c.Request.TLS != nil {
		url.Scheme = "https"
	}
	result := url.String()
	return &result
}

func newResponseRecord(res types.EnrichedRecord) responseRecord {
	return responseRecord{
		Surname:     res.Fio.Surname(),
		Name:        res.Fio.Name(),
		Patronymic:  res.Fio.Patronymic(),
		Sex:         res.Sex,
		Nationality: res.Nationality,
		Age:         res.Age,
	}
}
