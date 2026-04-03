package api

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/planitaicojp/jplaw-cli/internal/model"
)

type LawsParams struct {
	LawID            string
	LawNum           string
	LawNumEra        model.Era
	LawTitle         string
	LawType          []model.LawType
	CategoryCd       []string
	PromulgationFrom string
	PromulgationTo   string
	RepealStatus     []model.RepealStatus
	Limit            int
	Offset           int
	Order            string
}

func (c *Client) ListLaws(params *LawsParams) (*model.LawsResponse, error) {
	q := url.Values{}
	q.Set("response_format", "json")
	if params != nil {
		if params.LawID != "" {
			q.Set("law_id", params.LawID)
		}
		if params.LawNum != "" {
			q.Set("law_num", params.LawNum)
		}
		if params.LawNumEra != "" {
			q.Set("law_num_era", string(params.LawNumEra))
		}
		if params.LawTitle != "" {
			q.Set("law_title", params.LawTitle)
		}
		for _, lt := range params.LawType {
			q.Add("law_type", string(lt))
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
		for _, rs := range params.RepealStatus {
			q.Add("repeal_status", string(rs))
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
	path := fmt.Sprintf("/laws?%s", q.Encode())
	var resp model.LawsResponse
	if err := c.Get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
