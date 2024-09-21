package main

import "fmt"

func main() {
	error := StartServer(":8090")
	if error != nil {
		fmt.Println(error)
	}
}
