# midi-server

Examples of how to send MIDI message events via WebSocket.
Tested in Ubuntu 20.04 LTS/MacBook Arm64 arch.

First, connect the MIDI device, use `amidi -l` to view the connected
MIDI devices, and update the MIDI deviceID in the program if necessary.

## Golang Server

Install go and portmidi, then run:

```
go run main.go
```

## websocketd with test data

Install websocketd, then run:

```bash
websocketd -port=8080 -- script -pq data/midi-events.json.mat

#The script utility makes a typescript of everything printed on your terminal.  It is useful for students who need a hardcopy record of an interactive session as proof of an assignment, as the typescript file can be printed out later with lpr(1).
```

This is useful when you don't have a MIDI device connected but want to test websocket clients.
