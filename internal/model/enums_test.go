package model

import (
	"encoding/json"
	"testing"
)

func TestLawTypeLabel(t *testing.T) {
	tests := []struct {
		input    LawType
		expected string
	}{
		{LawTypeConstitution, "憲法"},
		{LawTypeAct, "法律"},
		{LawTypeCabinetOrder, "政令"},
		{LawTypeImperialOrder, "勅令"},
		{LawTypeMinisterialOrdinance, "府省令"},
		{LawTypeRule, "規則"},
		{LawTypeMisc, "その他"},
		{LawType("Unknown"), "Unknown"},
	}
	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			if got := tt.input.Label(); got != tt.expected {
				t.Errorf("Label() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestEraLabel(t *testing.T) {
	tests := []struct {
		input    Era
		expected string
	}{
		{EraMeiji, "明治"},
		{EraTaisho, "大正"},
		{EraShowa, "昭和"},
		{EraHeisei, "平成"},
		{EraReiwa, "令和"},
	}
	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			if got := tt.input.Label(); got != tt.expected {
				t.Errorf("Label() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestLawTypeJSON(t *testing.T) {
	type wrapper struct {
		Type LawType `json:"type"`
	}
	data := []byte(`{"type":"Act"}`)
	var w wrapper
	if err := json.Unmarshal(data, &w); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if w.Type != LawTypeAct {
		t.Errorf("Type = %q, want %q", w.Type, LawTypeAct)
	}
}

func TestLawTypeFromLabel(t *testing.T) {
	lt := LawTypeFromLabel("法律")
	if lt != LawTypeAct {
		t.Errorf("LawTypeFromLabel('法律') = %q, want %q", lt, LawTypeAct)
	}
	lt = LawTypeFromLabel("unknown")
	if lt != "" {
		t.Errorf("LawTypeFromLabel('unknown') = %q, want empty", lt)
	}
}
