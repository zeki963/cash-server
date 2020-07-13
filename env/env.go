package main

import (
	"fmt"
	"os"
)

func main() {
	var USER string
	USER = os.Getenv("USER")
	fmt.Println(USER)
}
