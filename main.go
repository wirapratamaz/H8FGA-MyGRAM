package main

import "github.com/wirapratamaz/H8FGA-MyGRAM/routers"

// @title FINAL PROJECT GOLANG HACTUV8
// @version 1.0
// @description This is a sample server todo server. You can visit the GitHub repository at https://github.com/LordGhostX/swag-gin-demo
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @query.collection.format multi
func main() {
	r := routers.StartApp()
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
