package cmdutil

import "github.com/planitaicojp/jplaw-cli/internal/model"

// ParseLawTypes converts Japanese or API law type names to model.LawType values.
func ParseLawTypes(inputs []string) []model.LawType {
	var result []model.LawType
	for _, lt := range inputs {
		apiType := model.LawTypeFromLabel(lt)
		if apiType == "" {
			result = append(result, model.LawType(lt))
		} else {
			result = append(result, apiType)
		}
	}
	return result
}

// ParseEra converts a Japanese or API era name to model.Era.
func ParseEra(input string) model.Era {
	if input == "" {
		return ""
	}
	era := model.EraFromLabel(input)
	if era == "" {
		return model.Era(input)
	}
	return era
}
