package main

import (
	"fmt"
	gap "github.com/byReqz/go-ask-password"
	"log"
)

func main() {
	user, err := gap.AskUser()
	if err != nil {
		log.Fatal(err)
	}
	pw, err := gap.AskPassword()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user + ":" + pw)
}
