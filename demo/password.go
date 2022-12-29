//nolint
package main

import (
	"fmt"
	gap "github.com/byReqz/go-ask-password"
	"log"
)

func main() {
	user, err := gap.AskUser("Username: ")
	if err != nil {
		log.Fatal(err)
	}
	pw, err := gap.AskPassword("Password: ")
	if err != nil {
		log.Fatal(err)
	}
	tf, err := gap.AskKey("2FA: ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user + ":" + pw + ":" + tf)
}
