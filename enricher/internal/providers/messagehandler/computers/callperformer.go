package computers

import (
	"context"
	"errors"
	"io"
	"net/http"
)

var (
	errNotOKStatusCode = errors.New("resource answered not OK status code")
	errWrongBody       = errors.New("incorrect body")
)

type HttpQueryPerformer struct {
	client http.Client
}

func NewHttpQueryPerformer(client http.Client) *HttpQueryPerformer {
	return &HttpQueryPerformer{client: client}
}

func (h *HttpQueryPerformer) PerformGetReq(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errNotOKStatusCode
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return io.ReadAll(resp.Body)
}
