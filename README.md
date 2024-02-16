# Net-Cat

### Description

This project recreates the **NetCat in a Server-Client Architecture** that can run in a server mode on a specified port listening for incoming connections, and it can be used in client mode, trying to connect to a specified port and transmitting information to the server. The project takes the form of a TCP connected group chat support a maximum of 10 clients.

### Instructions

To initialize a server, run the go program using  `go run . ` in the terminal. This will present you with a terminal GUI specifying the server IP and port. by default the port is 8989, but you can specify an exact port (ex. `go run . 3033`)

Alternative ways to run the server include using the bash file `./TCPChat.sh` or building the binary file using `go build` which can then be executed directly based on the binary file name `./TCPChat`

To connect as a client, run netcat in another terminal using the ip and port values `nc $IP $Port`. This will provide the user with a prompt requiring them to provide their username. Once a valid useame is obtained from the client, they are then allowed onto the group chat.

### Commands

One in the group chat, the following commands are avaliable: 

`--help` : provides the user with a list of the
`--changename` : allows the user to change their name
`--quit` : exits the chat
