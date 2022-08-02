package main

import "strings"

type arrayFlags []string

func (a *arrayFlags) String() string {
	return strings.Join(*a, ",")
}

func (a *arrayFlags) Set(v string) error {
	*a = append(*a, v)
	return nil
}

// func getReceiveCmd(args []string) (cmd receiveCmd) {

// 	channels := args[1 : len(args)-1]
// 	path := args[len(args)-1]

// 	cmd.channels = channels
// 	cmd.folderPath = path

// 	return
// }
