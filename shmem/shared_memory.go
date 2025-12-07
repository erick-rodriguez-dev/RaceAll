package shared_memory

import (
	"fmt"
	"sync"
)

const (
	FILE_MAP_READ = 0x0004
)

// AccStatus represents the game status
type AccStatus int32

const (
	AC_OFF AccStatus = iota
	AC_REPLAY
	AC_LIVE
	AC_PAUSE
)

// AccSessionType represents the session type
type AccSessionType int32

const (
	AC_UNKNOWN         AccSessionType = -1
	AC_PRACTICE        AccSessionType = 0
	AC_QUALIFY         AccSessionType = 1
	AC_RACE            AccSessionType = 2
	AC_HOTLAP          AccSessionType = 3
	AC_TIME_ATTACK     AccSessionType = 4
	AC_DRIFT           AccSessionType = 5
	AC_DRAG            AccSessionType = 6
	AC_HOTSTINT        AccSessionType = 7
	AC_HOTLAPSUPERPOLE AccSessionType = 8
)

// SessionTypeToString converts session type to string
func SessionTypeToString(sessionType AccSessionType) string {
	switch sessionType {
	case AC_UNKNOWN:
		return "Unknown"
	case AC_PRACTICE:
		return "Practice"
	case AC_QUALIFY:
		return "Qualify"
	case AC_RACE:
		return "Race"
	case AC_HOTLAP:
		return "Hotlap"
	case AC_TIME_ATTACK:
		return "Time attack"
	case AC_DRIFT:
		return "Drift"
	case AC_DRAG:
		return "Drag"
	case AC_HOTSTINT:
		return "Hotstint"
	case AC_HOTLAPSUPERPOLE:
		return "Hotlap superpole"
	default:
		return fmt.Sprintf("%d", sessionType)
	}
}

// AccFlagType represents flag types
type AccFlagType int32

const (
	AC_NO_FLAG AccFlagType = iota
	AC_BLUE_FLAG
	AC_YELLOW_FLAG
	AC_BLACK_FLAG
	AC_WHITE_FLAG
	AC_CHECKERED_FLAG
	AC_PENALTY_FLAG
	AC_GREEN_FLAG
	AC_BLACK_FLAG_WITH_ORANGE_CIRCLE
)

// FlagTypeToString converts flag type to string
func FlagTypeToString(flagType AccFlagType) string {
	switch flagType {
	case AC_NO_FLAG, AC_GREEN_FLAG:
		return "Green"
	case AC_BLUE_FLAG:
		return "Blue"
	case AC_YELLOW_FLAG:
		return "Yellow"
	case AC_BLACK_FLAG:
		return "Black"
	case AC_WHITE_FLAG:
		return "White"
	case AC_CHECKERED_FLAG:
		return "Checkered"
	case AC_PENALTY_FLAG:
		return "Penalty"
	case AC_BLACK_FLAG_WITH_ORANGE_CIRCLE:
		return "Orange"
	default:
		return fmt.Sprintf("%d", flagType)
	}
}

// PenaltyShortcut represents penalty types
type PenaltyShortcut int32

const (
	None PenaltyShortcut = iota
	DriveThrough_Cutting
	StopAndGo_10_Cutting
	StopAndGo_20_Cutting
	StopAndGo_30_Cutting
	Disqualified_Cutting
	RemoveBestLaptime_Cutting
	DriveThrough_PitSpeeding
	StopAndGo_10_PitSpeeding
	StopAndGo_20_PitSpeeding
	StopAndGo_30_PitSpeeding
	Disqualified_PitSpeeding
	RemoveBestLaptime_PitSpeeding
	Disqualified_IgnoredMandatoryPit
	PostRaceTime
	Disqualified_Trolling
	Disqualified_PitEntry
	Disqualified_PitExit
	Disqualified_WrongWay
	DriveThrough_IgnoredDriverStint
	Disqualified_IgnoredDriverStint
	Disqualified_ExceededDriverStintLimit
)

// AccTrackGripStatus represents track grip status
type AccTrackGripStatus int32

const (
	Green AccTrackGripStatus = iota
	Fast
	Optimum
	Greasy
	Damp
	Wet
	Flooded
)

// AccRainIntensity represents rain intensity
type AccRainIntensity int32

const (
	No_Rain AccRainIntensity = iota
	Dew
	Light_Rain
	Medium_Rain
	Heavy_Rain
	Thunderstorm
)

// AccRainIntensityToString converts rain intensity to string
func AccRainIntensityToString(intensity AccRainIntensity) string {
	switch intensity {
	case No_Rain:
		return "Dry"
	case Dew:
		return "Dew"
	case Light_Rain:
		return "Light"
	case Medium_Rain:
		return "Medium"
	case Heavy_Rain:
		return "Heavy"
	case Thunderstorm:
		return "Thunder"
	default:
		return ""
	}
}

// StructVector3 represents a 3D vector
type StructVector3 struct {
	X float32
	Y float32
	Z float32
}

func (v StructVector3) String() string {
	return fmt.Sprintf("X: %f, Y: %f, Z: %f", v.X, v.Y, v.Z)
}

// SPageFileGraphic contains graphical step data
type SPageFileGraphic struct {
	PacketId                 int32
	Status                   AccStatus
	SessionType              AccSessionType
	CurrentTime              [15]uint16
	LastTime                 [15]uint16
	BestTime                 [15]uint16
	Split                    [15]uint16
	CompletedLaps            int32
	Position                 int32
	CurrentTimeMs            int32
	LastTimeMs               int32
	BestTimeMs               int32
	SessionTimeLeft          float32
	DistanceTraveled         float32
	IsInPits                 int32
	CurrentSectorIndex       int32
	LastSectorTime           int32
	NumberOfLaps             int32
	TyreCompound             [33]uint16
	ReplayTimeMultiplier     float32 // Not used in ACC
	NormalizedCarPosition    float32
	ActiveCars               int32
	CarCoordinates           [60]StructVector3
	CarIds                   [60]int32
	PlayerCarID              int32
	PenaltyTime              float32
	Flag                     AccFlagType
	PenaltyType              PenaltyShortcut
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
	TCCut                    int32
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
	UsedFuelSinceRefuel      float32
	DeltaLapTime             [15]uint16
	DeltaLapTimeMillis       int32
	EstimatedLapTime         [15]uint16
	EstimatedLapTimeMillis   int32
	IsDeltaPositive          int32
	SplitTimeMillis          int32
	IsValidLap               int32
	FuelEstimatedLaps        float32
	TrackStatus              [33]uint16
	MandatoryPitStopsLeft    int32
	ClockTimeDaySeconds      float32
	BlinkerLeftOn            int32
	BlinkerRightOn           int32
	GlobalYellow             int32
	GlobalYellowSector1      int32
	GlobalYellowSector2      int32
	GlobalYellowSector3      int32
	GlobalWhite              int32
	GreenFlag                int32
	GlobalChequered          int32
	GlobalRed                int32
	MfdTyreSet               int32
	MfdFuelToAdd             float32
	MfdTyrePressureLF        float32
	MfdTyrePressureRF        float32
	MfdTyrePressureLR        float32
	MfdTyrePressureRR        float32
	TrackGripStatus          AccTrackGripStatus
	RainIntensity            AccRainIntensity
	RainIntensityIn10min     AccRainIntensity
	RainIntensityIn30min     AccRainIntensity
	CurrentTyreSet           int32
	StrategyTyreSet          int32
	GapAheadMillis           int32
	GapBehindMillis          int32
}

// SPageFilePhysics contains physics data
type SPageFilePhysics struct {
	PacketId             int32
	Gas                  float32
	Brake                float32
	Fuel                 float32
	Gear                 int32
	Rpms                 int32
	SteerAngle           float32
	SpeedKmh             float32
	Velocity             [3]float32
	AccG                 [3]float32
	WheelSlip            [4]float32
	WheelLoad            [4]float32 // Not used in ACC
	WheelPressure        [4]float32
	WheelAngularSpeed    [4]float32
	TyreWear             [4]float32 // Not used in ACC
	TyreDirtyLevel       [4]float32 // Not used in ACC
	TyreCoreTemperature  [4]float32
	CamberRad            [4]float32 // Not used in ACC
	SuspensionTravel     [4]float32
	Drs                  float32 // Not used in ACC
	TC                   float32
	Heading              float32
	Pitch                float32
	Roll                 float32 // Not used in ACC
	CgHeight             float32
	CarDamage            [5]float32
	NumberOfTyresOut     int32 // Not used in ACC
	PitLimiterOn         int32
	Abs                  float32
	KersCharge           float32 // Not used in ACC
	KersInput            float32 // Not used in ACC
	AutoShifterOn        int32
	RideHeight           [2]float32 // Not used in ACC
	TurboBoost           float32
	Ballast              float32 // Not used in ACC
	AirDensity           float32 // Not used in ACC
	AirTemp              float32
	RoadTemp             float32
	LocalAngularVelocity [3]float32
	FinalFF              float32
	PerformanceMeter     float32 // Not used in ACC
	EngineBrake          int32   // Not used in ACC
	ErsRecoveryLevel     int32   // Not used in ACC
	ErsPowerLevel        int32   // Not used in ACC
	ErsHeatCharging      int32   // Not used in ACC
	ErsIsCharging        int32   // Not used in ACC
	KersCurrentKJ        float32 // Not used in ACC
	DrsAvailable         int32   // Not used in ACC
	DrsEnabled           int32   // Not used in ACC
	BrakeTemperature     [4]float32
	Clutch               float32
	TyreTempI            [4]float32 // Not shown in ACC
	TyreTempM            [4]float32 // Not shown in ACC
	TyreTempO            [4]float32 // Not shown in ACC
	IsAiControlled       int32
	TyreContactPoint     [4]StructVector3
	TyreContactNormal    [4]StructVector3
	TyreContactHeading   [4]StructVector3
	BrakeBias            float32
	LocalVelocity        [3]float32
	P2PActivations       int32      // Not used in ACC
	P2PStatus            int32      // Not used in ACC
	CurrentMaxRpm        int32      // Maximum engine rpm
	Mz                   [4]float32 // Not shown in ACC
	Fx                   [4]float32 // Not shown in ACC
	Fy                   [4]float32 // Not shown in ACC
	SlipRatio            [4]float32
	SlipAngle            [4]float32
	TcinAction           int32      // Not used in ACC
	AbsInAction          int32      // Not used in ACC
	SuspensionDamage     [4]float32 // Not used in ACC
	TyreTemp             [4]float32 // Not used in ACC
	WaterTemp            float32
	BrakePressure        [4]float32
	FrontBrakeCompound   int32
	RearBrakeCompound    int32
	PadLife              [4]float32
	DiscLife             [4]float32
	IgnitionOn           int32
	StarterEngineOn      int32
	IsEngineRunning      int32
	KerbVibration        float32
	SlipVibrations       float32
	Gvibrations          float32
	AbsVibrations        float32
}

// SPageFileStatic contains static session data
type SPageFileStatic struct {
	SharedMemoryVersion      [15]uint16
	AssettoCorsaVersion      [15]uint16
	NumberOfSessions         int32
	NumberOfCars             int32
	CarModel                 [33]uint16
	Track                    [33]uint16
	PlayerName               [33]uint16
	PlayerSurname            [33]uint16
	PlayerNickname           [33]uint16
	SectorCount              int32
	MaxTorque                float32 // Not shown in ACC
	MaxPower                 float32 // Not shown in ACC
	MaxRpm                   int32
	MaxFuel                  float32
	SuspensionMaxTravel      [4]float32 // Not shown in ACC
	TyreRadius               [4]float32 // Not shown in ACC
	MaxTurboBoost            float32    // Not used in ACC
	AirTemperature           float32    // Not used in ACC
	RoadTemperature          float32    // Not used in ACC
	PenaltiesEnabled         int32
	AidFuelRate              float32
	AidTireRate              float32
	AidMechanicalDamage      float32
	AidAllowTyreBlankets     int32
	AidStability             float32
	AidAutoClutch            int32
	AidAutoBlip              int32
	HasDRS                   int32      // Not used in ACC
	HasERS                   int32      // Not used in ACC
	HasKERS                  int32      // Not used in ACC
	KersMaxJoules            float32    // Not used in ACC
	EngineBrakeSettingsCount int32      // Not used in ACC
	ErsPowerControllerCount  int32      // Not used in ACC
	TrackSplineLength        float32    // Not used in ACC
	TrackConfiguration       [33]uint16 // Not used in ACC
	ErsMaxJ                  float32    // Not used in ACC
	IsTimedRace              int32      // Not used in ACC
	HasExtraLap              int32      // Not used in ACC
	CarSkin                  [33]uint16 // Not used in ACC
	ReversedGridPositions    int32      // Not used in ACC
	PitWindowStart           int32
	PitWindowEnd             int32
	IsOnline                 int32
	DryTyresName             [33]uint16
	WetTyresName             [33]uint16
}

// ACCSharedMemory manages access to ACC shared memory
type ACCSharedMemory struct {
	physicsMap  string
	graphicsMap string
	staticMap   string

	PageFileStatic  *SPageFileStatic
	PageFilePhysics *SPageFilePhysics
	PageFileGraphic *SPageFileGraphic

	mu sync.RWMutex
}

var (
	instance *ACCSharedMemory
	once     sync.Once
)

// Instance returns the singleton instance
func Instance() *ACCSharedMemory {
	once.Do(func() {
		instance = &ACCSharedMemory{
			physicsMap:  "Local\\acpmf_physics",
			graphicsMap: "Local\\acpmf_graphics",
			staticMap:   "Local\\acpmf_static",
		}
		instance.ReadStaticPageFile(false)
		instance.ReadPhysicsPageFile(false)
		instance.ReadGraphicsPageFile(false)
	})
	return instance
}

// ReadGraphicsPageFile reads graphics data from shared memory
func (a *ACCSharedMemory) ReadGraphicsPageFile(fromCache bool) *SPageFileGraphic {
	a.mu.Lock()
	defer a.mu.Unlock()

	if fromCache && a.PageFileGraphic != nil {
		return a.PageFileGraphic
	}

	data := &SPageFileGraphic{}
	ToStruct(a.graphicsMap, data)

	a.PageFileGraphic = data
	return data
}

// ReadStaticPageFile reads static data from shared memory
func (a *ACCSharedMemory) ReadStaticPageFile(fromCache bool) *SPageFileStatic {
	a.mu.Lock()
	defer a.mu.Unlock()

	if fromCache && a.PageFileStatic != nil {
		return a.PageFileStatic
	}

	data := &SPageFileStatic{}
	ToStruct(a.staticMap, data)

	a.PageFileStatic = data
	return data
}

// ReadPhysicsPageFile reads physics data from shared memory
func (a *ACCSharedMemory) ReadPhysicsPageFile(fromCache bool) *SPageFilePhysics {
	a.mu.Lock()
	defer a.mu.Unlock()

	if fromCache && a.PageFilePhysics != nil {
		return a.PageFilePhysics
	}

	data := &SPageFilePhysics{}
	ToStruct(a.physicsMap, data)

	a.PageFilePhysics = data
	return data
}
