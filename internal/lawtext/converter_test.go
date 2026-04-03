package lawtext

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestConvertSimpleArticle(t *testing.T) {
	raw := json.RawMessage(`{
		"tag": "Law",
		"attr": {},
		"children": [
			{
				"tag": "LawBody",
				"attr": {},
				"children": [
					{
						"tag": "LawTitle",
						"attr": {},
						"children": ["テスト法"]
					},
					{
						"tag": "MainProvision",
						"attr": {},
						"children": [
							{
								"tag": "Article",
								"attr": {"Num": "1"},
								"children": [
									{
										"tag": "ArticleCaption",
										"attr": {},
										"children": ["（目的）"]
									},
									{
										"tag": "ArticleTitle",
										"attr": {},
										"children": ["第一条"]
									},
									{
										"tag": "Paragraph",
										"attr": {"Num": "1"},
										"children": [
											{
												"tag": "ParagraphNum",
												"attr": {},
												"children": []
											},
											{
												"tag": "ParagraphSentence",
												"attr": {},
												"children": [
													{
														"tag": "Sentence",
														"attr": {},
														"children": ["この法律は、テストを目的とする。"]
													}
												]
											}
										]
									}
								]
							}
						]
					}
				]
			}
		]
	}`)

	result, err := Convert(raw)
	if err != nil {
		t.Fatalf("Convert() error: %v", err)
	}
	if result == "" {
		t.Fatal("Convert() returned empty string")
	}
	if !strings.Contains(result, "第一条") {
		t.Errorf("result missing '第一条', got:\n%s", result)
	}
	if !strings.Contains(result, "この法律は、テストを目的とする。") {
		t.Errorf("result missing sentence, got:\n%s", result)
	}
}

func TestConvertEmpty(t *testing.T) {
	raw := json.RawMessage(`{}`)
	result, err := Convert(raw)
	if err != nil {
		t.Fatalf("Convert() error: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty string, got: %q", result)
	}
}

func TestConvertNull(t *testing.T) {
	raw := json.RawMessage(`null`)
	result, err := Convert(raw)
	if err != nil {
		t.Fatalf("Convert() error: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty string, got: %q", result)
	}
}
