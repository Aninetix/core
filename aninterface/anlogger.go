package aninterface

type AnLogger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}
