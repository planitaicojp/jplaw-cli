package api

import (
	"fmt"
	"net/url"

	"github.com/planitaicojp/jplaw-cli/internal/model"
)

type LawDataParams struct {
	Asof              string
	Elm               string
	JsonFormat        string
	LawFullTextFormat string
}

func (c *Client) GetLawData(idOrNum string, params *LawDataParams) (*model.LawDataResponse, error) {
	q := url.Values{}
	q.Set("response_format", "json")
	q.Set("law_full_text_format", "json")
	q.Set("json_format", "full")
	if params != nil {
		if params.Asof != "" {
			q.Set("asof", params.Asof)
		}
		if params.Elm != "" {
			q.Set("elm", params.Elm)
		}
		if params.JsonFormat != "" {
			q.Set("json_format", params.JsonFormat)
		}
		if params.LawFullTextFormat != "" {
			q.Set("law_full_text_format", params.LawFullTextFormat)
		}
	}
	path := fmt.Sprintf("/law_data/%s?%s", url.PathEscape(idOrNum), q.Encode())
	var resp model.LawDataResponse
	if err := c.Get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetLawFile(fileType, idOrNum, asof string) ([]byte, error) {
	q := url.Values{}
	if asof != "" {
		q.Set("asof", asof)
	}
	path := fmt.Sprintf("/law_file/%s/%s", url.PathEscape(fileType), url.PathEscape(idOrNum))
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	return c.GetRaw(path)
}
