package model

type RevisionInfo struct {
	LawRevisionID               string                `json:"law_revision_id" yaml:"law_revision_id"`
	LawType                     LawType               `json:"law_type" yaml:"law_type"`
	LawTitle                    string                `json:"law_title" yaml:"law_title"`
	LawTitleKana                string                `json:"law_title_kana" yaml:"law_title_kana"`
	Abbrev                      string                `json:"abbrev" yaml:"abbrev"`
	Category                    string                `json:"category" yaml:"category"`
	Updated                     string                `json:"updated" yaml:"updated"`
	AmendmentPromulgateDate     string                `json:"amendment_promulgate_date" yaml:"amendment_promulgate_date"`
	AmendmentEnforcementDate    string                `json:"amendment_enforcement_date" yaml:"amendment_enforcement_date"`
	AmendmentEnforcementComment string                `json:"amendment_enforcement_comment" yaml:"amendment_enforcement_comment"`
	AmendmentLawID              string                `json:"amendment_law_id" yaml:"amendment_law_id"`
	AmendmentLawTitle           string                `json:"amendment_law_title" yaml:"amendment_law_title"`
	AmendmentLawTitleKana       string                `json:"amendment_law_title_kana" yaml:"amendment_law_title_kana"`
	AmendmentLawNum             string                `json:"amendment_law_num" yaml:"amendment_law_num"`
	AmendmentType               AmendmentType         `json:"amendment_type" yaml:"amendment_type"`
	RepealStatus                RepealStatus          `json:"repeal_status" yaml:"repeal_status"`
	RepealDate                  *string               `json:"repeal_date" yaml:"repeal_date"`
	RemainInForce               bool                  `json:"remain_in_force" yaml:"remain_in_force"`
	Mission                     Mission               `json:"mission" yaml:"mission"`
	CurrentRevisionStatus       CurrentRevisionStatus `json:"current_revision_status" yaml:"current_revision_status"`
}

type RevisionsResponse struct {
	LawInfo   LawInfo        `json:"law_info" yaml:"law_info"`
	Revisions []RevisionInfo `json:"revisions" yaml:"revisions"`
}

type RevisionListRow struct {
	AmendmentDate  string `json:"amendment_date"`
	AmendmentTitle string `json:"amendment_title"`
	AmendmentNum   string `json:"amendment_num"`
	AmendmentType  string `json:"amendment_type"`
}
