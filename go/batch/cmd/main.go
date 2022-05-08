package main

import (
	"fmt"
)

type Result struct {
	Col1 string
	Col2 string
	Col3 string
}

func main() {
	var results []Result
	db.Raw("select col1, col2, col3 from testschema.sample").Scan(&results)
	for _, res := range results {
		fmt.Println(res)
	}
}
