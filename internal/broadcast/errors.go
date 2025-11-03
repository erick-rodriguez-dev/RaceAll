package broadcast

import (
	"RaceAll/internal/errors"
)

const moduleName = "broadcast"

func NewError(op string, err error) error {
	return errors.NewError(moduleName, op, err)
}

func NewErrorWithContext(op string, err error, ctx string) error {
	return errors.NewErrorWithContext(moduleName, op, err, ctx)
}

func NewValidationError(field string, value any, rule string) error {
	return errors.NewValidationError(moduleName, field, value, rule)
}

func ValidateCarIndex(index uint16) error {
	if index > 10000 {
		return NewValidationError("carIndex", index, "must be less than 10000")
	}
	return nil
}

func ValidateSessionType(sessionType RaceSessionType) error {
	switch sessionType {
	case RaceSessionTypePractice,
		RaceSessionTypeQualifying,
		RaceSessionTypeSuperpole,
		RaceSessionTypeRace,
		RaceSessionTypeHotlap,
		RaceSessionTypeHotstint,
		RaceSessionTypeHotlapSuperpole,
		RaceSessionTypeReplay:
		return nil
	default:
		return NewValidationError("sessionType", sessionType, "unknown session type")
	}
}

func ValidateSessionPhase(phase SessionPhase) error {
	switch phase {
	case SessionPhaseNone,
		SessionPhaseStarting,
		SessionPhasePreFormation,
		SessionPhaseFormationLap,
		SessionPhasePreSession,
		SessionPhaseSession,
		SessionPhaseSessionOver,
		SessionPhasePostSession,
		SessionPhaseResultUI:
		return nil
	default:
		return NewValidationError("sessionPhase", phase, "unknown session phase")
	}
}

func ValidateCarLocation(location CarLocationEnum) error {
	switch location {
	case CarLocationNone,
		CarLocationTrack,
		CarLocationPitlane,
		CarLocationPitEntry,
		CarLocationPitExit:
		return nil
	default:
		return NewValidationError("carLocation", location, "unknown car location")
	}
}

func ValidateBroadcastingEventType(eventType BroadcastingEventType) error {
	switch eventType {
	case BroadcastingEventTypeNone,
		BroadcastingEventTypeGreenFlag,
		BroadcastingEventTypeSessionOver,
		BroadcastingEventTypePenaltyCommMsg,
		BroadcastingEventTypeAccident,
		BroadcastingEventTypeLapCompleted,
		BroadcastingEventTypeBestSessionLap,
		BroadcastingEventTypeBestPersonalLap:
		return nil
	default:
		return NewValidationError("eventType", eventType, "unknown event type")
	}
}

func ValidateDriverCategory(category DriverCategory) error {
	switch category {
	case DriverCategoryBronze,
		DriverCategorySilver,
		DriverCategoryGold,
		DriverCategoryPlatinum,
		DriverCategoryError:
		return nil
	default:
		return NewValidationError("driverCategory", category, "unknown driver category")
	}
}
