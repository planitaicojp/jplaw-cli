package model

type LawType string

const (
	LawTypeConstitution         LawType = "Constitution"
	LawTypeAct                  LawType = "Act"
	LawTypeCabinetOrder         LawType = "CabinetOrder"
	LawTypeImperialOrder        LawType = "ImperialOrder"
	LawTypeMinisterialOrdinance LawType = "MinisterialOrdinance"
	LawTypeRule                 LawType = "Rule"
	LawTypeMisc                 LawType = "Misc"
)

var lawTypeLabels = map[LawType]string{
	LawTypeConstitution:         "憲法",
	LawTypeAct:                  "法律",
	LawTypeCabinetOrder:         "政令",
	LawTypeImperialOrder:        "勅令",
	LawTypeMinisterialOrdinance: "府省令",
	LawTypeRule:                 "規則",
	LawTypeMisc:                 "その他",
}

var labelToLawType = map[string]LawType{
	"憲法":  LawTypeConstitution,
	"法律":  LawTypeAct,
	"政令":  LawTypeCabinetOrder,
	"勅令":  LawTypeImperialOrder,
	"府省令": LawTypeMinisterialOrdinance,
	"規則":  LawTypeRule,
	"その他": LawTypeMisc,
}

func (lt LawType) Label() string {
	if label, ok := lawTypeLabels[lt]; ok {
		return label
	}
	return string(lt)
}

func LawTypeFromLabel(label string) LawType {
	if lt, ok := labelToLawType[label]; ok {
		return lt
	}
	return ""
}

type Era string

const (
	EraMeiji  Era = "Meiji"
	EraTaisho Era = "Taisho"
	EraShowa  Era = "Showa"
	EraHeisei Era = "Heisei"
	EraReiwa  Era = "Reiwa"
)

var eraLabels = map[Era]string{
	EraMeiji:  "明治",
	EraTaisho: "大正",
	EraShowa:  "昭和",
	EraHeisei: "平成",
	EraReiwa:  "令和",
}

var labelToEra = map[string]Era{
	"明治": EraMeiji,
	"大正": EraTaisho,
	"昭和": EraShowa,
	"平成": EraHeisei,
	"令和": EraReiwa,
}

func (e Era) Label() string {
	if label, ok := eraLabels[e]; ok {
		return label
	}
	return string(e)
}

func EraFromLabel(label string) Era {
	if e, ok := labelToEra[label]; ok {
		return e
	}
	return ""
}

type RepealStatus string

const (
	RepealStatusNone                RepealStatus = "None"
	RepealStatusRepeal              RepealStatus = "Repeal"
	RepealStatusExpire              RepealStatus = "Expire"
	RepealStatusSuspend             RepealStatus = "Suspend"
	RepealStatusLossOfEffectiveness RepealStatus = "LossOfEffectiveness"
)

type AmendmentType string

const (
	AmendmentTypeNew     AmendmentType = "1"
	AmendmentTypeAmended AmendmentType = "3"
	AmendmentTypeRepeal  AmendmentType = "8"
)

type Mission string

const (
	MissionNew     Mission = "New"
	MissionPartial Mission = "Partial"
)

type CurrentRevisionStatus string

const (
	CurrentRevisionStatusCurrentEnforced  CurrentRevisionStatus = "CurrentEnforced"
	CurrentRevisionStatusUnEnforced       CurrentRevisionStatus = "UnEnforced"
	CurrentRevisionStatusPreviousEnforced CurrentRevisionStatus = "PreviousEnforced"
	CurrentRevisionStatusRepeal           CurrentRevisionStatus = "Repeal"
)
