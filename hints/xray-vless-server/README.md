# Project

```
https://github.com/XTLS/Xray-core
```

# Build

Regular:

```
go build -o v2ray ./main
```

Android cross compilation example:

```
GOARCH='arm64' GOARM64='v8.0' GOOS='android' go build -o v2ray.arm ./main
```

For lowend i386 linux box:

```
GOOS=linux GOARCH=386 go build -o v2ray.i386 ./main
```

# Commands to configure

```
xray x25519

vless://${{UUID}}@${{IP}}:443?type=tcp&security=reality&pbk=${{PUBKEY}}&fp=chrome&sni=google.com&sid=aaaaaaaa&flow=xtls-rprx-vision#aaaaaaaa

https://habr.com/ru/articles/869340/
https://github.com/wi1dcard/v2ray-exporter
```

# Setup daemon

```
systemctl daemon-reload
systemctl list-units
```

# Read metrics

```
git clone https://github.com/fullstorydev/grpcurl.git
cd grpcurl/
go build ./cmd/...

GOOS=linux GOARCH=386 go build -o grpcurl.i386 ./cmd/...
```

Just idea:

```
for t in uplink downlink
do
  m="inbound>>>vpn>>>traffic>>>$t"
  o='xray.app.stats.command.StatsService/GetStats'
  addr='127.0.0.1:442'
  data='{"name":"'"$m"'","reset":false}'
  echo "$t $(./grpcurl.i386 -plaintext -d "$data" "$addr" "$o" | jq -r .stat.value)" |
  perl -pe 's/(\d)(\d{9}[^\d])/\1_\2/;s/(\d)(\d{6}[^\d])/\1_\2/;s/link//;s/down/dn/;s{ }{"."x(21-length($_))}e;s/ /./g'
done
```
