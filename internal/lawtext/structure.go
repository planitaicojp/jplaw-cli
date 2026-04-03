package lawtext

import "encoding/json"

type Node struct {
	Tag      string            `json:"tag"`
	Attr     map[string]string `json:"attr"`
	Children []json.RawMessage `json:"children"`
}

func parseNode(raw json.RawMessage) (*Node, error) {
	var node Node
	if err := json.Unmarshal(raw, &node); err != nil {
		return nil, err
	}
	return &node, nil
}

func childText(children []json.RawMessage) string {
	var result string
	for _, child := range children {
		var s string
		if json.Unmarshal(child, &s) == nil {
			result += s
			continue
		}
		node, err := parseNode(child)
		if err != nil {
			continue
		}
		result += childText(node.Children)
	}
	return result
}
