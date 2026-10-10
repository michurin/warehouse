Stop and disable `netctl`. Enable and start auto service

```sh
sudo netctl status wifi-f204 # enabled and online
sudo netctl stop wifi-f204
sudo netctl disable wifi-f204
sudo systemctl enable netctl-auto@wlo1.service
sudo systemctl start netctl-auto@wlo1.service
ip link # check the link is up
```

Manual switching

```sh
etctl-auto switch-to [PROFILE]
```
