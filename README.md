GoBot
========


## Description

> GoBot is my first stab at an IRC bot written in Go.

## Installation
```
    $ go get github.com/darthlukan/gobot
```
> Or:

```
    $ mkdir -p $GOPATH/src/github.com/darthlukan
    $ cd $GOPATH/src/github.com/darthlukan
    $ git clone git@github.com:darthlukan/gobot.git
    $ cd gobot
```

## Usage

> Edit the config.json file located in $GOPATH/src/github.com/darthlukan/gobot to your preferences.

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

- 1. Add commands that actually do something useful
- 2. Parse URLs pasted in channels to verify that they are legit.
- 3. Google Search
- 4. Logging
- 5. Tests would be nice >.>

## License

> GPLv2, see LICENSE file.
