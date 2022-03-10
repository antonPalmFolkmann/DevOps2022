package main

import (
	"github.com/antonPalmFolkmann/DevOps2022/storage"
)

func main() {

	db := storage.ConnectPsql()
	storage.Migrate(db)

	// go func() {
	// 	r := mux.NewRouter()
	// 	simulator.SetupRoutes(r)
	// 	log.Fatalln(http.ListenAndServe(":8081", r))
	// }()

	// // Setup minitwit "website"
	// r := mux.NewRouter()
	// minitwit.SetupRoutes(r)
	// log.Fatalln(http.ListenAndServe(":8080", r))
}
