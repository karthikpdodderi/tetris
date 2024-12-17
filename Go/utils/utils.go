package utils

import(
	"fmt"
)

func ClerLines(n int){
	fmt.Printf("\033[%dA\033[J", n)
}