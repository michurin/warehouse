# Network hole punching tool

![Network hole punching sequence diagram](sequence-diagram.png)

Try mermaid

```mermaid
sequenceDiagram
    actor A as Node A<br>192.168.0.1
    participant GA as Gateway A<br>1.1.1.1
    participant S as Server<br>200.4.4.4
    participant GB as Gateway B<br>2.2.2.2
    actor B as Node B<br>10.8.0.1

    A ->> GA: from: 192.168.0.1:1001<br>to: 200.2.2.2:1000
    GA ->> S: from: 1.1.1.1:X
    S ->> S: store 1.1.1.1:X<br>(nothing to response yet)

    B ->> GB: from: 10.8.0.1:1002<br>to: 200.2.2.2:1000
    GB ->> S: from: 2.2.2.2:Y
    S ->> S: store 2.2.2.2:Y and response "1.1.1.1:X"
    S ->> B: data: "1.1.1.1:X"
    B ->> B: start pinging 1.1.1.1:X from :1002<br>(see below)

    A ->> GA: next try
    GA ->> S: next try
    S ->> S: now we are able to send "2.2.2.2:Y" to Node A
    S ->> A: data: "2.2.2.2:Y"
    A ->> A: start pinging 2.2.2.2:Y from :1001

    B -X GA: "PING" to: "1.1.1.1:X"<br>(until we get "PONG" or "CLOSE")
    A -X GB: "PING" to: "2.2.2.2:Y"<br>(until we get "PONG" or "CLOSE")
    B ->> A: "PING"
    A ->> B: "PONG" (with retries too)
    B ->> B: stop listning on<br>first "PONG"
    B ->> A: "CLOSE" (with retries)<br>First "CLOSE" stops A
    A ->> A: stop listning and exit<br>after first "CLOSE"
    B ->> B: exit after emitting all "CLOSE"
```

## Related links

- [Peer-to-Peer Communication Across Network Address Translators](https://bford.info/pub/net/p2pnat/): fundamental work on P2P drilling
- [Setup OpenVPN](https://ubuntu.com/server/docs/service-openvpn)
