package api

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/planitaicojp/jplaw-cli/internal/model"
)

type KeywordParams struct {
	Keyword          string
	LawType          []model.LawType
	LawNumEra        model.Era
	Asof             string
	CategoryCd       []string
	PromulgationFrom string
	PromulgationTo   string
	Limit            int
	Offset           int
	Order            string
}

func (c *Client) SearchKeyword(params *KeywordParams) (*model.KeywordResponse, error) {
	q := url.Values{}
	q.Set("response_format", "json")
	if params != nil {
		q.Set("keyword", params.Keyword)
		for _, lt := range params.LawType {
			q.Add("law_type", string(lt))
		}
		if params.LawNumEra != "" {
			q.Set("law_num_era", string(params.LawNumEra))
		}
		if params.Asof != "" {
			q.Set("asof", params.Asof)
		}
		for _, cd := range params.CategoryCd {
			q.Add("category_cd", cd)
		}
		if params.PromulgationFrom != "" {
			q.Set("promulgation_date_from", params.PromulgationFrom)
		}
		if params.PromulgationTo != "" {
			q.Set("promulgation_date_to", params.PromulgationTo)
		}
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Offset > 0 {
			q.Set("offset", strconv.Itoa(params.Offset))
		}
		if params.Order != "" {
			q.Set("order", params.Order)
		}
	}
	path := fmt.Sprintf("/keyword?%s", q.Encode())
	var resp model.KeywordResponse
	if err := c.Get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
