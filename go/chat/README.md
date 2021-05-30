# Very simple long polling chat engine

Current status: proof of concept.

Target: (i) simple library to create web chats (serve side and JS SDK) and (ii) demo.

## Quick start

Run engine:

```
go run ./cmd/chat/...
```

Playing with `curl`:

```
curl -d '{"message":{"name":"Wirt","text":"ok"}}' localhost:8080/api/publish
curl -d '{"message":"Hi!"}' localhost:8080/api/publish
curl -d '{"id":0}' localhost:8080/api/poll
{"lastID":2,"messages":[{"message":"Hi!"},{"message":{"name":"Wirt","text":"ok"}}]}
curl -d '{"id":2}' localhost:8080/api/poll
...long polling...
```

UI: Just open `http://localhost:8080/` in your browser and chat.

## Low-level contract (draft)

### Send

Just text (has to be json in future):

### Poll

Request

- id of last message that we already have

Response

- list of message (most recent first)
- id of the most recent message

## TODO

- [] Isolate library part: Message (including validation) has to be put outside
- [x] JS SDK (jQuery free)
- [] Nicknames, colors (after Message type isolation)
- [] Chat window: limit number of messages, scrolling, etc.
- [] Client- and server-side check: remove ctl chars, care about empty messages, etc
- [] Contract: let send returns registered-as id: string vector increasing id
- [] Multiply chat rooms
- [] Logging, error handling, health checking, statistics
- [] TODOs in code
- [] Setup CI

## References

- [Known Issues and Best Practices for the Use of Long Polling and Streaming in Bidirectional HTTP](https://tools.ietf.org/id/draft-loreto-http-bidirectional-07.html)
- Similar projects: [github.com/jcuga/golongpoll](https://github.com/jcuga/golongpoll)
