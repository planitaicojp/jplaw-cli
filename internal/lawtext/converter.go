package lawtext

import (
	"encoding/json"
	"strings"
)

func Convert(raw json.RawMessage) (string, error) {
	if len(raw) == 0 || string(raw) == "{}" || string(raw) == "null" {
		return "", nil
	}
	node, err := parseNode(raw)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	renderNode(&b, node, 0)
	return b.String(), nil
}

func renderNode(b *strings.Builder, node *Node, depth int) {
	if node == nil {
		return
	}

	switch node.Tag {
	case "LawTitle":
		text := childText(node.Children)
		if text != "" {
			b.WriteString(text)
			b.WriteString("\n\n")
		}

	case "TOCLabel", "PartTitle", "ChapterTitle", "SectionTitle",
		"SubsectionTitle", "DivisionTitle":
		text := childText(node.Children)
		if text != "" {
			indent := strings.Repeat("  ", depth)
			b.WriteString(indent)
			b.WriteString(text)
			b.WriteString("\n\n")
		}

	case "ArticleTitle":
		text := childText(node.Children)
		if text != "" {
			b.WriteString(text)
		}

	case "ArticleCaption":
		text := childText(node.Children)
		if text != "" {
			b.WriteString(text)
			b.WriteString("\n")
		}

	case "ParagraphSentence", "ItemSentence", "Subitem1Sentence",
		"Subitem2Sentence", "Subitem3Sentence":
		text := childText(node.Children)
		if text != "" {
			indent := strings.Repeat("  ", depth)
			b.WriteString(indent)
			b.WriteString(text)
			b.WriteString("\n")
		}

	case "ParagraphNum":
		text := childText(node.Children)
		if text != "" {
			indent := strings.Repeat("  ", depth)
			b.WriteString(indent)
			b.WriteString(text)
			b.WriteString(" ")
		}

	case "ItemTitle":
		text := childText(node.Children)
		if text != "" {
			b.WriteString("    ")
			b.WriteString(text)
			b.WriteString(" ")
		}

	case "Subitem1Title", "Subitem2Title", "Subitem3Title":
		text := childText(node.Children)
		if text != "" {
			b.WriteString("      ")
			b.WriteString(text)
			b.WriteString(" ")
		}

	case "SupplProvisionLabel":
		text := childText(node.Children)
		if text != "" {
			b.WriteString("\n")
			b.WriteString(strings.Repeat("─", 40))
			b.WriteString("\n")
			b.WriteString(text)
			b.WriteString("\n\n")
		}

	case "Article":
		for _, child := range node.Children {
			childNode, err := parseNode(child)
			if err != nil {
				continue
			}
			renderNode(b, childNode, depth)
		}
		b.WriteString("\n")
		return

	case "Paragraph":
		for _, child := range node.Children {
			childNode, err := parseNode(child)
			if err != nil {
				continue
			}
			renderNode(b, childNode, 1)
		}
		return

	case "Item":
		for _, child := range node.Children {
			childNode, err := parseNode(child)
			if err != nil {
				continue
			}
			renderNode(b, childNode, 2)
		}
		return

	case "Subitem1", "Subitem2", "Subitem3":
		for _, child := range node.Children {
			childNode, err := parseNode(child)
			if err != nil {
				continue
			}
			renderNode(b, childNode, 3)
		}
		return
	}

	// Default: recurse into children
	for _, child := range node.Children {
		var s string
		if json.Unmarshal(child, &s) == nil {
			continue
		}
		childNode, err := parseNode(child)
		if err != nil {
			continue
		}
		renderNode(b, childNode, depth)
	}
}
