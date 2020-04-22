package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"

	"github.com/gdamore/tcell"
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
	tcell.SetEncodingFallback(tcell.EncodingFallbackUTF8)
	screen, e := tcell.NewScreen()
	if e != nil {
		panic(e)
	}
	if e = screen.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		panic(e)
	}
	err := InitGame(screen)
	if err != nil {
		panic(err)
	}
}
