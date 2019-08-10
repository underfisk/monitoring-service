package main

type LogSender struct {
	id int
	name string
	sentAt string
	applicationId int

}

type LogType struct {

}

/** JSON like**/
type LogPayload struct {

}

type Log struct {
	sender LogSender
	_type LogType
	payload LogPayload
}

