package broadcast

type DriverCategory byte

const (
	DriverCategoryBronze   DriverCategory = 0
	DriverCategorySilver   DriverCategory = 1
	DriverCategoryGold     DriverCategory = 2
	DriverCategoryPlatinum DriverCategory = 3
	DriverCategoryError    DriverCategory = 255
)

var driverCategoryNames = map[DriverCategory]string{
	DriverCategoryBronze:   "Bronze",
	DriverCategorySilver:   "Silver",
	DriverCategoryGold:     "Gold",
	DriverCategoryPlatinum: "Platinum",
	DriverCategoryError:    "Error",
}

func (d DriverCategory) String() string {
	if name, ok := driverCategoryNames[d]; ok {
		return name
	}
	return "Error"
}

type LapType byte

const (
	LapTypeError   LapType = 0
	LapTypeOutlap  LapType = 1
	LapTypeRegular LapType = 2
	LapTypeInlap   LapType = 3
)

var lapTypeNames = map[LapType]string{
	LapTypeError:   "Error",
	LapTypeOutlap:  "Outlap",
	LapTypeRegular: "Regular",
	LapTypeInlap:   "Inlap",
}

func (l LapType) String() string {
	if name, ok := lapTypeNames[l]; ok {
		return name
	}
	return "Error"
}

type CarLocationEnum byte

const (
	CarLocationNone     CarLocationEnum = 0
	CarLocationTrack    CarLocationEnum = 1
	CarLocationPitlane  CarLocationEnum = 2
	CarLocationPitEntry CarLocationEnum = 3
	CarLocationPitExit  CarLocationEnum = 4
)

var carLocationNames = map[CarLocationEnum]string{
	CarLocationNone:     "None",
	CarLocationTrack:    "Track",
	CarLocationPitlane:  "Pitlane",
	CarLocationPitEntry: "Pit Entry",
	CarLocationPitExit:  "Pit Exit",
}

func (c CarLocationEnum) String() string {
	if name, ok := carLocationNames[c]; ok {
		return name
	}
	return "None"
}

type SessionPhase byte

const (
	SessionPhaseNone         SessionPhase = 0
	SessionPhaseStarting     SessionPhase = 1
	SessionPhasePreFormation SessionPhase = 2
	SessionPhaseFormationLap SessionPhase = 3
	SessionPhasePreSession   SessionPhase = 4
	SessionPhaseSession      SessionPhase = 5
	SessionPhaseSessionOver  SessionPhase = 6
	SessionPhasePostSession  SessionPhase = 7
	SessionPhaseResultUI     SessionPhase = 8
)

var sessionPhaseNames = map[SessionPhase]string{
	SessionPhaseNone:         "None",
	SessionPhaseStarting:     "Starting",
	SessionPhasePreFormation: "Pre-Formation",
	SessionPhaseFormationLap: "Formation Lap",
	SessionPhasePreSession:   "Pre-Session",
	SessionPhaseSession:      "Session",
	SessionPhaseSessionOver:  "Session Over",
	SessionPhasePostSession:  "Post-Session",
	SessionPhaseResultUI:     "Result UI",
}

func (s SessionPhase) String() string {
	if name, ok := sessionPhaseNames[s]; ok {
		return name
	}
	return "None"
}

type RaceSessionType byte

const (
	RaceSessionTypePractice        RaceSessionType = 0
	RaceSessionTypeQualifying      RaceSessionType = 4
	RaceSessionTypeSuperpole       RaceSessionType = 9
	RaceSessionTypeRace            RaceSessionType = 10
	RaceSessionTypeHotlap          RaceSessionType = 11
	RaceSessionTypeHotstint        RaceSessionType = 12
	RaceSessionTypeHotlapSuperpole RaceSessionType = 13
	RaceSessionTypeReplay          RaceSessionType = 14
)

var raceSessionTypeNames = map[RaceSessionType]string{
	RaceSessionTypePractice:        "Practice",
	RaceSessionTypeQualifying:      "Qualifying",
	RaceSessionTypeSuperpole:       "Superpole",
	RaceSessionTypeRace:            "Race",
	RaceSessionTypeHotlap:          "Hotlap",
	RaceSessionTypeHotstint:        "Hotstint",
	RaceSessionTypeHotlapSuperpole: "Hotlap Superpole",
	RaceSessionTypeReplay:          "Replay",
}

func (r RaceSessionType) String() string {
	if name, ok := raceSessionTypeNames[r]; ok {
		return name
	}
	return "Unknown"
}

type BroadcastingEventType byte

const (
	BroadcastingEventTypeNone            BroadcastingEventType = 0
	BroadcastingEventTypeGreenFlag       BroadcastingEventType = 1
	BroadcastingEventTypeSessionOver     BroadcastingEventType = 2
	BroadcastingEventTypePenaltyCommMsg  BroadcastingEventType = 3
	BroadcastingEventTypeAccident        BroadcastingEventType = 4
	BroadcastingEventTypeLapCompleted    BroadcastingEventType = 5
	BroadcastingEventTypeBestSessionLap  BroadcastingEventType = 6
	BroadcastingEventTypeBestPersonalLap BroadcastingEventType = 7
)

var broadcastingEventTypeNames = map[BroadcastingEventType]string{
	BroadcastingEventTypeNone:            "None",
	BroadcastingEventTypeGreenFlag:       "Green Flag",
	BroadcastingEventTypeSessionOver:     "Session Over",
	BroadcastingEventTypePenaltyCommMsg:  "Penalty Communication",
	BroadcastingEventTypeAccident:        "Accident",
	BroadcastingEventTypeLapCompleted:    "Lap Completed",
	BroadcastingEventTypeBestSessionLap:  "Best Session Lap",
	BroadcastingEventTypeBestPersonalLap: "Best Personal Lap",
}

func (b BroadcastingEventType) String() string {
	if name, ok := broadcastingEventTypeNames[b]; ok {
		return name
	}
	return "None"
}

type NationalityEnum uint16

const (
	NationalityAny             NationalityEnum = 0
	NationalityItaly           NationalityEnum = 1
	NationalityGermany         NationalityEnum = 2
	NationalityFrance          NationalityEnum = 3
	NationalitySpain           NationalityEnum = 4
	NationalityGreatBritain    NationalityEnum = 5
	NationalityHungary         NationalityEnum = 6
	NationalityBelgium         NationalityEnum = 7
	NationalitySwitzerland     NationalityEnum = 8
	NationalityAustria         NationalityEnum = 9
	NationalityRussia          NationalityEnum = 10
	NationalityThailand        NationalityEnum = 11
	NationalityNetherlands     NationalityEnum = 12
	NationalityPoland          NationalityEnum = 13
	NationalityArgentina       NationalityEnum = 14
	NationalityMonaco          NationalityEnum = 15
	NationalityIreland         NationalityEnum = 16
	NationalityBrazil          NationalityEnum = 17
	NationalitySouthAfrica     NationalityEnum = 18
	NationalityPuertoRico      NationalityEnum = 19
	NationalitySlovakia        NationalityEnum = 20
	NationalityOman            NationalityEnum = 21
	NationalityGreece          NationalityEnum = 22
	NationalitySaudiArabia     NationalityEnum = 23
	NationalityNorway          NationalityEnum = 24
	NationalityTurkey          NationalityEnum = 25
	NationalitySouthKorea      NationalityEnum = 26
	NationalityLebanon         NationalityEnum = 27
	NationalityArmenia         NationalityEnum = 28
	NationalityMexico          NationalityEnum = 29
	NationalitySweden          NationalityEnum = 30
	NationalityFinland         NationalityEnum = 31
	NationalityDenmark         NationalityEnum = 32
	NationalityCroatia         NationalityEnum = 33
	NationalityCanada          NationalityEnum = 34
	NationalityChina           NationalityEnum = 35
	NationalityPortugal        NationalityEnum = 36
	NationalitySingapore       NationalityEnum = 37
	NationalityIndonesia       NationalityEnum = 38
	NationalityUSA             NationalityEnum = 39
	NationalityNewZealand      NationalityEnum = 40
	NationalityAustralia       NationalityEnum = 41
	NationalitySanMarino       NationalityEnum = 42
	NationalityUAE             NationalityEnum = 43
	NationalityLuxembourg      NationalityEnum = 44
	NationalityKuwait          NationalityEnum = 45
	NationalityHongKong        NationalityEnum = 46
	NationalityColombia        NationalityEnum = 47
	NationalityJapan           NationalityEnum = 48
	NationalityAndorra         NationalityEnum = 49
	NationalityAzerbaijan      NationalityEnum = 50
	NationalityBulgaria        NationalityEnum = 51
	NationalityCuba            NationalityEnum = 52
	NationalityCzechRepublic   NationalityEnum = 53
	NationalityEstonia         NationalityEnum = 54
	NationalityGeorgia         NationalityEnum = 55
	NationalityIndia           NationalityEnum = 56
	NationalityIsrael          NationalityEnum = 57
	NationalityJamaica         NationalityEnum = 58
	NationalityLatvia          NationalityEnum = 59
	NationalityLithuania       NationalityEnum = 60
	NationalityMacau           NationalityEnum = 61
	NationalityMalaysia        NationalityEnum = 62
	NationalityNepal           NationalityEnum = 63
	NationalityNewCaledonia    NationalityEnum = 64
	NationalityNigeria         NationalityEnum = 65
	NationalityNorthernIreland NationalityEnum = 66
	NationalityPapuaNewGuinea  NationalityEnum = 67
	NationalityPhilippines     NationalityEnum = 68
	NationalityQatar           NationalityEnum = 69
	NationalityRomania         NationalityEnum = 70
	NationalityScotland        NationalityEnum = 71
	NationalitySerbia          NationalityEnum = 72
	NationalitySlovenia        NationalityEnum = 73
	NationalityTaiwan          NationalityEnum = 74
	NationalityUkraine         NationalityEnum = 75
	NationalityVenezuela       NationalityEnum = 76
	NationalityWales           NationalityEnum = 77
	NationalityIran            NationalityEnum = 78
	NationalityBahrain         NationalityEnum = 79
	NationalityZimbabwe        NationalityEnum = 80
	NationalityChineseTaipei   NationalityEnum = 81
	NationalityChile           NationalityEnum = 82
	NationalityUruguay         NationalityEnum = 83
	NationalityMadagascar      NationalityEnum = 84
)

var nationalityNames = map[NationalityEnum]string{
	NationalityAny:             "Any",
	NationalityItaly:           "Italy",
	NationalityGermany:         "Germany",
	NationalityFrance:          "France",
	NationalitySpain:           "Spain",
	NationalityGreatBritain:    "Great Britain",
	NationalityHungary:         "Hungary",
	NationalityBelgium:         "Belgium",
	NationalitySwitzerland:     "Switzerland",
	NationalityAustria:         "Austria",
	NationalityRussia:          "Russia",
	NationalityThailand:        "Thailand",
	NationalityNetherlands:     "Netherlands",
	NationalityPoland:          "Poland",
	NationalityArgentina:       "Argentina",
	NationalityMonaco:          "Monaco",
	NationalityIreland:         "Ireland",
	NationalityBrazil:          "Brazil",
	NationalitySouthAfrica:     "South Africa",
	NationalityPuertoRico:      "Puerto Rico",
	NationalitySlovakia:        "Slovakia",
	NationalityOman:            "Oman",
	NationalityGreece:          "Greece",
	NationalitySaudiArabia:     "Saudi Arabia",
	NationalityNorway:          "Norway",
	NationalityTurkey:          "Turkey",
	NationalitySouthKorea:      "South Korea",
	NationalityLebanon:         "Lebanon",
	NationalityArmenia:         "Armenia",
	NationalityMexico:          "Mexico",
	NationalitySweden:          "Sweden",
	NationalityFinland:         "Finland",
	NationalityDenmark:         "Denmark",
	NationalityCroatia:         "Croatia",
	NationalityCanada:          "Canada",
	NationalityChina:           "China",
	NationalityPortugal:        "Portugal",
	NationalitySingapore:       "Singapore",
	NationalityIndonesia:       "Indonesia",
	NationalityUSA:             "USA",
	NationalityNewZealand:      "New Zealand",
	NationalityAustralia:       "Australia",
	NationalitySanMarino:       "San Marino",
	NationalityUAE:             "UAE",
	NationalityLuxembourg:      "Luxembourg",
	NationalityKuwait:          "Kuwait",
	NationalityHongKong:        "Hong Kong",
	NationalityColombia:        "Colombia",
	NationalityJapan:           "Japan",
	NationalityAndorra:         "Andorra",
	NationalityAzerbaijan:      "Azerbaijan",
	NationalityBulgaria:        "Bulgaria",
	NationalityCuba:            "Cuba",
	NationalityCzechRepublic:   "Czech Republic",
	NationalityEstonia:         "Estonia",
	NationalityGeorgia:         "Georgia",
	NationalityIndia:           "India",
	NationalityIsrael:          "Israel",
	NationalityJamaica:         "Jamaica",
	NationalityLatvia:          "Latvia",
	NationalityLithuania:       "Lithuania",
	NationalityMacau:           "Macau",
	NationalityMalaysia:        "Malaysia",
	NationalityNepal:           "Nepal",
	NationalityNewCaledonia:    "New Caledonia",
	NationalityNigeria:         "Nigeria",
	NationalityNorthernIreland: "Northern Ireland",
	NationalityPapuaNewGuinea:  "Papua New Guinea",
	NationalityPhilippines:     "Philippines",
	NationalityQatar:           "Qatar",
	NationalityRomania:         "Romania",
	NationalityScotland:        "Scotland",
	NationalitySerbia:          "Serbia",
	NationalitySlovenia:        "Slovenia",
	NationalityTaiwan:          "Taiwan",
	NationalityUkraine:         "Ukraine",
	NationalityVenezuela:       "Venezuela",
	NationalityWales:           "Wales",
	NationalityIran:            "Iran",
	NationalityBahrain:         "Bahrain",
	NationalityZimbabwe:        "Zimbabwe",
	NationalityChineseTaipei:   "Chinese Taipei",
	NationalityChile:           "Chile",
	NationalityUruguay:         "Uruguay",
	NationalityMadagascar:      "Madagascar",
}

func (n NationalityEnum) String() string {
	if name, ok := nationalityNames[n]; ok {
		return name
	}
	return "Any"
}
