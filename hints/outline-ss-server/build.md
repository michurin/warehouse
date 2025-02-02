```
openssl rand -base64 24

git clone https://github.com/Jigsaw-Code/outline-ss-server
```

```
services:
  - listeners:
      - type: tcp
        address: "[::]:9000"
      - type: udp
        address: "[::]:9000"
    keys:
      - id: user-0
        cipher: chacha20-ietf-poly1305
        secret: Secret0
    dialer:
      # fwmark can be used in conjunction with other Linux networking features like cgroups, network namespaces, and TC (Traffic Control) for sophisticated network management.
      # Value of 0 disables fwmark (SO_MARK) (Linux Only)
      fwmark: 0
```

```
./outline-ss-server -config cmd/outline-ss-server/config_example.yml

GOOS=linux GOARCH=386 go build ./cmd/...
scp outline-ss-server cmd/outline-ss-server/config_example.yml a@YOURHOST:~/ss
./outline-ss-server -config config_example.yml -metrics :9999

nohup ./outline-ss-server -config config.yml -metrics :9999 &

curl http://localhost:9999/metrics
```

`https://github.com/fionn/shadowsocks-url/blob/master/ss_url.py`

```
import base64
"ss://"+base64.b64encode(bytes("chacha20-ietf-poly1305:Secret0@YOURHOST:9000", "ascii")).decode()
```

```
sudo iptables -L INPUT --line-numbers -n
sudo iptables -I INPUT 9 -p tcp --dport 9000 -j ACCEPT
```

vi:ft=markdown
