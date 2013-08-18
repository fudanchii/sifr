Still Work in Progress

sifr were supposed to be an irc framework.  
It's written in Go, and developed as a base for ircbot in mind.

Here be explained how sifr architectured:

Sifr consisted in 3 parts,
- The main program (sifr)
- irc package
- skill package

Currently the main program is just a glue between irc module and skill module, and it may serves as an example on how irc module being used.

TODO:
- Better message parsing
- SSL support
- Command authorization
- More skills!

Licensed under MIT/X11
