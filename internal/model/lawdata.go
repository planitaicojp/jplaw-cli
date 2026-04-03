package model

import "encoding/json"

type LawDataResponse struct {
	LawInfo           LawInfo            `json:"law_info" yaml:"law_info"`
	RevisionInfo      RevisionInfo       `json:"revision_info" yaml:"revision_info"`
	LawFullText       json.RawMessage    `json:"law_full_text" yaml:"law_full_text"`
	AttachedFilesInfo *AttachedFilesInfo `json:"attached_files_info,omitempty" yaml:"attached_files_info,omitempty"`
}

type AttachedFilesInfo struct {
	ImageData     string         `json:"image_data,omitempty" yaml:"image_data,omitempty"`
	AttachedFiles []AttachedFile `json:"attached_files,omitempty" yaml:"attached_files,omitempty"`
}

type AttachedFile struct {
	LawRevisionID string `json:"law_revision_id" yaml:"law_revision_id"`
	Src           string `json:"src" yaml:"src"`
	Updated       string `json:"updated" yaml:"updated"`
}
