GoBot
========


## Description

> GoBot is my first stab at an IRC bot written in Go.

## Installation
```
    $ go get github.com/darthlukan/gobot
```

## Usage

> Please edit $GOPATH/src/github.com/darthlukan/gobot/main.go and change the channel,
> server, botNick, and botUser variables to your preference before running the bot!

> After those variables have been edited, you can run:
```
    $ go install .    # Note the '.'
    $ gobot
```

## in-channel interaction

> For now, there aren't really any "real" commands, but if you prefix with "!", the bot will spit out a message.

> Example:
```
    !slap SomeUser really hard
    >> *$botNick slaps SomeUser really hard, FOR SCIENCE!
```

## TODO

- 1. Reorganize the code
- 2. Use a config file
- 3. Add commands that actually do something useful
- 4. Parse URLs pasted in channels to verify that they are legit.
- 5. Google Search
- 6. Logging
- 7. Tests would be nice >.>

## License

> GPLv2, see LICENSE file.
