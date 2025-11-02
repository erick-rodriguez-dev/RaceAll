package sharedmemory

import (
	"unsafe"
)

const (
	// Shared memory names
	AccSharedMemoryName   = "Local\\acpmf_physics"
	AccGraphicsMemoryName = "Local\\acpmf_graphics"
	AccStaticMemoryName   = "Local\\acpmf_static"

	// Page file sizes
	PhysicsPageFileSize  = int(unsafe.Sizeof(Physics{}))
	GraphicsPageFileSize = int(unsafe.Sizeof(Graphics{}))
	StaticPageFileSize   = int(unsafe.Sizeof(Static{}))
)

// Status enum
type ACStatus int32

const (
	ACOff ACStatus = iota
	ACReplay
	ACLive
	ACPause
)

// Session Type enum
type ACSessionType int32

const (
	ACUnknown    ACSessionType = -1
	ACPractice   ACSessionType = 0
	ACQualify    ACSessionType = 1
	ACRace       ACSessionType = 2
	ACHotlap     ACSessionType = 3
	ACTimeAttack ACSessionType = 4
	ACDrift      ACSessionType = 5
	ACDrag       ACSessionType = 6
)

// String returns the string representation of the session type
func (s ACSessionType) String() string {
	switch s {
	case ACUnknown:
		return "Unknown"
	case ACPractice:
		return "Practice"
	case ACQualify:
		return "Qualify"
	case ACRace:
		return "Race"
	case ACHotlap:
		return "Hotlap"
	case ACTimeAttack:
		return "TimeAttack"
	case ACDrift:
		return "Drift"
	case ACDrag:
		return "Drag"
	default:
		return "Unknown"
	}
}

// Flag Type enum
type ACFlagType int32

const (
	ACNoFlag      ACFlagType = 0
	ACBlueFlag    ACFlagType = 1
	ACYellowFlag  ACFlagType = 2
	ACBlackFlag   ACFlagType = 3
	ACWhiteFlag   ACFlagType = 4
	ACCheckedFlag ACFlagType = 5
	ACPenaltyFlag ACFlagType = 6
)

// Physics data structure
type Physics struct {
	PacketId            int32
	Gas                 float32
	Brake               float32
	Fuel                float32
	Gear                int32
	Rpms                int32
	SteerAngle          float32
	SpeedKmh            float32
	Velocity            [3]float32
	AccG                [3]float32
	WheelSlip           [4]float32
	WheelLoad           [4]float32
	WheelsPressure      [4]float32
	WheelAngularSpeed   [4]float32
	TyreWear            [4]float32
	TyreDirtyLevel      [4]float32
	TyreCoreTemperature [4]float32
	CamberRAD           [4]float32
	SuspensionTravel    [4]float32
	Drs                 float32
	TC                  float32
	Heading             float32
	Pitch               float32
	Roll                float32
	CgHeight            float32
	CarDamage           [5]float32
	NumberOfTyresOut    int32
	PitLimiterOn        int32
	Abs                 float32
	KersCharge          float32
	KersInput           float32
	AutoShifterOn       int32
	RideHeight          [2]float32
	TurboBoost          float32
	Ballast             float32
	AirDensity          float32
	AirTemp             float32
	RoadTemp            float32
	LocalAngularVel     [3]float32
	FinalFF             float32
	PerformanceMeter    float32
	EngineBrake         int32
	ErsRecoveryLevel    int32
	ErsPowerLevel       int32
	ErsHeatCharging     int32
	ErsIsCharging       int32
	KersCurrentKJ       float32
	DrsAvailable        int32
	DrsEnabled          int32
	BrakeTemp           [4]float32
	Clutch              float32
	TyreTempI           [4]float32
	TyreTempM           [4]float32
	TyreTempO           [4]float32
	IsAIControlled      int32
	TyreContactPoint    [4][3]float32
	TyreContactNormal   [4][3]float32
	TyreContactHeading  [4][3]float32
	BrakeBias           float32
	LocalVelocity       [3]float32
	P2PActivations      int32
	P2PStatus           int32
	CurrentMaxRpm       int32
	Mz                  [4]float32
	Fx                  [4]float32
	Fy                  [4]float32
	SlipRatio           [4]float32
	SlipAngle           [4]float32
	TcinAction          int32
	AbsInAction         int32
	SuspensionDamage    [4]float32
	TyreTemp            [4]float32
	WaterTemp           float32
	BrakePressure       [4]float32
	FrontBrakeCompound  int32
	RearBrakeCompound   int32
	PadLife             [4]float32
	DiscLife            [4]float32
	IgnitionOn          int32
	StarterEngineOn     int32
	IsEngineRunning     int32
	KerbVibration       float32
	SlipVibrations      float32
	GVibrations         float32
	AbsVibrations       float32
}

// Graphics data structure
type Graphics struct {
	PacketId                 int32
	Status                   ACStatus
	Session                  ACSessionType
	CurrentTime              [15]uint16
	LastTime                 [15]uint16
	BestTime                 [15]uint16
	Split                    [15]uint16
	CompletedLaps            int32
	Position                 int32
	ICurrentTime             int32
	ILastTime                int32
	IBestTime                int32
	SessionTimeLeft          float32
	DistanceTraveled         float32
	IsInPit                  int32
	CurrentSectorIndex       int32
	LastSectorTime           int32
	NumberOfLaps             int32
	TyreCompound             [33]uint16
	ReplayTimeMultiplier     float32
	NormalizedCarPosition    float32
	ActiveCars               int32
	CarCoordinates           [60][3]float32
	CarID                    [60]int32
	PlayerCarID              int32
	PenaltyTime              float32
	Flag                     ACFlagType
	PenaltyShortcut          int32
	IdealLineOn              int32
	IsInPitLane              int32
	SurfaceGrip              float32
	MandatoryPitDone         int32
	WindSpeed                float32
	WindDirection            float32
	IsSetupMenuVisible       int32
	MainDisplayIndex         int32
	SecondaryDisplayIndex    int32
	TC                       int32
	TCCUT                    int32
	EngineMap                int32
	ABS                      int32
	FuelXLap                 float32
	RainLights               int32
	FlashingLights           int32
	LightsStage              int32
	ExhaustTemperature       float32
	WiperLV                  int32
	DriverStintTotalTimeLeft int32
	DriverStintTimeLeft      int32
	RainTyres                int32
	SessionIndex             int32
	UsedFuel                 float32
	DeltaLapTime             [15]uint16
	IDeltaLapTime            int32
	EstimatedLapTime         [15]uint16
	IEstimatedLapTime        int32
	IsDeltaPositive          int32
	ISplit                   int32
	IsValidLap               int32
	FuelEstimatedLaps        float32
	TrackStatus              [33]uint16
	MissingMandatoryPits     int32
	Clock                    float32
	DirectionLightsLeft      int32
	DirectionLightsRight     int32
	GlobalYellow             int32
	GlobalYellow1            int32
	GlobalYellow2            int32
	GlobalYellow3            int32
	GlobalWhite              int32
	GlobalGreen              int32
	GlobalChequered          int32
	GlobalRed                int32
	MfdTyreSet               int32
	MfdFuelToAdd             float32
	MfdTyrePressureLF        float32
	MfdTyrePressureRF        float32
	MfdTyrePressureLR        float32
	MfdTyrePressureRR        float32
	TrackGripStatus          int32
	RainIntensity            int32
	RainIntensityIn10min     int32
	RainIntensityIn30min     int32
	CurrentTyreSet           int32
	StrategyTyreSet          int32
}

// Static data structure
type Static struct {
	SMVersion                [15]uint16
	ACVersion                [15]uint16
	NumberOfSessions         int32
	NumCars                  int32
	CarModel                 [33]uint16
	Track                    [33]uint16
	PlayerName               [33]uint16
	PlayerSurname            [33]uint16
	PlayerNick               [33]uint16
	SectorCount              int32
	MaxTorque                float32
	MaxPower                 float32
	MaxRpm                   int32
	MaxFuel                  float32
	SuspensionMaxTravel      [4]float32
	TyreRadius               [4]float32
	MaxTurboBoost            float32
	Deprecated1              float32
	Deprecated2              float32
	PenaltiesEnabled         int32
	AidFuelRate              float32
	AidTireRate              float32
	AidMechanicalDamage      float32
	AidAllowTyreBlankets     int32
	AidStability             float32
	AidAutoClutch            int32
	AidAutoBlip              int32
	HasDRS                   int32
	HasERS                   int32
	HasKERS                  int32
	KersMaxJ                 float32
	EngineBrakeSettingsCount int32
	ErsPowerControllerCount  int32
	TrackSPlineLength        float32
	TrackConfiguration       [33]uint16
	ErsMaxJ                  float32
	IsTimedRace              int32
	HasExtraLap              int32
	CarSkin                  [33]uint16
	ReversedGridPositions    int32
	PitWindowStart           int32
	PitWindowEnd             int32
	IsOnline                 int32
	DryTyresName             [33]uint16
	WetTyresName             [33]uint16
}
