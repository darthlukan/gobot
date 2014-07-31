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

## Logging
> Gobot can now Log Channels.
> Set the "LogDir" in config.json
> NOTE: the Directory Must be Writable by the user executing the bot.

## TODO

- 1. Add commands that actually do something useful
- 2. Google Search
- 3. ~~~Logging~~~
- 4. Tests would be nice >.>

## License

> GPLv2, see LICENSE file.


[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/darthlukan/gobot/trend.png)](https://bitdeli.com/free "Bitdeli Badge")

