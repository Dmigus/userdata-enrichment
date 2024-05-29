// Package computers предназначен для реализации вычислителей по ФИО
package computers

import (
	"bff/pkg/types"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

const nameKey = "name"

type (
	callPerformer interface {
		PerformGetReq(ctx context.Context, url string) ([]byte, error)
	}
	AgifyComputer struct {
		urlTemplate   *url.URL
		callPerformer callPerformer
	}
)

func NewAgifyComputer(hostname string, callPerformer callPerformer) (*AgifyComputer, error) {
	urlTemplate, err := url.Parse(hostname)
	if err != nil {
		return nil, err
	}
	return &AgifyComputer{urlTemplate: urlTemplate, callPerformer: callPerformer}, nil
}

func (a *AgifyComputer) Get(ctx context.Context, fio types.FIO) (types.Age, error) {
	queryURL := a.getURL(fio)
	bodyBytes, err := a.callPerformer.PerformGetReq(ctx, queryURL)
	if err != nil {
		return 0, fmt.Errorf("error getting info from agify: %w", err)
	}
	var answer struct {
		AgeField int `json:"age"`
	}
	err = json.Unmarshal(bodyBytes, &answer)
	if err != nil {
		return 0, fmt.Errorf("error getting info from agify: %w", errWrongBody)
	}
	return answer.AgeField, nil
}

func (a *AgifyComputer) getURL(fio types.FIO) string {
	queryURL := *a.urlTemplate
	q := queryURL.Query()
	q.Add(nameKey, fio.Name())
	queryURL.RawQuery = q.Encode()
	return queryURL.String()
}
