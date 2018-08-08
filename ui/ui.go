package ui

import (
	"fmt"
)

// UI interface definition ----------------------------------------
type JrepackUi interface {
	Error(message string)
	OnEnd(message string)
}

// default UI implementation --------------------------------------

type defaultUi struct {
}

func (ui defaultUi) Error(message string) {
	panic(message)
}

func (ui defaultUi) OnEnd(message string) {
	fmt.Println(message)
}

// UI implementation setter and getter ----------------------------

var (
	someUi JrepackUi = defaultUi{}
)

func Set(ui JrepackUi) {
	someUi = ui
}

func Current() JrepackUi {
	return someUi
}
