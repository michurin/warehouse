# tbot(?)

Successor of [cnbot](https://github.com/michurin/cnbot).

## Main ideas

- Contract must be simple and flexible
- New features of [Telegram bot API](https://core.telegram.org/bots/api) has to be available instantly without changing of code of the bot
- Bot has to manage subprocesses: timeouts, etc
- Bot has to manage API call: [rate limits](https://core.telegram.org/bots/faq#my-bot-is-hitting-limits-how-do-i-avoid-this), etc
- Configuration must be simpler
- Code must be testable and has to be covered
- Functionality has to be observable and has to provide ability to add metrics and monitoring by adding middleware without code changing

## Run

```sh
BOT_TOKEN=your_bot_token go run ./tbot/...
```

## Application structure

```
   Telegram infrastructure
             ^                             ............. crons
        HTTP :                        HTTP :             scripts
             :                             v             any other
.=BOT================================================.   asynchronous
|            API           | HTTP server for         |
|..........................| asynchronous messaging  |
| polling for : sending    |                         |
| updates     : messages  <-- send data from req     |
`===================================================='
    |             ^    ^  send stdout     |
    |             |    `---------.        | request params
    | message     | send         |        | as command line positional args
    v data        | stdout       |        v
........................        ......................
: run script for every :        : long-running       :
: message              :        : script             :
:......................:        :....................:
```
