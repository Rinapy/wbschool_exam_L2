package main

import (
	"dev06/cut"
	"fmt"
)

func main() {
	_, data := cut.ParseFlag()
	fmt.Println(data)

	//c := cut.CutApp{
	//	Cfg:  cfg,
	//	Line: data,
	//}
	//if res, err := cut.parseF("1-8"); err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(res)
	//}
}
