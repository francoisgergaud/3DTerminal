package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

func main() {
	fmt.Println("terminal: " + os.Getenv("TERM"))
	tcell.SetEncodingFallback(tcell.EncodingFallbackUTF8)
	game := InitGame()
	game.Start()
}
