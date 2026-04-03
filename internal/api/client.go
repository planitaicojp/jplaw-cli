package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	cerrors "github.com/planitaicojp/jplaw-cli/internal/errors"
	"github.com/planitaicojp/jplaw-cli/internal/model"
)

var UserAgent = "jplaw-cli/dev"

const (
	defaultTimeout = 30 * time.Second
	maxRetries     = 3
)

type Client struct {
	HTTP    *http.Client
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{
		HTTP:    &http.Client{Timeout: defaultTimeout},
		BaseURL: baseURL,
	}
}

func (c *Client) Get(path string, result any) error {
	resp, err := c.do(http.MethodGet, path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) GetRaw(path string) ([]byte, error) {
	resp, err := c.do(http.MethodGet, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) do(method, path string) (*http.Response, error) {
	url := c.BaseURL + path

	start := time.Now()

	var resp *http.Response
	for attempt := 0; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return nil, fmt.Errorf("リクエスト作成エラー: %w", err)
		}
		req.Header.Set("User-Agent", UserAgent)
		req.Header.Set("Accept", "application/json")

		if attempt == 0 {
			debugLogRequest(req)
		}

		resp, err = c.HTTP.Do(req)
		if err != nil {
			if attempt == maxRetries {
				return nil, &cerrors.NetworkError{Err: err}
			}
			time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
			continue
		}
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			if attempt < maxRetries {
				resp.Body.Close()
				time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
				continue
			}
		}
		break
	}

	elapsed := time.Since(start)
	if resp == nil {
		return nil, &cerrors.NetworkError{Err: fmt.Errorf("リトライ後もレスポンスなし")}
	}
	debugLogResponse(resp, elapsed)

	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp)
	}
	return resp, nil
}

func parseAPIError(resp *http.Response) error {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var errInfo model.ErrorInfo
	if json.Unmarshal(body, &errInfo) == nil && errInfo.Message != "" {
		if resp.StatusCode == 404 {
			return &cerrors.NotFoundError{Resource: "法令", ID: errInfo.Message}
		}
		return &cerrors.APIError{
			StatusCode: resp.StatusCode,
			Code:       errInfo.Code,
			Message:    errInfo.Message,
		}
	}
	if resp.StatusCode == 404 {
		return &cerrors.NotFoundError{Resource: "法令", ID: ""}
	}
	return &cerrors.APIError{
		StatusCode: resp.StatusCode,
		Message:    string(body),
	}
}
