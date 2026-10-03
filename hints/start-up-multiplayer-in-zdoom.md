# How to start up multiplayer in zdoom

Tested on Linux gzdoom `g4.14.2-m`.

## Server side

Simplest:

```sh
gzdoom -host 2
```

Extended:

```sh
gzdoom -iwad ~/.config/gzdoom/heretic.wad -skill 1 +map E3M1 -host 2
```

## Client

Simplest:

```sh
gzdoom -join IP
```

Reach example:

```sh
DISPLAY=:0 gzdoom -iwad ~/.config/gzdoom/heretic.wad -skill 1 +map E3M1 -join 213.7.81.190
```

## Network

Server have to be reachable by **UDP** port **5029**.

`iptable` rule looks like:

```
-A INPUT -p udp --dport 5029 -j ACCEPT
```
