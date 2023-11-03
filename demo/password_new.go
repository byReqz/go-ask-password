// nolint
package main

import (
	"fmt"
	"log"

	gap "github.com/byReqz/go-ask-password"
)

func main() {
	user, err := gap.NewScan(gap.Options{
		Plain:   true,
		Prefix:  "Username: ",
		BreakAt: 1,
	})
	if err != nil {
		log.Fatal(err)
	}
	pw, err := gap.NewScan(gap.Options{
		Prefix:      "Password: ",
		Substitute:  "*",
		Placeholder: "TAB for reveal",
	})
	fmt.Println(user + ":" + pw)

	user, err = gap.AskUser("Username: ")
	if err != nil {
		log.Fatal(err)
	}
	pw, err = gap.AskPassword("Password: ")
	if err != nil {
		log.Fatal(err)
	}
	tf, err := gap.AskKey("2FA: ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user + ":" + pw + ":" + tf)
}
