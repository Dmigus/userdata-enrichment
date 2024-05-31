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
	NationalityComputer struct {
		urlTemplate   *url.URL
		callPerformer CallPerformer
	}
	nationalizeResult struct {
		CountryField []struct {
			CountryId string `json:"country_id"`
		} `json:"country"`
	}
)

func NewNationalityComputer(hostname string, callPerformer CallPerformer) (*NationalityComputer, error) {
	urlTemplate, err := url.Parse(hostname)
	if err != nil {
		return nil, err
	}
	return &NationalityComputer{urlTemplate: urlTemplate, callPerformer: callPerformer}, nil
}

func (a *NationalityComputer) Get(ctx context.Context, fio types.FIO) (types.Nationality, error) {
	queryURL := a.getURL(fio)
	bodyBytes, err := a.callPerformer.PerformGetReq(ctx, queryURL)
	if err != nil {
		return "", fmt.Errorf("error getting info from nationalize: %w", err)
	}
	var answer nationalizeResult
	err = json.Unmarshal(bodyBytes, &answer)
	if err != nil || len(answer.CountryField) == 0 {
		return "", fmt.Errorf("error getting info from nationalize: %w", errWrongBody)
	}
	mostFreqCountry := answer.CountryField[0]
	return mostFreqCountry.CountryId, nil
}

func (a *NationalityComputer) getURL(fio types.FIO) string {
	queryURL := *a.urlTemplate
	q := queryURL.Query()
	q.Add(nameKey, fio.Name())
	queryURL.RawQuery = q.Encode()
	return queryURL.String()
}
