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
	"github.com/ajanicij/goduckgo/goduckgo"
	"github.com/darthlukan/cakeday"
	"github.com/thoj/go-ircevent"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const delay = 40

type Config struct {
	Server, Channel, BotUser, BotNick, Trigger, WeatherKey, LogDir, WikiLink, Homepage, Forums string
}

var phrases = []string{
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
}

func RandomString() string {
	return phrases[rand.Intn(len(phrases))]
}

// Begin Bot Channel Logging.
func ChannelLogger(Log string, UserNick string, message string) {
	STime := time.Now().UTC().Format(time.ANSIC)
	log := strings.Replace(Log, "#", "", 1)
	logFile := fmt.Sprintf("%s.log", log)

	//Open the file for writing With Append Flag to create file persistence
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_SYNC, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	//And Write the Logs with timestamps :)
	n, err := io.WriteString(f, fmt.Sprintf("%v > %v: %v\n", STime, UserNick, message))
	if err != nil {
		fmt.Println(n, err)
	}
}

func LogDir(CreateDir string) {
	//Check if the LogDir Exists. And if not Create it.
	if _, err := os.Stat(CreateDir); os.IsNotExist(err) {
		fmt.Printf("No such file or directory: %s\n", CreateDir)
		os.Mkdir(CreateDir, 0777)
	} else {
		fmt.Printf("Its There: %s\n", CreateDir)
	}
}

func LogFile(CreateFile string) {
	log := strings.Replace(CreateFile, "#", "", 1)
	logFile := fmt.Sprintf("%s.log", log)
	//Check if the Log File for the Channel(s) Exists if not create it
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		fmt.Printf("Log File %s Doesn't Exist. Creating Log File.\n", logFile)
		os.Create(logFile)
		fmt.Printf("Log File %s Created.\n", logFile)
	} else {
		fmt.Printf("Log File %s Exists.\n", logFile)
	}
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
		case strings.Contains(cmd, "weather"):
			msg = WeatherCmd()
		case strings.Contains(cmd, "cakeday"):
			msg = CakeDayCmd(msgArray[1])
		case strings.Contains(cmd, "ddg"), strings.Contains(cmd, "search"):
			query := strings.Join(msgArray[1:], " ")
			msg = WebSearch(query)
		default:
			msg = GenericVerbCmd(cmd, msgArray[1])
		}
	} else {
		switch {
		case strings.Contains(msgArray[0], "help"):
			msg = HelpCmd()
		case strings.Contains(msgArray[0], "wiki"):
			msg = WikiCmd(config)
		case strings.Contains(msgArray[0], "homepage"):
			msg = HomePageCmd(config)
		case strings.Contains(msgArray[0], "forums"):
			msg = ForumCmd(config)
		default:
			msg = "I did not understand your command. Try '!slap Setsuna-Xero really hard'"
		}
	}
	return msg
}

// Commands

// GenericVerbCmd returns a message string based on the supplied cmd (a verb).
func GenericVerbCmd(cmd, extra string) string {
	// This should give us something like:
	//     "Snuffles slaps $USER, FOR SCIENCE!"
	// If given the command:
	//     "!slap $USER"
	randPhrase := RandomString()
	return fmt.Sprintf("\x01"+"ACTION %v %v, %v\x01", cmd, extra, randPhrase)
}

// WeatherCmd is NYI
func WeatherCmd() string {
	// weatherArray := strings.Split(msgArray[1], " ", 2)
	// query := strings.Join(weatherArray[0], "")
	// msg = QueryWeather(query, config)
	return fmt.Sprintf("Look outside, this feature isn't implemented just yet.\n")
}

// CakeDayCmd returns a string containing the Reddit cakeday of a user
// upon success, or an error string on failure.
func CakeDayCmd(user string) string {
	var msg string
	// !cakeday $USER
	responseString, err := cakeday.Get(user)
	if err != nil {
		msg = fmt.Sprintf("I caught an error: %v\n", err)
	} else {
		// >> Reddit Cake Day for $USER is: $CAKEDAY
		msg = fmt.Sprintf("%v\n", responseString)
	}
	return msg
}

func HelpCmd() string {
	return fmt.Sprintf("Available commands: !help, !ddg/search !weather (NYI), !cakeday, !VERB\n")
}

func WikiCmd(config *Config) string {
	return fmt.Sprintf("(Channel Wiki)[%s]\n", config.WikiLink)
}

func HomePageCmd(config *Config) string {
	return fmt.Sprintf("(Channel Homepage)[%s]\n", config.Homepage)
}

func ForumCmd(config *Config) string {
	return fmt.Sprintf("(Channel Forums)[%s]\n", config.Forums)
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

func WebSearch(query string) string {
	msg, err := goduckgo.Query(query)
	if err != nil {
		return fmt.Sprintf("DDG Error: %v\n", err)
	}

	switch {
	case len(msg.RelatedTopics) > 0:
		return fmt.Sprintf("First Topical Result: [ %s ]( %s )\n", msg.RelatedTopics[0].FirstURL, msg.RelatedTopics[0].Text)
	case len(msg.Results) > 0:
		return fmt.Sprintf("First External result: [ %s ]( %s )\n", msg.Results[0].FirstURL, msg.Results[0].Text)
	case len(msg.Redirect) > 0:
		return fmt.Sprintf("Redirect result: %s\n", UrlTitle(msg.Redirect))
	default:
		return fmt.Sprintf("Query: '%s' returned no results.\n", query)
	}
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
			LogDir(config.LogDir)
			LogFile(config.LogDir + e.Arguments[0])
		}
		message := fmt.Sprintf("%s has joined", e.Nick)
		go ChannelLogger(config.LogDir+e.Arguments[0], e.Nick, message)
	})
	conn.AddCallback("PART", func(e *irc.Event) {
		pmessage := "parted"
		message := e.Message()
		go ChannelLogger(config.LogDir+config.Channel, e.Nick+"@"+e.Host, pmessage+" "+"("+message+")")
	})
	conn.AddCallback("QUIT", func(e *irc.Event) {
		qmessage := "has quit"
		message := e.Message()
		go ChannelLogger(config.LogDir+config.Channel, e.Nick+"@"+e.Host, qmessage+" "+"("+message+")")
	})

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		var response string
		message := e.Message()
		switch {
		case strings.Contains(message, config.Trigger) && strings.Index(message, config.Trigger) == 0:
			// This is a command, parse it.
			response = ParseCmds(message, config)
		case strings.Contains(message, "http://"), strings.Contains(message, "https://"), strings.Contains(message, "www."):
			response = UrlTitle(message)
		}

		if len(response) > 0 {
			conn.Privmsg(config.Channel, response)
		}

		if len(message) > 0 {
			if e.Arguments[0] != config.BotNick {
				go ChannelLogger(config.LogDir+e.Arguments[0], e.Nick+": ", message)
			}
		}
	})
}

// Connect tries up to three times to get a connection to the server
// and channel, hopefully with a nil err value at some point.
// Returns error
func Connect(conn *irc.Connection, config *Config) error {
	var err error

	for attempt := 1; attempt <= 3; attempt++ {
		time.Sleep(delay * time.Second)
		if err = conn.Connect(config.Server); err != nil {
			fmt.Println("Connection attempt %v failed, trying again...", attempt)
		} else {
			break
		}
	}
	return err
}

func main() {

	rand.Seed(64)
	// Read the config file and populate our Config struct.
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
	err = Connect(conn, config)

	if err != nil {
		fmt.Println("Failed to connect.")
		// Without a connection, we're useless, panic and die.
		panic(err)
	}

	AddCallbacks(conn, config)
	conn.Loop()
}
