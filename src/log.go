package core


type LogType int32

const (
	ERROR LogType = 0
	LOG LogType = 1
)


type Log struct {
	id int
	name string
	applicationId int
	_type LogType
	payload interface {}
}

