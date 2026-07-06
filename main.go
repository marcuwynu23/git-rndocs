package main

import "github.com/marcuwynu23/git-rndocs/cmd"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Execute(version)
}
