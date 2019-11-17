package tools

import "gopkg.in/gookit/color.v1"

func ShowCommonMessage(str string) {
	color.Info.Println(str)
}

func ShowSuccessMessage(str string) {
	color.Green.Println(str)
}

func ShowErrorMessage(str string) {
	color.Red.Println(str)
}

func ShowError(err error) {
	ShowErrorMessage(err.Error())
}
