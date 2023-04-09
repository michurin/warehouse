# Prometheus exporter

## Project status

It is not ready for use solution. However, it is very easy to steal parts of code from it and fit to your environment.

Beware, some code is good for tiny projects only. Somewhere it's hard to test or looks like dirty hack.
Somewhere it even smells like unsafe and goroutine leaking prone, see comments.

## Supported metrics overview

### Load average

```
app_loadavg{loadavg="1"} 0.2
app_loadavg{loadavg="10"} 0.41
app_loadavg{loadavg="5"} 0.25
```

### CPU utilization details: total and per core

```
app_cpu{cpu="cpu",mx="guest"} 0
app_cpu{cpu="cpu",mx="guest_nice"} 0
app_cpu{cpu="cpu",mx="idle"} 1.78813829e+08
app_cpu{cpu="cpu",mx="iowait"} 121112
app_cpu{cpu="cpu",mx="irq"} 126654
app_cpu{cpu="cpu",mx="nice"} 349406
app_cpu{cpu="cpu",mx="softirq"} 32175
app_cpu{cpu="cpu",mx="steal"} 0
app_cpu{cpu="cpu",mx="system"} 1.446909e+06
app_cpu{cpu="cpu",mx="user"} 1.60626e+06
app_cpu{cpu="cpu0",mx="guest"} 0
app_cpu{cpu="cpu0",mx="guest_nice"} 0
...
```

### Sensors

```
app_temperature{hw="acpitz",term="null"} 27.8
app_temperature{hw="coretemp",term="Core_0"} 38
app_temperature{hw="coretemp",term="Core_1"} 39
app_temperature{hw="coretemp",term="Core_2"} 40
app_temperature{hw="coretemp",term="Core_3"} 39
app_temperature{hw="coretemp",term="Core_4"} 39
app_temperature{hw="coretemp",term="Core_5"} 39
app_temperature{hw="coretemp",term="Core_6"} 40
app_temperature{hw="coretemp",term="Core_7"} 39
app_temperature{hw="coretemp",term="Package_id_0"} 42
app_temperature{hw="iwlwifi_1",term="null"} 44
app_temperature{hw="nvme",term="Composite"} 37.85
```

### Block devices

```
app_block{device="nvme0n1",m="ms",mx="discard"} 0
app_block{device="nvme0n1",m="ms",mx="flush"} 34545
app_block{device="nvme0n1",m="ms",mx="in_queue"} 5.525806e+06
app_block{device="nvme0n1",m="ms",mx="io"} 1.402814e+06
app_block{device="nvme0n1",m="ms",mx="read"} 61212
app_block{device="nvme0n1",m="ms",mx="write"} 5.430048e+06
app_block{device="nvme0n1",m="reqs",mx="discard"} 0
app_block{device="nvme0n1",m="reqs",mx="discard_merges"} 0
app_block{device="nvme0n1",m="reqs",mx="flush"} 48885
app_block{device="nvme0n1",m="reqs",mx="read"} 144510
app_block{device="nvme0n1",m="reqs",mx="read_merges"} 60609
app_block{device="nvme0n1",m="reqs",mx="write"} 873222
app_block{device="nvme0n1",m="reqs",mx="write_merges"} 334028
app_block{device="nvme0n1",m="reqs_gauge",mx="in_flight"} 0
app_block{device="nvme0n1",m="sectors",mx="discard"} 0
app_block{device="nvme0n1",m="sectors",mx="read"} 7.093052e+06
app_block{device="nvme0n1",m="sectors",mx="write"} 2.2400963e+07
app_block{device="sda",m="ms",mx="discard"} 0
app_block{device="sda",m="ms",mx="flush"} 0
app_block{device="sda",m="ms",mx="in_queue"} 138190
app_block{device="sda",m="ms",mx="io"} 126190
app_block{device="sda",m="ms",mx="read"} 88202
app_block{device="sda",m="ms",mx="write"} 49987
app_block{device="sda",m="reqs",mx="discard"} 0
app_block{device="sda",m="reqs",mx="discard_merges"} 0
app_block{device="sda",m="reqs",mx="flush"} 0
app_block{device="sda",m="reqs",mx="read"} 22885
app_block{device="sda",m="reqs",mx="read_merges"} 3
app_block{device="sda",m="reqs",mx="write"} 5256
app_block{device="sda",m="reqs",mx="write_merges"} 237
app_block{device="sda",m="reqs_gauge",mx="in_flight"} 0
app_block{device="sda",m="sectors",mx="discard"} 0
app_block{device="sda",m="sectors",mx="read"} 1.915504e+06
app_block{device="sda",m="sectors",mx="write"} 7.315008e+06
```

### Local network details per interface

```
app_network{if="lo",metric="collisions"} 0
app_network{if="lo",metric="multicast"} 0
app_network{if="lo",metric="rx_bytes"} 7.821599142e+09
app_network{if="lo",metric="rx_compressed"} 0
app_network{if="lo",metric="rx_crc_errors"} 0
app_network{if="lo",metric="rx_dropped"} 0
app_network{if="lo",metric="rx_errors"} 0
app_network{if="lo",metric="rx_fifo_errors"} 0
app_network{if="lo",metric="rx_frame_errors"} 0
app_network{if="lo",metric="rx_length_errors"} 0
app_network{if="lo",metric="rx_missed_errors"} 0
app_network{if="lo",metric="rx_nohandler"} 0
app_network{if="lo",metric="rx_over_errors"} 0
app_network{if="lo",metric="rx_packets"} 2.487579e+06
app_network{if="lo",metric="tx_aborted_errors"} 0
app_network{if="lo",metric="tx_bytes"} 7.821599142e+09
app_network{if="lo",metric="tx_carrier_errors"} 0
app_network{if="lo",metric="tx_compressed"} 0
app_network{if="lo",metric="tx_dropped"} 0
app_network{if="lo",metric="tx_errors"} 0
app_network{if="lo",metric="tx_fifo_errors"} 0
app_network{if="lo",metric="tx_heartbeat_errors"} 0
app_network{if="lo",metric="tx_packets"} 2.487579e+06
app_network{if="lo",metric="tx_window_errors"} 0
app_network{if="wlo1",metric="collisions"} 0
app_network{if="wlo1",metric="multicast"} 0
app_network{if="wlo1",metric="rx_bytes"} 2.560707471e+09
app_network{if="wlo1",metric="rx_compressed"} 0
app_network{if="wlo1",metric="rx_crc_errors"} 0
...
```

### Processes and kernel context switch count

```
app_proc{mx="ctxt"} 1.97630333e+08
app_proc{mx="processes"} 50954
app_proc{mx="procs_blocked"} 0
app_proc{mx="procs_running"} 2
```

### My router statistics (SNMP)

```
app_router{if="bridge",octets="in"} 2.914962302e+09
app_router{if="bridge",octets="out"} 4.194052036e+09
app_router{if="ether1",octets="in"} 4.18180629e+08
app_router{if="ether1",octets="out"} 2.880240305e+09
app_router{if="wlan1",octets="in"} 2.941688676e+09
app_router{if="wlan1",octets="out"} 4.030649641e+09
app_router{if="wlan2",octets="in"} 2.453019438e+09
app_router{if="wlan2",octets="out"} 2.463335355e+09
```

### Checking infrastructure of my internet provider using TPC probing (draft)

```
app_network_reachable{point="hop_2"} 1
app_network_reachable{point="hop_3"} 1
```

### Hardware interrupts (very custom)

You can find legend in you `/proc/interrupts`

```
app_interrupts{mx="CAL"} 1.655089e+06
app_interrupts{mx="DFR"} 0
app_interrupts{mx="ERR"} 0
app_interrupts{mx="IWI"} 337612
app_interrupts{mx="LOC"} 3.62964e+06
app_interrupts{mx="MCE"} 0
app_interrupts{mx="MCP"} 383
app_interrupts{mx="MIS"} 0
app_interrupts{mx="NMI"} 73
app_interrupts{mx="NPI"} 0
app_interrupts{mx="PCI-MSI-acpi"} 79585
app_interrupts{mx="PCI-MSI-ahci"} 0
app_interrupts{mx="PCI-MSI-cascade"} 0
app_interrupts{mx="PCI-MSI-i915"} 507751
app_interrupts{mx="PCI-MSI-mei_me"} 62
app_interrupts{mx="PCI-MSI-pcie-dpc"} 0
app_interrupts{mx="PCI-MSI-rtc0"} 0
app_interrupts{mx="PCI-MSI-snd_hda_intel"} 0
app_interrupts{mx="PCI-MSI-xhci_hcd"} 109740
app_interrupts{mx="PCI-MSIX-iwlwifi"} 79037
app_interrupts{mx="PCI-MSIX-nvme"} 60572
app_interrupts{mx="PIN"} 0
app_interrupts{mx="PIW"} 0
app_interrupts{mx="PMI"} 73
app_interrupts{mx="RES"} 197993
app_interrupts{mx="RTR"} 15
app_interrupts{mx="SPU"} 0
app_interrupts{mx="THR"} 0
app_interrupts{mx="TLB"} 1.322019e+06
app_interrupts{mx="TRM"} 0
```

### Mikrotik router connections information by client

```
app_mikrotik_clients{client="A31_2g",dir="null",mx="signal_to_noise"} 63
app_mikrotik_clients{client="A31_2g",dir="null",mx="strength"} -56
app_mikrotik_clients{client="A31_2g",dir="null",mx="uptime"} 3.5947e+06
app_mikrotik_clients{client="A31_2g",dir="rx",mx="bit_per_sec"} 6.5e+07
app_mikrotik_clients{client="A31_2g",dir="rx",mx="octets"} 8.9289804e+07
app_mikrotik_clients{client="A31_2g",dir="rx",mx="pkt"} 572956
app_mikrotik_clients{client="A31_2g",dir="tx",mx="bit_per_sec"} 7.22e+07
app_mikrotik_clients{client="A31_2g",dir="tx",mx="octets"} 2.713791687e+09
app_mikrotik_clients{client="A31_2g",dir="tx",mx="pkt"} 2.112439e+06
app_mikrotik_clients{client="Taurus",dir="null",mx="signal_to_noise"} 43
app_mikrotik_clients{client="Taurus",dir="null",mx="strength"} -74
app_mikrotik_clients{client="Taurus",dir="null",mx="uptime"} 6.7899e+06
app_mikrotik_clients{client="Taurus",dir="rx",mx="bit_per_sec"} 4.333e+08
app_mikrotik_clients{client="Taurus",dir="rx",mx="octets"} 3.06944188e+08
app_mikrotik_clients{client="Taurus",dir="rx",mx="pkt"} 2.3371426e+07
app_mikrotik_clients{client="Taurus",dir="tx",mx="bit_per_sec"} 3.25e+08
app_mikrotik_clients{client="Taurus",dir="tx",mx="octets"} 2.540330108e+09
app_mikrotik_clients{client="Taurus",dir="tx",mx="pkt"} 6.8735705e+07
```

## Setup systemd

Minimal `/etc/systemd/system/home-box-prometheus-exporter.service`

```ini
[Unit]
Description=Home box Prometheus exporter
After=network.target

[Service]
Type=simple
Restart=always
RestartSec=1
User=nobody
ExecStart=/opt/home-box-prometheus-exporter/home-box-prometheus-exporter

[Install]
WantedBy=multi-user.target
```

```sh
sudo cp home-box-prometheus-exporter.service /etc/systemd/system
systemctl daemon-reload
sudo systemctl start home-box-prometheus-exporter
sudo systemctl status home-box-prometheus-exporter
```

```sh
sudo systemctl enable home-box-prometheus-exporter
sudo systemctl enable prometheus
```
