# tbot(?)

Successor of [cnbot](https://github.com/michurin/cnbot).

## Main ideas

- Contract must be simpler and more flexible
- New features of [Telegram bot API](https://core.telegram.org/bots/api) has to be available without changing of code of the bot
- Bot has to manage subprocesses: timeouts, rate limits etc
- Configuration must be simpler
- Code must be testable and has to be covered
- Functionality has to be observable and has to provide ability to add metrics and monitoring by adding middleware without code changing

## Run

```
BOT_TOKEN=your_bot_token go run ./tbot/...
```
