# Demo scripts

## How to run

Add your Telegram bot token to `demo_bot.env`.

```sh
sudo docker compose up
```

## Tips and hints

Remove image after editing compose-file:

```sh
sudo docker compose images cnbot
```

```
CONTAINER                    REPOSITORY                 TAG                 IMAGE ID            SIZE
chbot-echo-example-cnbot-1   chbot-echo-example-cnbot   latest              8d79354c2325        515MB
```

```sh
sudo docker image rm -f chbot-echo-example-cnbot:latest
```

Enter to container to debug, install new packages, edit scripts etc:

```sh
sudo docker compose exec cnbot /bin/bash
```

```
root@ba1a8da9eab4:/app#
```
