package model

import (
	"fmt"
	"strings"
)

/*BoardItem is a single card on a lane*/
type BoardItem struct {
	ID       int32
	Name     string
	Severity SeverityType
}

/*SeverityType enum*/
type SeverityType uint8

/*Values for SeverityType*/
const (
	LOW SeverityType = iota
	MEDIUM
	HIGH
)

var severityNames = map[SeverityType]string{LOW: "low", MEDIUM: "medium", HIGH: "high"}

func (s SeverityType) String() string {
	if name, ok := severityNames[s]; ok {
		return name
	}
	return fmt.Sprintf("severity(%d)", uint8(s))
}

/*Marker returns the short symbol shown next to the card id*/
func (s SeverityType) Marker() string {
	switch s {
	case MEDIUM:
		return "*"
	case HIGH:
		return "!"
	default:
		return ""
	}
}

/*ParseSeverity converts "low", "medium" or "high" into a SeverityType*/
func ParseSeverity(value string) (SeverityType, error) {
	for severity, name := range severityNames {
		if strings.EqualFold(value, name) {
			return severity, nil
		}
	}
	return LOW, fmt.Errorf("unknown severity %q, use low, medium or high", value)
}
