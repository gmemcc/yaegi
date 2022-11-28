package main

import (
	"log"

	"github.com/gmemcc/yaegi/interp"
)

func main() {
	log.SetFlags(log.Lshortfile)
	i := interp.New(interp.Options{})
	i.Use(interp.Symbols)
	if _, err := i.Eval(`import "github.com/gmemcc/yaegi/interp"`); err != nil {
		log.Fatal(err)
	}
	if _, err := i.Eval(`i := interp.New(interp.Options{})`); err != nil {
		log.Fatal(err)
	}
	if _, err := i.Eval(`i.Eval("println(42)")`); err != nil {
		log.Fatal(err)
	}
}
