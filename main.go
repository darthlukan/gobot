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
	"encoding/json"
	"fmt"
	"github.com/thoj/go-ircevent"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Config struct {
	Server, Channel, BotUser, BotNick string
}

// ParseCmds takes PRIVMSG strings containing a preceding bang "!"
// and attempts to turn them into an ACTION that makes sense.
// Returns a msg string.
func ParseCmds(cmdMsg string) string {
	var (
		msg      string
		msgArray []string
		cmdArray []string
	)

	cmdArray = strings.SplitAfterN(cmdMsg, "!", 2)

	if len(cmdArray) > 0 {
		msgArray = strings.SplitN(cmdArray[1], " ", 2)
	}

	if len(msgArray) > 1 {
		cmd := fmt.Sprintf("%vs", msgArray[0])

		// This should give us something like:
		//     "Snuffles slaps $USER, FOR SCIENCE!"
		// If given the command:
		//     "!slap $USER"
		msg = fmt.Sprintf("\x01"+"ACTION %v %v, FOR SCIENCE!\x01", cmd, msgArray[1])
	} else {
		msg = "I did not understand your command. Try '!slap Setsuna-Xero really hard'"
	}
	return msg
}

// UrlTitle attempts to extract the title of the page that a
// pasted URL points to.
// Returns a string message with the title and URL on success, or a string
// with an error message on failure.
func UrlTitle(msg string) string {
	var (
		newMsg, url, title, word string
	)

	regex, _ := regexp.Compile(`<title[^>]*>([^<]+)<\/title>`)

	msgArray := strings.Split(msg, " ")

	for _, word = range msgArray {
		if strings.Contains(word, "http") || strings.Contains(word, "www") {
			url = word
			break
		}
	}

	resp, err := http.Get(word)

	if err != nil {
		newMsg = fmt.Sprintf("Could not resolve URL %v, beware...\n", word)
		return newMsg
	}

	defer resp.Body.Close()

	rawBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		newMsg = fmt.Sprintf("Could not read response Body of %v ...", word)
		return newMsg
	}

	body := string(rawBody)
	title = regex.FindString(body)
	newMsg = fmt.Sprintf("[ %v ]->( %v )", title, url)

	return newMsg
}

func QueryGoogle(query string) string {
	var results string
	// TODO: Logic!
	return results
}

// AddCallbacks is a single function that does what it says.
// It's merely a way of decluttering the main function.
func AddCallbacks(conn *irc.Connection, config *Config) {
	conn.AddCallback("001", func(e *irc.Event) {
		conn.Join(config.Channel)
	})

	conn.AddCallback("JOIN", func(e *irc.Event) {
		if e.Nick == config.BotNick {
			conn.Privmsg(config.Channel, "Hello everybody, I'm a bot")
		}
	})

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		var response string
		message := e.Message()

		if strings.Contains(message, "!") && strings.Index(message, "!") == 0 {
			// This is a command, parse it.
			response = ParseCmds(message)
		}

		if strings.Contains(message, "http") || strings.Contains(message, "www") {
			response = UrlTitle(message)
		}

		if len(response) > 0 {
			conn.Privmsg(config.Channel, response)
		}
	})
}

// Connect tries up to three times to get a connection to the server
// and channel, hopefully with a nil err value at some point.
// Returns error
func Connect(conn *irc.Connection, config *Config) error {
	var err error

	for attempt := 1; attempt <= 3; attempt++ {
		if err = conn.Connect(config.Server); err != nil {
			fmt.Println("Connection attempt %v failed, trying again...", attempt)
			continue
		} else {
			break
		}
	}
	return err
}

func main() {

	// Read the config file and populate our Config struct.
	file, err := os.Open("config.json")

	if err != nil {
		fmt.Println("Couldn't read config file, dying...")
		panic(err)
	}

	decoder := json.NewDecoder(file)
	config := &Config{}
	decoder.Decode(&config)

	conn := irc.IRC(config.BotNick, config.BotUser)
	err = Connect(conn, config)

	if err != nil {
		fmt.Println("Failed to connect.")
		// Without a connection, we're useless, panic and die.
		panic(err)
	}

	AddCallbacks(conn, config)
	conn.Loop()
}
