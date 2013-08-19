Still Work in Progress

sifr were supposed to be an irc framework.  
It's written in Go, and developed as a base for ircbot in mind.

Here be explained how sifr architectured:

Sifr consisted in 3 parts,
- The main program (sifr)
- irc package
- skill package

Currently the main program is just a glue between irc package and skill package, it may serves as an example on how irc package being used.

irc package is abstractized from 3 entities (struct): `Client`, `User`, and `Message`.
Client, being the main controller get initialized with method `Connect` with server address and User object passed as parameter.
Here the client will connect to the server, create channels, and run two goroutines, one for handling input then send message via channel, and the other to receive message, process it and then send the output back to server. This process stopped if input reader receives error.

The Message processing system were partially stolen from [go-ircevent](https://github.com/thoj/go-ircevent)'s callback system.

skill package enclosed all irc message handler, exposed via skill.ActivateFor function.

TODO:
- Better message parsing
- SSL support
- Command authorization
- Reconnect
- Anti-flood
- Message filter/preprocessor, run before message get passed to MessageHandler
- More skills!

Licensed under MIT/X11
