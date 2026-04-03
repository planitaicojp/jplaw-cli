package model

type Sentence struct {
	Position string `json:"position" yaml:"position"`
	Text     string `json:"text" yaml:"text"`
}

type KeywordItem struct {
	LawInfo      LawInfo      `json:"law_info" yaml:"law_info"`
	RevisionInfo RevisionInfo `json:"revision_info" yaml:"revision_info"`
	Sentences    []Sentence   `json:"sentences" yaml:"sentences"`
}

type KeywordResponse struct {
	TotalCount    int64         `json:"total_count" yaml:"total_count"`
	SentenceCount int64         `json:"sentence_count" yaml:"sentence_count"`
	NextOffset    *int64        `json:"next_offset" yaml:"next_offset"`
	Items         []KeywordItem `json:"items" yaml:"items"`
}

type KeywordListRow struct {
	LawNum   string `json:"law_num"`
	LawTitle string `json:"law_title"`
	Match    string `json:"match"`
}
