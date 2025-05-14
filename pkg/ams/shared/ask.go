// -*- Mode: Go; indent-tabs-mode: t -*-
/*
 * This file is part of AMS SDK
 * Copyright 2021 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of the Lesser GNU General Public License version 3, as published
 * by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY, SATISFACTORY
 * QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.  See the Lesser GNU General Public
 * License for more details.
 *
 * You should have received a copy of the Lesser GNU General Public License along
 * with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
	fmt.Print(question)
	pwd, _ := terminal.ReadPassword(0)
	fmt.Println("")

	return string(pwd)
}

// Ask a question on the output stream and read the answer from the input stream
func askQuestion(question, defaultAnswer string) string {
	fmt.Print(question)

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
	fmt.Fprint(os.Stderr, "Invalid input, try again.\n\n")
}
