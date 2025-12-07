package structs

// DriverCategory represents the driver skill category
type DriverCategory int

const (
	Bronze   DriverCategory = 0
	Silver   DriverCategory = 1
	Gold     DriverCategory = 2
	Platinum DriverCategory = 3
	Error    DriverCategory = 255
)

// LapType represents the type of lap
type LapType int

const (
	LapTypeError LapType = 0
	Outlap       LapType = 1
	Regular      LapType = 2
	Inlap        LapType = 3
)

// CarLocationEnum represents the location of a car on track
type CarLocationEnum int

const (
	LocationNone CarLocationEnum = 0
	Track        CarLocationEnum = 1
	Pitlane      CarLocationEnum = 2
	PitEntry     CarLocationEnum = 3
	PitExit      CarLocationEnum = 4
)

// SessionPhase represents the current phase of the session
type SessionPhase int

const (
	PhaseNone    SessionPhase = 0
	Starting     SessionPhase = 1
	PreFormation SessionPhase = 2
	FormationLap SessionPhase = 3
	PreSession   SessionPhase = 4
	Session      SessionPhase = 5
	SessionOver  SessionPhase = 6
	PostSession  SessionPhase = 7
	ResultUI     SessionPhase = 8
)

// RaceSessionType represents the type of racing session
type RaceSessionType int

const (
	Practice        RaceSessionType = 0
	Qualifying      RaceSessionType = 4
	Superpole       RaceSessionType = 9
	Race            RaceSessionType = 10
	Hotlap          RaceSessionType = 11
	Hotstint        RaceSessionType = 12
	HotlapSuperpole RaceSessionType = 13
	Replay          RaceSessionType = 14
)

// BroadcastingCarEventType represents types of car events
type BroadcastingCarEventType int

const (
	EventNone        BroadcastingCarEventType = 0
	GreenFlag        BroadcastingCarEventType = 1
	EventSessionOver BroadcastingCarEventType = 2
	PenaltyCommMsg   BroadcastingCarEventType = 3
	Accident         BroadcastingCarEventType = 4
	LapCompleted     BroadcastingCarEventType = 5
	BestSessionLap   BroadcastingCarEventType = 6
	BestPersonalLap  BroadcastingCarEventType = 7
)

// NationalityEnum represents driver/team nationality
type NationalityEnum int

const (
	Any             NationalityEnum = 0
	Italy           NationalityEnum = 1
	Germany         NationalityEnum = 2
	France          NationalityEnum = 3
	Spain           NationalityEnum = 4
	GreatBritain    NationalityEnum = 5
	Hungary         NationalityEnum = 6
	Belgium         NationalityEnum = 7
	Switzerland     NationalityEnum = 8
	Austria         NationalityEnum = 9
	Russia          NationalityEnum = 10
	Thailand        NationalityEnum = 11
	Netherlands     NationalityEnum = 12
	Poland          NationalityEnum = 13
	Argentina       NationalityEnum = 14
	Monaco          NationalityEnum = 15
	Ireland         NationalityEnum = 16
	Brazil          NationalityEnum = 17
	SouthAfrica     NationalityEnum = 18
	PuertoRico      NationalityEnum = 19
	Slovakia        NationalityEnum = 20
	Oman            NationalityEnum = 21
	Greece          NationalityEnum = 22
	SaudiArabia     NationalityEnum = 23
	Norway          NationalityEnum = 24
	Turkey          NationalityEnum = 25
	SouthKorea      NationalityEnum = 26
	Lebanon         NationalityEnum = 27
	Armenia         NationalityEnum = 28
	Mexico          NationalityEnum = 29
	Sweden          NationalityEnum = 30
	Finland         NationalityEnum = 31
	Denmark         NationalityEnum = 32
	Croatia         NationalityEnum = 33
	Canada          NationalityEnum = 34
	China           NationalityEnum = 35
	Portugal        NationalityEnum = 36
	Singapore       NationalityEnum = 37
	Indonesia       NationalityEnum = 38
	USA             NationalityEnum = 39
	NewZealand      NationalityEnum = 40
	Australia       NationalityEnum = 41
	SanMarino       NationalityEnum = 42
	UAE             NationalityEnum = 43
	Luxembourg      NationalityEnum = 44
	Kuwait          NationalityEnum = 45
	HongKong        NationalityEnum = 46
	Colombia        NationalityEnum = 47
	Japan           NationalityEnum = 48
	Andorra         NationalityEnum = 49
	Azerbaijan      NationalityEnum = 50
	Bulgaria        NationalityEnum = 51
	Cuba            NationalityEnum = 52
	CzechRepublic   NationalityEnum = 53
	Estonia         NationalityEnum = 54
	Georgia         NationalityEnum = 55
	India           NationalityEnum = 56
	Israel          NationalityEnum = 57
	Jamaica         NationalityEnum = 58
	Latvia          NationalityEnum = 59
	Lithuania       NationalityEnum = 60
	Macau           NationalityEnum = 61
	Malaysia        NationalityEnum = 62
	Nepal           NationalityEnum = 63
	NewCaledonia    NationalityEnum = 64
	Nigeria         NationalityEnum = 65
	NorthernIreland NationalityEnum = 66
	PapuaNewGuinea  NationalityEnum = 67
	Philippines     NationalityEnum = 68
	Qatar           NationalityEnum = 69
	Romania         NationalityEnum = 70
	Scotland        NationalityEnum = 71
	Serbia          NationalityEnum = 72
	Slovenia        NationalityEnum = 73
	Taiwan          NationalityEnum = 74
	Ukraine         NationalityEnum = 75
	Venezuela       NationalityEnum = 76
	Wales           NationalityEnum = 77
	Iran            NationalityEnum = 78
	Bahrain         NationalityEnum = 79
	Zimbabwe        NationalityEnum = 80
	ChineseTaipei   NationalityEnum = 81
	Chile           NationalityEnum = 82
	Uruguay         NationalityEnum = 83
	Madagascar      NationalityEnum = 84
)
