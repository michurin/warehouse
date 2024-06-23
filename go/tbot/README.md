# cnbot

The goal of this project is to provide a way
to alive Telegram bots by scripting that
even simpler than CGI scripts.
All you need to write is a script (on any language)
that is complying with extremely simple contract.

All interactions are based on `stdout` stream, arguments and environment variables.

The engine recognize multimedia and images and cares about concurrency and races.

It also provides simple API for asynchronous messaging from crons and such things.

It manages tasks (subprocesses), control timeouts, send signals and provides abilities to
run long-running tasks like long image/video conversions and/or downloading.

One engine is able to manage several different bots.

## What is to for

This bot engine has proven itself in monitoring, alerting system monitoring and managing tasks.

## Quick start

### Run simplest one-line bot

Let's start our firs bot. First things first, you need to create bot and get it's token.
It is free, just follow [instructions](https://core.telegram.org/bots#how-do-i-create-a-bot).

Just run one command to invoke the simplest bot:

```sh
tb_mybot_token='TOKEN' tb_mybot_script=/usr/bin/echo tb_mybot_long_running_script=/usr/bin/echo tb_mybot_ctrl_addr=:9999 go run ./cmd/cnbot/...
```

Don't worry, we will use configuration file further. The engine is able to use both files and direct environment variables.

- `tb_YOURBOTNAME_token` is a token your are given: `digits:long_string`
- `tb_YOURBOTNAME_script` is a command to run. We use the standard system command `echo`. I can be located elsewhere in your system. Try to say `whereis echo` to fine it
- `tb_YOURBOTNAME_long_running_script` let it be the same command. We consider it later
- `tb_YOURBOTNAME_ctrl_addr` we consider it soon

Run this command with correct variables and try to say something to you bot. You will be echoed by it.

### Put your configuration to file

You may as well put your configuration into env-file. The format of file is literally the same as `systemd` use.
So you are able to load it in `systemd` files as well. For example:

```sh
# let's name it config.env
tb_mybot_token='TOKEN'
tb_mybot_script=/usr/bin/echo
tb_mybot_long_running_script=/usr/bin/echo
tb_mybot_ctrl_addr=:9999
```

Now just start bot like this:

```sh
go run ./cmd/cnbot/... config.env
```

### Drive multiply bots

It is possible to drive more than one bots at the same time. Just use different prefixes for variables:

```sh
tb_mybot_token='TOKEN'
tb_mybot_script=/usr/bin/echo
...
tb_thenextbot_token='TOKEN'
tb_thenextbot_script=/usr/bin/echo
...
```

### Your first script (finding out your UserID)

Let's look at the script, that shows its arguments and environment variables:

```sh
#!/bin/sh

echo "Args: $@"
echo "Environment:"
env | sort | grep tg_
```

Name it `mybot.sh` and mention it in configuration variable `tb_mybot_script=./mybot.sh`. Restart the bot and say to it `Hellp bot!`.
It will reply to you something like that:

```
Args: hello bot!
Environment:
tg_message_chat_first_name=Alexey
tg_message_chat_id=153333328
tg_message_chat_last_name=Michurin
tg_message_chat_type=private
tg_message_chat_username=AlexeyMichurin
tg_message_date=1717171717
tg_message_from_first_name=Alexey
tg_message_from_id=153333328
tg_message_from_is_bot=false
tg_message_from_language_code=en
tg_message_from_last_name=Michurin
tg_message_from_username=AlexeyMichurin
tg_message_message_id=4554
tg_message_text=Hello bot!
tg_update_id=513333387
tg_x_build=development (devel)
tg_x_ctrl_addr=:9999
```

You can see that your message has been put to arguments in convenient normalized form, and you have a bunch of useful variables
with additional information. We will consider them further. At this point we just figure out then our user id is `tg_message_from_id=153333328`.
We will use this information very soon.

### Asynchronous messaging

You are free to send messages from anywhere: from cron jobs, from init scripts... Try it just from command line:

```sh
curl -qs http://localhost:9999/?to=153333328 -d 'OK!'
```

If you bot is running, you will obtain the message `OK!` in you Telegram client.

```
╭──────────╮
│ OK!      │
╰──────────╯
```

Do not forget to use *your* user id from previous section.

It makes sense what variable `tb_mybot_ctrl_addr=:9999` is for. It defines a control interface for external interactions with bot engine.

### Call arbitrary Telegram API methods

You can call whatever method you want. Full list of methods can be found in the
[official Telegram bot API documentation](https://core.telegram.org/bots/api).

For example, you can obtain information about your bot
(using method [getMe](https://core.telegram.org/bots/api#getme)):

```sh
curl -qs http://localhost:9999/method/getMe | jq
```

The response will look like this:

```json
{
  "ok": true,
  "result": {
    "id": 223333386,
    "is_bot": true,
    "first_name": "Your Bot",
    "username": "your_bot",
    "can_join_groups": true,
    "can_read_all_group_messages": false,
    "supports_inline_queries": false,
    "can_connect_to_business": false
  }
}
```

It enables you to send extended messages. For example, you can send a message with buttons
(method [sendMessage](https://core.telegram.org/bots/api#sendmessage)):

```sh
curl -qs http://localhost:9999/sendMessage -F chat_id=153333328 -F text='Select search engine' -F reply_markup='{"inline_keyboard":[[{"text":"Google","url":"https://www.google.com/"}, {"text":"DuckDuckGo","url":"https://duckduckgo.com/"}]]}'
```

You will receive message with two clickable buttons:

```
╭───────────────────────────╮
│ Select search engine      │
├─────────────┬─────────────┤
│ Google     ↗│ DuckDuckGo ↗│
╰─────────────┴─────────────╯
```

Do not forget to change `user_id`.

> Please note that you can use any prefixes in URLs.
> URLs `http://localhost:9999/sendMessage` and `http://localhost:9999/ANITHING/sendMessage` are equal.
> It allows you to put engine's API behind prefix.

### Lets play with images (synchronously and asynchronously; our second script)

Bot recognizes media type of input. It will send text:

```sh
echo 'Hello!' | curl -qs http://localhost:9999/?to=153333328 --data-binary '@-'
```

However, it will send you image:

```sh
curl -qs https://github.githubassets.com/favicons/favicon.png | curl -qs http://localhost:9999/?to=153333328 --data-binary '@-'
```

> Please use the `--data-binary` option for binary data. Option `-d` corrupts EOLs.

### Formatted text and putting all together: script, that considers commands

Let's extend our `mybot.sh` like that:

```sh
#!/bin/sh

LOG=/dev/null

CTRL="http://localhost$tg_x_ctrl_addr"
FROM="$tg_message_from_id"

crl() {
    url="http://localhost$tg_x_ctrl_addr/$1"
    shift
    echo "====== calling $url $@" >>$LOG
    curl -qs "$url" "$@" >>$LOG 2>&1
    echo >>$LOG
    echo >>$LOG
}

case "$1" in
    debug)
        echo '%!PRE'
        echo "Args: $@"
        echo "Environment:"
        env | sort | grep tg_
        echo "CTRL=$CTRL"
        echo "FROM=$FROM"
        echo "LOG=$LOG"
        ;;
    about)
        echo '%!PRE'
        curl -qs http://localhost:9999/method/getMe | jq
        ;;
    two)
        crl ?to=$FROM -d 'OK ONE!'
        crl ?to=$FROM -d 'OK TWO!!'
        ;;
    buttons)
        crl sendMessage \
            -F chat_id=$FROM \
            -F text='Select search engine' \
            -F reply_markup='{"inline_keyboard":[[{"text":"Google","url":"https://www.google.com/"},{"text":"DuckDuckGo","url":"https://duckduckgo.com/"}]]}'
        ;;
    image)
        curl -qs https://github.com/fluidicon.png |
            crl ?to=$FROM --data-binary '@-'
        ;;
    reaction)
        crl setMessageReaction -F chat_id=$FROM -F message_id=$tg_message_message_id -F reaction='[{"type":"emoji","emoji":"👾"}]'
        ;;
    help)
        crl sendMessage -F chat_id=$FROM -F text='
Known commands:

- `debug` — show args, environment and vars
- `about` — reslut of getMe
- `two` — one request, two responses
- `buttons` — message with buttons
- `image` — show image
- `reaction` — show reaction
' -F parse_mode=Markdown
        ;;
    *)
        crl sendMessage -F chat_id=$FROM -F text='Invalid command. Say `help`.' -F parse_mode=Markdown
        ;;
esac
```

## Advanced topics

### Process management: concurrency, timeouts, signals, long-running tasks

### Uploading and downloading

### Environment details

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

## Running

The process itself does not try to be immortal. It dies on fatal issues that can not be solved by process itself. Like network problems.
It is believed that the process will be restart by systemd or stuff like that according the proper way with timeouts, logging, notifications, alerting.

TODO: example of systemd file

## Known issues

- It wasn't tested on MS Windows and FreeBSD

## Develop

### Main ideas

- Contract must be simple and flexible
- New features of [Telegram bot API](https://core.telegram.org/bots/api) has to be available instantly without changing of code of the bot
- Bot has to manage subprocesses: timeouts, etc
- Bot has to manage API call: [rate limits](https://core.telegram.org/bots/faq#my-bot-is-hitting-limits-how-do-i-avoid-this), etc
- Configuration must be simple
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
