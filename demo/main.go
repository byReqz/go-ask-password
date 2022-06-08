package main

import (
	"fmt"
	gap "github.com/byReqz/go-ask-password"
	"log"
)

func main() {
	pw, err := gap.AskPassword("Password: ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pw)
}
