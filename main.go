// Package main is the entry point for the TickTask CLI application.
// It delegates all command processing to the cmd package.
package main

import "ticktask/cmd"

func main() {
	cmd.Execute()
}
