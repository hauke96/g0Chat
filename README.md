# g0Ch@
g0Ch@: A simple terminal based chat written in go (unfortunately github does not allow @ in the repo name thats why the repo is called `g0Chat` and not `g0Ch@` :/)

The chat is a linux only application, it'll propably *not* run on Windows and *maybe* on Mac. Same for the server which might also run on Windows and Mac, but there's no guarantee.

## How to build
Just execute the build script via `sh compile.sh`. Windows and Mac user may do this manually.

## How to start the server
Simply execute the `g0Ch@_server` file and the server will start on port 10000 by default (for parameters s. "Server Parameter" below). When it doesn't, check the output of it (so better run the server in a terminal ;) ).

## How to start the client
Simply execute the `g0Ch@_client` file, choose on username and enter the server data (IP, port).

## Parameter
There'er some parameter like `-l`, that allows you to specify the amount of messages that are stored locally (and that are also displayed) and `-h, --help` that shows the help page.

All parameters can be used with the syntax [parameter]=[value], e.g. `--username=Hugo` and there's also a short term for those long parameters (it's always the first letter, so `-u=Hugo` would be equivalent to the one before).
### Client Parameter
* `-u, --username : `The username/nickname.
* `-i, --ip       : `The IP of the g0Ch@ server.
* `-p, --port     : `The port of the g0Ch@ server (usually 10000).
* `-l, --limit    : `Value for the size of the message buffer (how many messages are stored). The default is 50.
* `-c, --channel  : `The channel you want to talk in.
* `-h, --help     : `Shows this kind of list.

The `-l` or `--limit` can also be used without value. In this case you will be asked for a value before you can enter a chatroom. If the parameter is not given, the default value of 50 will be used.

### Server Parameter
* `-p, --port     :`Defines a custom port. The default one (used by not using this parameter) is 1000.0

### How to chat
Seriously?^^ Get a live lol

### How to leave the chat
Simply enter `exit` as message and everything will be fine. Cancelling the chat with `Ctrl+C` may have unwanted effects on your terminal (nothing to worry about, so just try it out to see what happens --> s. below).

## Problems
### Terminal hacking
It's not a bug or something, it's just not fancy. At the moment I have to kind of hack the terminal with `stty` to be able to grab typed characters before the user presses enter.

Let's say a user writes a sentence, while doing so, a message comes in, the screen becomes cleared and the message is printed onto the screen. Normally the users input is gone now. It's buffered (so available for the internal system) but not visible anymore.

To prevent this, I have to disable buffering, grab characters directly when they are typed and display them manually. And this is ugly I think. If there's a better method out there, please let me know or commit a pull request.
