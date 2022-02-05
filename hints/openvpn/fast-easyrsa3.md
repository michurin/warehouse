[article](https://ubuntu.com/server/docs/service-openvpn)

# CA

    ./easyrsa init-pki
    ./easyrsa build-ca nopass

# server key

    ./easyrsa gen-req owl nopass  -> pki/private/owl.key
    ./easyrsa gen-dh              -> pki/dh.pem
    ./easyrsa sign-req server owl -> pki/issued/owl.crt

    cp pki/dh.pem pki/ca.crt pki/issued/owl.crt pki/private/owl.key /etc/openvpn/server

    cd /etc/openvpn/server
    openvpn --genkey --secret ta.key

# setup

    sysctl net.ipv4.ip_forward=1

# client

    ./easyrsa gen-req mac nopass  -> pki/private/mac.key
    ./easyrsa sign-req client mac -> pki/issued/mac.crt
    https://raw.githubusercontent.com/graysky2/ovpngen/master/ovpngen
    sh ovpngen owl pki/ca.crt pki/issued/mac.crt pki/private/mac.key ta.key >mac.ovpn

    cp pki/ca.crt pki/issued/mac.crt --> client
    cp mac.ovpn                      --> client

# RUN

    openvpn /etc/openvpn/server/server.conf
