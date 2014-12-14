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

type Config struct {
	Server, Channel, BotUser, BotNick, WeatherKey, LogDir string
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
}

func RandomString() string {
	return phrases[rand.Intn(len(phrases))]
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

	cmdArray = strings.SplitAfterN(cmdMsg, "!", 2)

	if len(cmdArray) > 0 {
		msgArray = strings.SplitN(cmdArray[1], " ", 2)
	}

	if len(msgArray) > 1 {
		cmd := fmt.Sprintf("%vs", msgArray[0])

		if strings.Contains(cmd, "weather") {
			msg = WeatherCmd()
		} else if strings.Contains(cmd, "cakeday") {
			msg = CakeDayCmd(msgArray[1])
		} else {
			msg = GenericVerbCmd(cmd, msgArray[1])
		}
	} else {
		if strings.Contains(msgArray[0], "help") {
			msg = HelpCmd()
		} else {
			msg = "I did not understand your command. Try '!slap Setsuna-Xero really hard'"
		}
	}
	return msg
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
	log := strings.TrimPrefix(CreateFile, "#")
	//Check if the Log File for the Channel(s) Exists if not create it
	if _, err := os.Stat(log + ".log"); os.IsNotExist(err) {
		fmt.Printf("Log File " + log + ".log Doesn't Exist. Creating Log File.\n")
		os.Create(log + ".log")
		fmt.Printf("Log File " + log + ".log Created.\n")
	} else {
		fmt.Printf("Log File Exists.\n")
	}
}

// Begin Bot Channel Logging.
func ChannelLogger(Log string, UserNick string, message string) {
	STime := time.Now().UTC().Format(time.ANSIC)
	log := strings.TrimPrefix(Log, "#")

	//Open the file for writing With Append Flag to create file persistence
	f, err := os.OpenFile(log+".log", os.O_RDWR|os.O_APPEND|os.O_SYNC, 0666)
	if err != nil {
		fmt.Println(err)
	}
	//And Write the Logs with timestamps :)
	n, err := io.WriteString(f, fmt.Sprintf("%v > %v: %v\n", STime, UserNick, message))
	if err != nil {
		fmt.Println(n, err)
	}
	f.Close()
}

// Commands

// GenericVerbCmd returns a message string based on the supplied cmd (a verb).
func GenericVerbCmd(cmd, extra string) string {
	// This should give us something like:
	//     "Snuffles slaps $USER, FOR SCIENCE!"
	// If given the command:
	//     "!slap $USER"
	randPhrase := RandomString()
	msg := fmt.Sprintf("\x01"+"ACTION %v %v, %v\x01", cmd, extra, randPhrase)
	return msg
}

// WeatherCmd is NYI
func WeatherCmd() string {
	// weatherArray := strings.Split(msgArray[1], " ", 2)
	// query := strings.Join(weatherArray[0], "")
	// msg = QueryWeather(query, config)
	msg := "Look outside, this feature isn't implemented just yet.\n"
	return msg
}

// CakeDayCmd returns a string containing the Reddit cakeday of a user
// upon success, or an error string on failure.
func CakeDayCmd(user string) string {
	// !cakeday $USER
	responseString, err := cakeday.Get(user)
	if err != nil {
		msg := fmt.Sprintf("I caught an error: %v\n", err)
		return msg
	}

	// >> Reddit Cake Day for $USER is: $CAKEDAY
	msg := fmt.Sprintf("%v\n", responseString)
	return msg
}

func HelpCmd() string {
	msgp1 := "Available commands: !help, !weather (NYI), !cakeday, !VERB\n"
	msgp2 := "!help: Display this help message.\n"
	msgp3 := "!weather <city>: Not yet implemented.\n"
	msgp4 := "!cakeday <username>: Get the Reddit Cake Day for the requested user.\n"
	msgp5 := "!VERB <msg>: Perform the selected verb in <msg> context. Example: !slap setkeh\n"
	msg := fmt.Sprintf("%s %s %s %s %s", msgp1, msgp2, msgp3, msgp4, msgp5)
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

	regex, _ := regexp.Compile(`<title>(.+?)<\/title>`)

	msgArray := strings.Split(msg, " ")

	for _, word = range msgArray {
		if strings.Contains(word, "http") || strings.Contains(word, "www") {
			url = word
			break
		}
	}

	// Band-AID. TODO: Fix this properly.
	if strings.Contains(url, "imgur") || strings.Contains(url, ".jpg") || strings.Contains(url, ".png") || strings.Contains(url, ".gif") {
		newMsg = fmt.Sprintf("Cannot resolve Image / Imgur links right now, beware...\n")
		return newMsg
	}

	resp, err := http.Get(url)

	if err != nil {
		newMsg = fmt.Sprintf("Could not resolve URL %v, beware...\n", url)
		return newMsg
	}

	defer resp.Body.Close()

	rawBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		newMsg = fmt.Sprintf("Could not read response Body of %v ...\n", url)
		return newMsg
	}

	body := string(rawBody)
	title = regex.FindStringSubmatch(body)[1]
	newMsg = fmt.Sprintf("[ %v ]->( %v )\n", title, url)

	return newMsg
}

//func QueryWeather(query string, config *Config) string {
//	var beginUrl string = "http://api.worldweatheronline.com/free/v1/weather.ashx?q="
//	var endUrl string = "&format=json&num_of_days=1&date=today&includelocation=yes&show_comments=no&key="

//	city := query
//	url := fmt.Sprintf("%v%v%v%v", beginUrl, city, endUrl, config.WeatherKey)
//	client := http.Client()
//	response, err := client.Get(url)

//	if err != nil {
//		return fmt.Sprintf("Caught error: %v", err.Error())
//	}
//	weatherJson := json.Decoder(response.Body)
//	weather := fmt.Sprintf("Weather for %v: %vC", weatherJson)
//	return weather

//}

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
			LogDir(config.LogDir)
			LogFile(config.LogDir + e.Arguments[0])
		}
		message := " has joined"
		ChannelLogger(config.LogDir+e.Arguments[0], e.Nick, message)
	})
	conn.AddCallback("PART", func(e *irc.Event) {
		message := " has parted"
		ChannelLogger(config.LogDir+e.Arguments[0], e.Nick, message)
	})
	conn.AddCallback("QUIT", func(e *irc.Event) {
		message := " has quit"
		ChannelLogger(config.LogDir+e.Arguments[0], e.Nick, message)
	})

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		var response string
		message := e.Message()

		if strings.Contains(message, "!") && strings.Index(message, "!") == 0 {
			// This is a command, parse it.
			response = ParseCmds(message, config)
		}

		if strings.Contains(message, "http") || strings.Contains(message, "www") {
			response = UrlTitle(message)
		}

		if len(response) > 0 {
			conn.Privmsg(config.Channel, response)
		}

		if len(message) > 0 {
			ChannelLogger(config.LogDir+e.Arguments[0], e.Nick+": ", message)
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

	rand.Seed(64)
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
