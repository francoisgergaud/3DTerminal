package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
)

func main() {
	// CPU profiling by default
	// go func() {
	// 	log.Printf("Starting Server! \t Go to http://localhost:6060/debug/pprof/\n")
	// 	err := http.ListenAndServe("localhost:6060", nil)
	// 	if err != nil {
	// 		log.Printf("Failed to start the server! Error: %v", err)
	// 	}
	// }()
	fmt.Println("terminal: " + os.Getenv("TERM"))
	err := InitGame()
	if err != nil {
		panic(err)
	}
}
