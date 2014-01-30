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

var (
	channel          = "#tinfoilhats"
	server           = "irc.freenode.net:6667"
	botNick, botUser = "Snuffles", "Snuffles"
)

func parseCmds(cmdMsg string) string {
	cmdArray := strings.SplitAfterN(cmdMsg, "!", 2)
	msgArray := strings.SplitN(cmdArray[1], " ", 2)
	cmd := fmt.Sprintf("%vs", msgArray[0])

	// This should give us something like:
	// "Snuffles slaps $USER, FOR SCIENCE!"
	// If given the command: "!slap $USER"
	msg := fmt.Sprintf("\x01"+"ACTION %v %v, FOR SCIENCE!\x01", cmd, msgArray[1])
	return msg
}

func main() {
	conn := irc.IRC(botNick, botUser)
	err := conn.Connect(server)

	if err != nil {
		fmt.Println("Failed to connect.")
	}

	conn.AddCallback("001", func(e *irc.Event) {
		conn.Join(channel)
	})

	conn.AddCallback("JOIN", func(e *irc.Event) {
		conn.Privmsg(channel, "Hello, I'm a bot")
	})

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		var response string

		if strings.Contains(e.Message, "!") && strings.Index(e.Message, "!") == 0 {
			// This is a command, parse it.
			response = parseCmds(e.Message)
			conn.Privmsg(channel, response)
		}
	})

	conn.Loop()
}
