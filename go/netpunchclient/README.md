```sh
go build ./...
```

```ini
[Unit]
Description=Netpunch client
After=network.target

[Service]
Type=simple
ExecStart=/opt/netpunchclient/netpunchclient a your-public-control-node.com:9999 your-secret /usr/bin/openvpn 192.168.2.1 192.168.2.2 /opt/netpunchclient/secret.key
Restart=always
RestartSec=60

[Install]
WantedBy=multi-user.target
```

```sh
systemctl daemon-reload
```

```sh
systemctl -l status netpunchclient
systemctl start netpunchclient
systemctl -l status netpunchclient
systemctl enable netpunchclient
systemctl -l status netpunchclient
```
