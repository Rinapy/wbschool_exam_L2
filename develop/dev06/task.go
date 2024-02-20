package main

import (
	"dev06/cut"
)

func main() {
	cfg, data := cut.ParseFlag()

	c := cut.CutApp{
		Cfg:  cfg,
		Line: data,
	}
	c.Run()
	//if res, err := cut.parseF("1-8"); err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(res)
	//}
}
