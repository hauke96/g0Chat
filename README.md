# g0Ch@
g0Ch@: A simple terminal based chat written in go (unfortunately github does not allow @ in the repo name thats why the repo is called `g0Chat` and not `g0Ch@` :/)

## How to start the server
Simply execute the `g0Ch@_server` file and the server will start on port 10000. When it doesn't, check the output of it (so start in terminal).

## How to start the client
Simply execute the `g0Ch@_client` file, choose on username and enter the server data (IP, port).

### Parameter
There'er some parameter like `-l`, that allows you to specify the amount of messages that are stored locally (and that are also displayed) and `-h, --help` that shows the help page.

All parameters can be used with the syntax [parameter]=[value], e.g. `--username=Hugo` and there's also a short term for those long parameters (it's always the first letter, so `-u=Hugo` would be equivalent to the one before).
#### All Parameters
* `-u, --username : `The username/nickname.
* `-i, --ip       : `The IP of the g0Ch@ server.
* `-p, --port     : `The port of the g0Ch@ server (usually 10000).
* `-l, --limit    : `Value for the size of the message buffer (how many messages are stored). The default is 50.
* `-c, --channel  : `The channel you want to talk in.

The `-l` or `--limit` can also be used without value. In this case you will be asked for a value before you can enter a chatroom. If the parameter is not given, the default value of 50 will be used.

### How to chat
Seriously?^^ Get a live lol

### How to leave the chat
Simply enter `exit` as message and everything will be fine. Cancelling the chat with `Ctrl+C` may have unwanted effects on your terminal (nothing to worry about, so just try it out to see what happens --> s. below).

## Problems
### Can't enter anything after killing the chat client.
Just close and re-open the terminal. g0Ch@ reads the input of the user but also updates the view. This cant be done normally so the `stty` settings must be changed a bit which causes "invisible input". When you know a bette way of doing this, please create a ticket or a pull-request :)
