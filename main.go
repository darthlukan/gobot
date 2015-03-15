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
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const delay = 40

type Config struct {
	Admins     []string
	Server     string
	Channel    string
	BotUser    string
	BotNick    string
	Trigger    string
	WeatherKey string
	LogDir     string
	WikiLink   string
	Homepage   string
	Forums     string
}

var quips = []string{
	"FOR SCIENCE!",
	"because... reasons.",
	"it's super effective!",
	"because... why not?",
	"was it good for you?",
	"given the alternative, yep, worth it!",
	"don't ask...",
	"then makes a sandwich.",
	"oh noes!",
	"did I do that?",
	"why must you turn this place into a house of lies!",
	"really???",
	"LLLLEEEEEERRRRRROOOOYYYY JEEEENNNKINNNS!",
	"DOH!",
	"Giggity!",
}

func RandomQuip() string {
	return quips[rand.Intn(len(quips))]
}

// ParseCmds takes PRIVMSG strings containing a preceding bang "!"
// and attempts to turn them into an ACTION that makes sense.
// Returns a msg string.
func ParseCmds(cmdMsg string, config *Config) string {
	var (
		msg      string
		msgArray []string
		cmdArray []string
	)

	cmdArray = strings.SplitAfterN(cmdMsg, config.Trigger, 2)

	if len(cmdArray) > 0 {
		msgArray = strings.SplitN(cmdArray[1], " ", 2)
	}

	if len(msgArray) > 1 {
		cmd := fmt.Sprintf("%vs", msgArray[0])
		switch {
		case strings.Contains(cmd, "cakeday"):
			msg = CakeDayCmd(msgArray[1])
		case strings.Contains(cmd, "ddg"), strings.Contains(cmd, "search"):
			query := strings.Join(msgArray[1:], " ")
			msg = SearchCmd(query)
		case strings.Contains(cmd, "convtemp"):
			query := strings.Join(msgArray[1:], " ")
			msg = ConvertTempCmd(query)
		default:
			msg = GenericVerbCmd(cmd, msgArray[1])
		}
	} else {
		switch {
		case strings.Contains(msgArray[0], "help"):
			msg = HelpCmd(config.Trigger)
		case strings.Contains(msgArray[0], "wiki"):
			msg = WikiCmd(config)
		case strings.Contains(msgArray[0], "homepage"):
			msg = HomePageCmd(config)
		case strings.Contains(msgArray[0], "forums"):
			msg = ForumCmd(config)
		default:
			msg = "I get it, you're just a human.  Try '!help'"
		}
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

	regex, _ := regexp.Compile(`(?i)<title>(.*?)<\/title>`)

	msgArray := strings.Split(msg, " ")

	for _, word = range msgArray {
		if strings.Contains(word, "http") || strings.Contains(word, "www") {
			url = word
			break
		}
	}

	resp, err := http.Get(url)

	if err != nil {
		return fmt.Sprintf("Could not resolve URL %v, beware...\n", url)
	}

	defer resp.Body.Close()

	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Could not read response Body of %v ...\n", url)
	}

	body := string(rawBody)
	noNewLines := strings.Replace(body, "\n", "", -1)
	noCarriageReturns := strings.Replace(noNewLines, "\r", "", -1)
	notSoRawBody := noCarriageReturns

	titleMatch := regex.FindStringSubmatch(notSoRawBody)
	if len(titleMatch) > 1 {
		title = strings.TrimSpace(titleMatch[1])
	} else {
		title = fmt.Sprintf("Title Resolution Failure")
	}
	newMsg = fmt.Sprintf("[ %v ]( %v )\n", title, url)

	return newMsg
}

// AddCallbacks is a single function that does what it says.
// It's merely a way of decluttering the main function.
func AddCallbacks(conn *irc.Connection, config *Config) {
	log := fmt.Sprintf("%s%s", config.LogDir, config.Channel)

	conn.AddCallback("001", func(e *irc.Event) {
		conn.Join(config.Channel)
	})

	conn.AddCallback("JOIN", func(e *irc.Event) {
		if e.Nick == config.BotNick {
			conn.Privmsg(config.Channel, "Hello everybody, I'm a bot")
			LogDir(config.LogDir)
			LogFile(config.LogDir + e.Arguments[0])
		}
		message := fmt.Sprintf("%s has joined", e.Nick)
		go ChannelLogger(log, e.Nick, message)
	})
	conn.AddCallback("PART", func(e *irc.Event) {
		message := fmt.Sprintf("has parted (%s)", e.Message())
		nick := fmt.Sprintf("%s@%s", e.Nick, e.Host)
		go ChannelLogger(log, nick, message)
	})
	conn.AddCallback("QUIT", func(e *irc.Event) {
		message := fmt.Sprintf("has quit (%v)", e.Message)
		nick := fmt.Sprintf("%s@%s", e.Nick, e.Host)
		go ChannelLogger(log, nick, message)
	})

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		var response string
		message := e.Message()
		if strings.Contains(message, config.Trigger) && strings.Index(message, config.Trigger) == 0 {
			response = ParseCmds(message, config)
		}
		if strings.Contains(message, "http://") || strings.Contains(message, "https://") || strings.Contains(message, "www.") {
			response = UrlTitle(message)
		}

		if strings.Contains(message, "quit") {
			QuitCmd(config.Admins, e.Nick)
		}

		if len(response) > 0 {
			conn.Privmsg(config.Channel, response)
		}

		if len(message) > 0 {
			if e.Arguments[0] != config.BotNick {
				go ChannelLogger(log, e.Nick+": ", message)
			}
		}
	})
}

func main() {

	rand.Seed(64)
	file, err := os.Open("config.json")

	if err != nil {
		fmt.Println("Couldn't read config file, dying...")
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	decoder.Decode(&config)

	conn := irc.IRC(config.BotNick, config.BotUser)
	err = conn.Connect(config.Server)

	if err != nil {
		fmt.Println("Failed to connect.")
		panic(err)
	}

	AddCallbacks(conn, config)
	conn.Loop()
}
