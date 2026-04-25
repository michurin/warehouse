- `ssh-reverse-tunnel.service` — `/etc/systemd/system/ssh-reverse-tunnel.service`
- `cpupower.service` — CPU frequency control

### Clear

```sh
sudo journalctl --flush --rotate
sudo journalctl --vacuum-time=1d
```

### Run with limits

```
systemd-run --user --scope -p MemoryMax=500M -p MemorySwapMax=0G google-chrome-stable
```