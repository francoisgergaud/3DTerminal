package main

import (
	"flag"
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
	var mode = flag.String("mode", "local", "possible mode: 'local', 'remote', 'remoteClient', 'remoteServer'")
	var remoteAddress = flag.String("address", "127.0.0.1:9836", "remote-server host-port")
	var serverPort = flag.String("port", "9836", "remote-server host-port")
	flag.Parse()
	game := NewGame()
	var err error
	if *mode == "local" {
		err = game.InitLocalGame()
	} else if *mode == "remote" {
		err = game.InitRemoteGame(*serverPort)
	} else if *mode == "remoteClient" {
		err = game.InitRemoteClient(*remoteAddress)
	} else if *mode == "remoteServer" {
		err = game.InitRemoteServer(*serverPort)
	}
	if err != nil {
		panic(err)
	}
}
