package main

import (
	"fmt"
	"test/action"
)

func main(){
	test := 1
	switch test {
		case 0:
			action.FirstTCROSS()
		case 1:
			action.AccessTCROSS()
		case 2:
			action.TCROSSCache()
		default:
			fmt.Println("Please Select the correct scenario.")
	}
}
