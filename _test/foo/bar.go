package foo

import "github.com/gmemcc/yaegi/_test/foo/boo"

var Bar = "BARR"
var Boo = boo.Boo

func init() { println("init foo") }
