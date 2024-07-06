package routes

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Metrics(writer http.ResponseWriter, request *http.Request, logger *slog.Logger) {

	// Metrics descriptions
	wireguard_rx_bytes := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "wireguard_rx_bytes",
		Help: "Number of received bytes",
	}, []string{"ifname", "publicKey"})
	wireguard_tx_bytes := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "wireguard_tx_bytes",
		Help: "Number of transferred bytes",
	}, []string{"ifname", "publicKey"})
	wireguard_keep_alive := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wireguard_keep_alive",
		Help: "Number of keep alive seconds",
	}, []string{"ifname", "publicKey"})
	wireguard_last_handshake := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "wireguard_last_handshake",
		Help: "Timestamp of last handshake",
	}, []string{"ifname", "publicKey"})
	wireguard_peers_count := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wireguard_peers_count",
		Help: "Number of active peers",
	}, []string{"ifname"})

	// Scrape stats from CLI
	peers, err := Wireguard{}.Show()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Peers counter
	peers_count := make(map[string]int)

	// Handle metrics
	for _, peer := range peers {

		// Count peers
		_, exists := peers_count[peer.Ifname]
		if exists {
			peers_count[peer.Ifname] += 1
		} else {
			peers_count[peer.Ifname] = 1
		}

		wireguard_tx_bytes.With(prometheus.Labels{
			"ifname":    peer.Ifname,
			"publicKey": peer.Publickey,
		}).Add(float64(peer.Txbytes))
		wireguard_rx_bytes.With(prometheus.Labels{
			"ifname":    peer.Ifname,
			"publicKey": peer.Publickey,
		}).Add(float64(peer.Rxbytes))
		wireguard_keep_alive.With(prometheus.Labels{
			"ifname":    peer.Ifname,
			"publicKey": peer.Publickey,
		}).Add(float64(peer.Keepalive))
		wireguard_last_handshake.With(prometheus.Labels{
			"ifname":    peer.Ifname,
			"publicKey": peer.Publickey,
		}).Add(float64(peer.Lasthandshake))
	}

	// Register wireguard_peers_count
	for key, value := range peers_count {
		wireguard_peers_count.With(prometheus.Labels{
			"ifname": key,
		}).Set(float64(value))
	}

	prometheus.Register(wireguard_rx_bytes)
	prometheus.Register(wireguard_tx_bytes)
	prometheus.Register(wireguard_keep_alive)
	prometheus.Register(wireguard_last_handshake)
	prometheus.Register(wireguard_peers_count)

	// Logging request
	logger.Info("Successfully handled", "url", request.RequestURI, "method", request.Method, "remote", request.RemoteAddr, "code", 200)

	// Publish metrics via promhttp
	promhttp.Handler().ServeHTTP(writer, request)
}

type Wireguard struct{}

// Calls the command "wg show all dump" and parses the output
// Returns the Peer structure
func (w Wireguard) Show() ([]Peer, error) {

	// Collect data by exec "wg show all dump"
	cmd := exec.Command("sudo", "/usr/bin/wg", "show", "all", "dump")
	output, _ := cmd.CombinedOutput()
	rc := cmd.ProcessState.ExitCode()
	if rc != 0 {
		return nil, fmt.Errorf("Fatal error failed execute command 'wg show all dump'. Return code is " + strconv.Itoa(rc))
	}

	// Parse result
	lines := strings.Split(string(output), "\n") // Split by \n (new line)

	var result []Peer

	// Foreach all lines
	for i := range lines {
		substrings := strings.Fields(lines[i]) // Split by space
		if len(substrings) > 5 {

			// Convert values
			lastHandShake, err := strconv.Atoi(substrings[6])
			if err != nil {
				fmt.Println(err.Error())
			}
			rxBytes, _ := strconv.Atoi(substrings[6])
			txBytes, _ := strconv.Atoi(substrings[7])
			keepAlive, _ := strconv.Atoi(substrings[8])

			result = append(result, Peer{
				Ifname:        substrings[0],
				Publickey:     substrings[1],
				Endpoint:      substrings[3],
				Allowedhosts:  substrings[4],
				Lasthandshake: lastHandShake,
				Rxbytes:       rxBytes,
				Txbytes:       txBytes,
				Keepalive:     keepAlive,
			})
		}
	}
	return result, nil
}

type Peer struct {
	Ifname        string
	Publickey     string
	Endpoint      string
	Allowedhosts  string
	Lasthandshake int
	Txbytes       int
	Rxbytes       int
	Keepalive     int
}
