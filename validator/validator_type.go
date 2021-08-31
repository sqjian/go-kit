package validator

//go:generate stringer -type=ValidatorType  -linecomment
type ValidatorType int

const (
	_ ValidatorType = iota
	Json
)
