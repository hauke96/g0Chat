# g0Ch@
g0Ch@: A simple terminal based chat written in go (unfortunately github does not allow @ in the repo name thats why the repo is called `g0Chat` and not `g0Ch@` :/)

## How to start the server
Simply execute the `g0Ch@_server` file and the server will start on port 10000. When it doesn't, check the output of it (so start in terminal).

## How to start the client
Simply execute the `g0Ch@_client` file, choose on username and enter the server data (IP, port).

### Parameter
There'er some parameter like `-b`, that allows you to specify the amount of messages that are stored locally (and that are also displayed) and `-h, --help` that shows the help page.

### How to chat
Get a live

### How to leave the chat
Simply enter `exit` as message and everything will be fine.

## Problems
### Can't enter anything after killing the chat client.
Just close and re-open the terminal. g0Ch@ reads the input of the user but also updates the view. This cant be done normally so the `stty` settings must be changed a bit which causes "invisible input". When you know a bette way of doing this, please create a ticket or a pull-request :)
