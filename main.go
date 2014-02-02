/*
GoBot

An IRC bot written in Go.

Copyright (C) 2014  Brian C. Tomlinson

Contact: brian.tomlinson@linux.com

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License along
with this program; if not, write to the Free Software Foundation, Inc.,
51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
*/
package main

import (
	"fmt"
	"github.com/thoj/go-ircevent"
	"strings"
)

// TODO: These should really be in a config file...
var (
	channel          = "#tinfoilhats"
	server           = "irc.freenode.net:6667"
	botNick, botUser = "Snuffles_test", "Snuffles_test"
)

// ParseCmds takes PRIVMSG strings containing a preceding bang "!"
// and attempts to turn them into an ACTION that makes sense.
// Returns a msg string.
func ParseCmds(cmdMsg string) string {
	cmdArray := strings.SplitAfterN(cmdMsg, "!", 2)
	msgArray := strings.SplitN(cmdArray[1], " ", 2)
	cmd := fmt.Sprintf("%vs", msgArray[0])

	// This should give us something like:
	//     "Snuffles slaps $USER, FOR SCIENCE!"
	// If given the command:
	//     "!slap $USER"
	msg := fmt.Sprintf("\x01"+"ACTION %v %v, FOR SCIENCE!\x01", cmd, msgArray[1])
	return msg
}

// AddCallbacks is a single function that does what it says.
// It's merely a way of decluttering the main function.
func AddCallbacks(conn *irc.Connection) {
	conn.AddCallback("001", func(e *irc.Event) {
		conn.Join(channel)
	})

	conn.AddCallback("JOIN", func(e *irc.Event) {
		if e.Nick == botNick {
			conn.Privmsg(channel, "Hello everybody, I'm a bot")
		}
	})

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		var response string

		if strings.Contains(e.Message, "!") && strings.Index(e.Message, "!") == 0 {
			// This is a command, parse it.
			response = ParseCmds(e.Message)
			conn.Privmsg(channel, response)
		}
	})
}

// Connect tries up to three times to get a connection to the server
// and channel, hopefully with a nil err value at some point.
// Returns error
func Connect(conn *irc.Connection) error {
	var err error

	for attempt := 1; attempt <= 3; attempt++ {
		if err = conn.Connect(server); err != nil {
			fmt.Println("Connection attempt %v failed, trying again...", attempt)
			continue
		} else {
			break
		}
	}
	return err
}

func main() {
	conn := irc.IRC(botNick, botUser)
	err := Connect(conn)

	if err != nil {
		fmt.Println("Failed to connect.")
		// Without a connection, we're useless, panic and die.
		panic(err)
	}

	AddCallbacks(conn)
	conn.Loop()
}
