package main

import "fmt"

func main() {
	fmt.Println("The server port is 8090")
	error := StartServer(":8090")
	if error != nil {
		fmt.Println(error)
	}
}
