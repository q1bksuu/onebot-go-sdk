package entity

type ActionResponseStatus string

const (
	StatusOK     ActionResponseStatus = "ok"
	StatusAsync  ActionResponseStatus = "async"
	StatusFailed ActionResponseStatus = "failed"
)

type ActionResponseRetcode int

const (
	RetcodeSuccess ActionResponseRetcode = 0
	RetcodeAsync   ActionResponseRetcode = 1
)
