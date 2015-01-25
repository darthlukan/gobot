GoBot
========

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/darthlukan/gobot?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/darthlukan/gobot.svg?branch=master)](https://travis-ci.org/darthlukan/gobot)


## Description

> GoBot is my first stab at an IRC bot written in Go. The goal is for it to become less of a "dumb bot" and more
of an semi-clever assistant.

## Installation

**Important**: Requires Go>=1.3

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

> For now, the available commands are fairly limited, here's what's available:

* !search or !ddg: Execute a search via DuckDuckGo

```
    !search Los Angeles
    !ddg New York
```

> !bangs support has also been included, so you can get results from Google as well as many other sources:

```
    !search !google weather in Los Angeles
    !search !archwiki i3
    !ddg !godoc cakeday
```

* !cakeday: Find the Reddit cakeday for a user

```
    !cakeday darthlukan
```

* !<VERB>: Echo the <VERB> and add a random quip. 

```
    !slap SomeUser really hard
    >> *$botNick slaps SomeUser really hard, FOR SCIENCE!
```

## Logging
> Gobot can now Log Channels.

> Set the "LogDir" in config.json

> NOTE: the Directory Must be Writable by the user executing the bot.

## TODO

- 1. Add commands that actually do something useful
- 2. ~~DuckDuckGo Search~~
- 3. ~~Logging~~
- 4. Tests would be nice >.>

## License

> GPLv2, see LICENSE file.


[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/darthlukan/gobot/trend.png)](https://bitdeli.com/free "Bitdeli Badge")

