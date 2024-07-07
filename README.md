# Prometheus WireGuard Exporter

A Prometheus exporter for [WireGuard](https://www.wireguard.com), written in GoLang. This tool execute command `sudo /usr/bin/wg show all dump` and parse results in a format [OpenMetrics](https://github.com/OpenObservability/OpenMetrics/blob/main/specification/OpenMetrics.md) that [Prometheus](https://prometheus.io/) can understand.

### Metrics

All metrics are aveliable on `/metrics` path. As example you can use `curl 127.0.0.1:9172/metrics`

```ebnf
# HELP wireguard_keep_alive Number of keep alive seconds
# TYPE wireguard_keep_alive gauge
wireguard_keep_alive{ifname="wg0",publicKey="5CL0j1HMDxK0ftLz2ANXaZ5w2y5aOe9dV3tKF4llhlI="} 5
# HELP wireguard_last_handshake Timestamp of last handshake
# TYPE wireguard_last_handshake counter
wireguard_last_handshake{ifname="wg0",publicKey="5CL0j1HMDxK0ftLz2ANXaZ5w2y5aOe9dV3tKF4llhlI="} 2.9867196e+07
# HELP wireguard_peers_count Number of active peers
# TYPE wireguard_peers_count gauge
wireguard_peers_count{ifname="wg0"} 1
# HELP wireguard_rx_bytes Number of received bytes
# TYPE wireguard_rx_bytes counter
wireguard_rx_bytes{ifname="wg0",publicKey="5CL0j1HMDxK0ftLz2ANXaZ5w2y5aOe9dV3tKF4llhlI="} 2.9867196e+07
# HELP wireguard_tx_bytes Number of transferred bytes
# TYPE wireguard_tx_bytes counter
wireguard_tx_bytes{ifname="wg0",publicKey="5CL0j1HMDxK0ftLz2ANXaZ5w2y5aOe9dV3tKF4llhlI="} 2.746656e+06
```

For example this is how you edit your WireGuard configuration file:

```ini
[Interface]
Address = 192.168.0.8/24
PrivateKey = ******uWogB6cn/Mw********j3dZeB54C1SIn4iG2k=

[Peer]
PublicKey = 5CL0j1HMDxK0ftLz2ANXaZ5w2y5aOe9dV3tKF4llhlI=
AllowedIPs = 192.168.0.0/24,192.168.1.0/24
PersistentKeepalive = 5
Endpoint = wg0.example.com:443
```

## Flags
❗SUDO RULES❗️
If the `wireguard-exporter` is not launched as `root`, then the user must have sudo rights
```shell
someuser ALL=(root) NOPASSWD: /usr/bin/wg
```
You can simply run the `wireguard-exporter` file with standard settings on any host that has utility `wg`

You can call `wireguard-exporter -h` to print all available flags

```shell
root@myhost:~# ./wireguard-exporter -h
Usage of ./wireguard-exporter-linux-amd64:
  --log.format string
        Output format of log messages. One of: [text, json] (default "text")
  --log.level string
        Only log messages with the given severity or above. One of: [debug, info, warn, error] (default "info")
  --web.listem.address string
        Address to listen (default "0.0.0.0:9172")
```

### Build for another platform
If your platform is not in the releases, you can build it yourself, for this you only need a Go compiler. Just execute this shell script by changing `GOOS` and `GOARCH` with according [GoLang specification](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63)
```shell
GOOS=linux # REQUIRED REPLACE
GOARCH=amd64 # REQUIRED REPLACE

export GOARCH=$GOARCH
export GOOS=$GOOS
file="wireguard-exporter-$GOOS-$GOARCH"
go build -ldflags "-s -w" -o $file main.go
```
If everything went well, just transfer the file `wireguard-exporter-...-...` to your host and execute him


### Simple systemd service file

❗SUDO RULES❗️
If the `wireguard-exporter` is not launched as `root`, then the user must have sudo rights
```shell
someuser ALL=(root) NOPASSWD: /usr/bin/wg
```

Now add the exporter to the Prometheus exporters as usual. I recommend to start it as a service. It's necessary to run it as root or configure a sudo rule (if there is a non-root way to call `wg show all dump` please let me know). My systemd service file is like this one:

```ini
[Unit]
Description=Prometheus WireGuard Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=someuser
Group=someuser
Type=simple
ExecStart=/usr/local/bin/wireguard-exporter \
    --log.format text \
    --log.level info \
    --web.listem.address 0.0.0.0:9172
Restart=on-failure
RestartSec=15s

[Install]
WantedBy=multi-user.target
```

