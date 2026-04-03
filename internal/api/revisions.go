package api

import (
	"fmt"
	"net/url"

	"github.com/planitaicojp/jplaw-cli/internal/model"
)

type RevisionsParams struct {
	AmendmentDateFrom string
	AmendmentDateTo   string
}

func (c *Client) GetRevisions(idOrNum string, params *RevisionsParams) (*model.RevisionsResponse, error) {
	q := url.Values{}
	q.Set("response_format", "json")
	if params != nil {
		if params.AmendmentDateFrom != "" {
			q.Set("amendment_date_from", params.AmendmentDateFrom)
		}
		if params.AmendmentDateTo != "" {
			q.Set("amendment_date_to", params.AmendmentDateTo)
		}
	}
	path := fmt.Sprintf("/law_revisions/%s?%s", url.PathEscape(idOrNum), q.Encode())
	var resp model.RevisionsResponse
	if err := c.Get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
