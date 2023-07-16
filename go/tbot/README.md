# tbot(?)

Successor of [cnbot](https://github.com/michurin/cnbot).

## What is to for

## Quick start

## Advanced configuration

### Configuration file(s)

### Secure API token

### Launch multiply bots

## Examples

### Simple scripts

#### Simplest text

#### Preformatted text

### Advanced control API

#### Advanced formatting

#### Multimedia and documents in response

### Long-running scripts

## Reference guide

### Processing updates. Overview

### Script's arguments

### Environment variables

#### Turning telegram payload to environment variables

```json
{
  "ok": true,
  "result": [
    {
      "message": {
        "caption": "Hi!",
        "chat": {
          "first_name": "Alexey",
          "id": 150000000,
          "last_name": "Michurin",
          "type": "private",
          "username": "AlexeyMichurin"
        },
        "date": 1600000000,
        "from": {
          "first_name": "Alexey",
          "id": 150000000,
          "is_bot": false,
          "language_code": "en",
          "last_name": "Michurin",
          "username": "AlexeyMichurin"
        },
        "message_id": 2222,
        "photo": [
          {
            "file_id": "aaa0",
            "file_size": 2444,
            "file_unique_id": "id0",
            "height": 90,
            "width": 90
          },
          {
            "file_id": "aaa1",
            "file_size": 4888,
            "file_unique_id": "id1",
            "height": 128,
            "width": 128
          }
        ]
      },
      "update_id": 500000000
    }
  ]
}
```

```sh
tg_message_caption=Hi!
tg_message_chat_first_name=Alexey
tg_message_chat_id=150000000
tg_message_chat_last_name=Michurin
tg_message_chat_type=private
tg_message_chat_username=AlexeyMichurin
tg_message_date=1600000000
tg_message_from_first_name=Alexey
tg_message_from_id=150000000
tg_message_from_is_bot=false
tg_message_from_language_code=en
tg_message_from_last_name=Michurin
tg_message_from_username=AlexeyMichurin
tg_message_message_id=2222
tg_message_photo=tg_message_photo_0 tg_message_photo_1
tg_message_photo_0_file_id=aaa0
tg_message_photo_0_file_size=2444
tg_message_photo_0_file_unique_id=id0
tg_message_photo_0_height=90
tg_message_photo_0_width=90
tg_message_photo_1_file_id=aaa1
tg_message_photo_1_file_size=4888
tg_message_photo_1_file_unique_id=id1
tg_message_photo_1_height=128
tg_message_photo_1_width=128
tg_update_id=500000000
```

#### Build-in variables (`x`-variables)

- `tg_x_build`
- `tg_x_ctrl_addr`
- `tg_x_to`

#### System variables

### Working directory

### Timeouts

### Concurrency

## Known issues

- It wasn't tested on MS Windows and FreeBSD

## Develop

### Main ideas

- Contract must be simple and flexible
- New features of [Telegram bot API](https://core.telegram.org/bots/api) has to be available instantly without changing of code of the bot
- Bot has to manage subprocesses: timeouts, etc
- Bot has to manage API call: [rate limits](https://core.telegram.org/bots/faq#my-bot-is-hitting-limits-how-do-i-avoid-this), etc
- Configuration must be simpler
- Code must be testable and has to be covered
- Functionality has to be observable and has to provide ability to add metrics and monitoring by adding middleware without code changing

### Embedding

### Application structure

(horrible ASCII art warning)

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
