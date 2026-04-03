package api

import (
	"fmt"
	"net/url"
)

func (c *Client) GetAttachment(revisionID, src string) ([]byte, error) {
	q := url.Values{}
	if src != "" {
		q.Set("src", src)
	}
	path := fmt.Sprintf("/attachment/%s", url.PathEscape(revisionID))
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	return c.GetRaw(path)
}
