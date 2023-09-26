package main

import (
	"log"
	"pkv/api/src/router"
)

//	@title			DPV API
//	@version		1.0
//	@description	API to get data about DPV
//	@termsOfService	http://dpfv.de/

//	@contact.name	API Support
//	@contact.url	http://dpfv.de/
//	@contact.email	support@dpfv.de

//	@license.name	MIT License
//	@license.url	https://opensource.org/licenses/MIT

//	@host		localhost:8080
//	@BasePath	/

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	server := router.NewServer("config.yml")
	log.Fatal(server.ListenAndServe())
}
