//nolint
package main

import (
	"fmt"
	gap "github.com/byReqz/go-ask-password"
	"log"
)

func main() {
	tk, err := gap.ScanSecret("Token: ", "-", "Tab to toggle visibility")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tk)
}
