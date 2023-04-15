package main

import "github.com/wirapratamaz/H8FGA-MyGRAM/routers"

func main() {
	r := routers.StartApp()
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
