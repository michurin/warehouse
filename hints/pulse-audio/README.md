# Setting up Logitech, Inc. Webcam C310

- Device: `ID 046d:081b Logitech, Inc. Webcam C310`
- Problem: sound corruption due to wrong sample rate
- Solution: `~/.pulse/daemon.conf`

```
default-sample-rate = 16000
alternate-sample-rate = 16000
```
