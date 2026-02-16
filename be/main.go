package main

import (
	"api/src/config"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

// func init() {
// 	key := make([]byte, 64)

// 	if _, err := rand.Read(key); err != nil {
// 		log.Fatal(err)
// 	}

// 	base64str := base64.StdEncoding.EncodeToString(key)

// 	fmt.Print(base64str)
// }

func main() {
	config.LoadConfig()
	r := router.Generate()
	fmt.Printf("listening at port :%d", config.ApiPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.ApiPort), r))
}
