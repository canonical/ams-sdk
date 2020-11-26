// -*- Mode: Go; indent-tabs-mode: t -*-
// AMS - Anbox Management Service
// Copyright 2018 Canonical Ltd.  All rights reserved.

package shared

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

var stdin = bufio.NewReader(os.Stdin)

// AskForBool asks a question and expect a yes/no answer.
func AskForBool(question string, defaultAnswer string) bool {
	for {
		answer := askQuestion(question, defaultAnswer)

		if StringInSlice(strings.ToLower(answer), []string{"yes", "y"}) {
			return true
		} else if StringInSlice(strings.ToLower(answer), []string{"no", "n"}) {
			return false
		}

		invalidInput()
	}
}

// AskForPassword asks the user to enter a password.
func AskForPassword(question string) string {
	fmt.Printf(question)
	pwd, _ := terminal.ReadPassword(0)
	fmt.Println("")

	return string(pwd)
}

// Ask a question on the output stream and read the answer from the input stream
func askQuestion(question, defaultAnswer string) string {
	fmt.Printf(question)

	return readAnswer(defaultAnswer)
}

// Read the user's answer from the input stream, trimming newline and providing a default.
func readAnswer(defaultAnswer string) string {
	answer, _ := stdin.ReadString('\n')
	answer = strings.TrimSuffix(answer, "\n")
	answer = strings.TrimSpace(answer)
	if answer == "" {
		answer = defaultAnswer
	}

	return answer
}

// Print an invalid input message on the error stream
func invalidInput() {
	fmt.Fprintf(os.Stderr, "Invalid input, try again.\n\n")
}
