package model

type LawInfo struct {
	LawType      LawType `json:"law_type" yaml:"law_type"`
	LawID        string  `json:"law_id" yaml:"law_id"`
	LawNum       string  `json:"law_num" yaml:"law_num"`
	LawNumEra    Era     `json:"law_num_era" yaml:"law_num_era"`
	LawNumYear   int     `json:"law_num_year" yaml:"law_num_year"`
	LawNumType   LawType `json:"law_num_type" yaml:"law_num_type"`
	LawNumNum    string  `json:"law_num_num" yaml:"law_num_num"`
	Promulgation string  `json:"promulgation_date" yaml:"promulgation_date"`
}

type LawItem struct {
	LawInfo             LawInfo      `json:"law_info" yaml:"law_info"`
	RevisionInfo        RevisionInfo `json:"revision_info" yaml:"revision_info"`
	CurrentRevisionInfo RevisionInfo `json:"current_revision_info" yaml:"current_revision_info"`
}

type LawsResponse struct {
	TotalCount int64     `json:"total_count" yaml:"total_count"`
	Count      int64     `json:"count" yaml:"count"`
	NextOffset *int64    `json:"next_offset" yaml:"next_offset"`
	Laws       []LawItem `json:"laws" yaml:"laws"`
}

type LawListRow struct {
	LawID        string `json:"law_id" header:"法令ID"`
	LawNum       string `json:"law_num" header:"法令番号"`
	LawTitle     string `json:"law_title" header:"法令名"`
	LawType      string `json:"law_type" header:"法令種別"`
	Promulgation string `json:"promulgation_date" header:"公布日"`
}
