package kingpinOp

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

func Read(app *kingpin.Application,arg[] string)(string){
	result:=kingpin.MustParse(app.Parse(arg))
	return result
}