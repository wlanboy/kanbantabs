package model

/*BoardItem struct*/
type BoardItem struct {
	ID       int32
	Name     string
	Severity SeverityType
}

/*SeverityType enum*/
type SeverityType uint8

/*const for SeverityType enum*/
const (
	LOW = iota
	MEDIUM
	HIGH
)
