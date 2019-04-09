package main

import (
	"fmt"
	"github.com/jvogelsb/golang-challenge-1-drum_machine/drum"
	"log"
)

func main() {
	fmt.Println("TESTING")
	var file = "fixtures/pattern_5.splice"
	_, err := drum.DecodeFile(file)
	if err != nil {
		log.Printf("ERROR: Unable to parse file: %s", file)
	}

	fmt.Println("Parsed File!!!!!")

}