package main

import (
	githubactions "github.com/sethvargo/go-githubactions"

	"github.com/ethanthatonekid/gitcord"
)

func main() {
	action := githubactions.New()
	err := requireconditional.Run(action)
	if err != nil {
		action.Fatalf("%v", err)
	}
}
