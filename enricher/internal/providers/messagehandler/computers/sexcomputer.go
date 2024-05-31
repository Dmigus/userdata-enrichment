// Package computers предназначен для реализации вычислителей по ФИО
package computers

import (
	"context"
	"encoding/json"
	"enrichstorage/pkg/types"
	"fmt"
	"net/url"
)

type (
	SexComputer struct {
		urlTemplate   *url.URL
		callPerformer CallPerformer
	}
)

func NewSexComputer(hostname string, callPerformer CallPerformer) (*SexComputer, error) {
	urlTemplate, err := url.Parse(hostname)
	if err != nil {
		return nil, err
	}
	return &SexComputer{urlTemplate: urlTemplate, callPerformer: callPerformer}, nil
}

func (a *SexComputer) Get(ctx context.Context, fio types.FIO) (types.Sex, error) {
	queryURL := a.getURL(fio)
	bodyBytes, err := a.callPerformer.PerformGetReq(ctx, queryURL)
	if err != nil {
		return "", fmt.Errorf("error getting info from genderize: %w", err)
	}
	var answer struct {
		GenderField string `json:"gender"`
	}
	err = json.Unmarshal(bodyBytes, &answer)
	if err != nil {
		return "", fmt.Errorf("error getting info from genderize: %w", errWrongBody)
	}
	return answer.GenderField, nil
}

func (a *SexComputer) getURL(fio types.FIO) string {
	queryURL := *a.urlTemplate
	q := queryURL.Query()
	q.Add(nameKey, fio.Name())
	queryURL.RawQuery = q.Encode()
	return queryURL.String()
}
