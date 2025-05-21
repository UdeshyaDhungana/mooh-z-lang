package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/udeshyadhungana/interprerer/app/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("नमस्कार %s मुजी!\n", user.Username)
	fmt.Println("यो \"मुजी\" भाषा हो। तल लेख् मुजी 👇")
	repl.Start(os.Stdin, os.Stdout)
}
