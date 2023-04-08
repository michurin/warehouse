package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var log = fmt.Println

func ls(p string) [][2]string {
	lst, err := os.ReadDir(p)
	if err != nil {
		log(fmt.Sprintf("cannot read path: %s: %s", p, err.Error()))
		return nil
	}
	res := make([][2]string, len(lst))
	for i, v := range lst {
		n := v.Name()
		res[i] = [2]string{path.Join(p, n), n}
	}
	return res
}

func readFloat(f string) (float64, error) {
	v, err := os.ReadFile(f)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseFloat(strings.TrimSpace(string(v)), 64)
	if err != nil {
		return 0, err
	}
	return float64(n), nil
}

func readLabel(f string) (string, error) {
	v, err := os.ReadFile(f)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(v)), nil
}

func safeGaugeSet(opts prometheus.GaugeOpts, v float64) {
	m := prometheus.NewGauge(opts)
	err := prometheus.DefaultRegisterer.Register(m)
	if err != nil {
		are := &prometheus.AlreadyRegisteredError{}
		if errors.As(err, are) {
			m = are.ExistingCollector.(prometheus.Gauge)
		} else {
			log(fmt.Sprintf("ERROR: Reg new metric: %s: %s", m.Desc().String(), err.Error()))
			panic(err)
		}
	} else {
		log(fmt.Sprintf("Reg new metric: %s", m.Desc().String()))
	}
	m.Set(v)
}

func safeLabels(x map[string]string) map[string]string {
	res := map[string]string{}
	for k, v := range x {
		if v == "" {
			v = "null"
		}
		res[k] = strings.ReplaceAll(v, " ", "_")
	}
	return res
}

func updateLocalNetwork() {
	for _, ifs := range ls("/sys/class/net") {
		for _, f := range ls(path.Join(ifs[0], "statistics")) {
			v, err := readFloat(f[0])
			if err != nil {
				log(f[0], err)
				continue
			}
			safeGaugeSet(prometheus.GaugeOpts{
				Name:        "app_network",
				ConstLabels: map[string]string{"if": ifs[1], "metric": f[1]},
			}, v)
		}
	}
}

func updateSensors() {
	for _, hw := range ls("/sys/class/hwmon") {
		name := ""
		input := map[string]float64{}
		label := map[string]string{}
		all := map[string]struct{}{}
	LOOP:
		for _, h := range ls(hw[0]) {
			n := h[1]
			switch {
			case n == "name":
				v, err := readLabel(h[0])
				if err != nil {
					log("SKIP LOOP error:", h[0], err)
					continue LOOP
				}
				name = v
			case strings.HasPrefix(n, "temp") && strings.HasSuffix(n, "_input"):
				v, err := readFloat(h[0])
				if err != nil {
					log("SKIP LOOP error:", h[0], err)
					continue LOOP
				}
				if v == 0 { // somehow temperature in tempXX_input can be zero
					log("SKIP LOOP value:", v, v == 0, h[0])
					continue LOOP
				}
				k := n[4 : len(n)-6]
				all[k] = struct{}{}
				input[k] = v / 1000.
			case strings.HasPrefix(n, "temp") && strings.HasSuffix(n, "_label"):
				v, err := readLabel(h[0])
				if err != nil {
					log("SKIP LOOP error:", h[0], err)
					continue LOOP
				}
				k := n[4 : len(n)-6]
				all[k] = struct{}{}
				label[k] = v
			}
		}
		for k := range all {
			f, ok := input[k]
			if !ok {
				log("SKIP key not found in input:", name, k)
				continue
			}
			l, ok := label[k]
			if !ok {
				log("FALLBACK key not found in labels:", name, k)
			}
			safeGaugeSet(prometheus.GaugeOpts{
				Name:        "app_temperature",
				ConstLabels: safeLabels(map[string]string{"hw": name, "term": l}),
			}, f)
		}
	}
}

func updateBlocks() {
	// https://docs.kernel.org/block/stat.html
	for _, hw := range ls("/sys/block") {
		c, err := readLabel(path.Join(hw[0], "stat"))
		if err != nil {
			log("SKIP block device", hw[0], err)
			continue
		}
		f := strings.Fields(c)
		for i, g := range [][2]string{
			{"read", "reqs"},            // number of read I/Os processed
			{"read_merges", "reqs"},     // number of read I/Os merged with in-queue I/O
			{"read", "sectors"},         // number of sectors read
			{"read", "ms"},              // total wait time for read requests
			{"write", "reqs"},           // number of write I/Os processed
			{"write_merges", "reqs"},    // number of write I/Os merged with in-queue I/O
			{"write", "sectors"},        // number of sectors write
			{"write", "ms"},             // total wait time for write requests
			{"in_flight", "reqs_gauge"}, // number of I/Os currently in flight
			{"io", "ms"},                // total time this block device has been active
			{"in_queue", "ms"},          // total wait time for all requests
			{"discard", "reqs"},         // number of discard I/Os processed
			{"discard_merges", "reqs"},  // number of discard I/Os merged with in-queue I/O
			{"discard", "sectors"},      // number of sectors discard
			{"discard", "ms"},           // total wait time for discard requests
			{"flush", "reqs"},           // number of flush I/Os processed
			{"flush", "ms"},             // total wait time for flush requests
		} {
			n, err := strconv.ParseFloat(strings.TrimSpace(f[i]), 64)
			if err != nil {
				log("SKIP block device metric error:", hw[0], g[0], g[1], f[i], err)
				continue
			}
			safeGaugeSet(prometheus.GaugeOpts{
				Name:        "app_block",
				ConstLabels: safeLabels(map[string]string{"device": hw[1], "mx": g[0], "m": g[1]}),
			}, float64(n))
		}
	}
}

func updateCPU() {
	// https://docs.kernel.org/filesystems/proc.html#miscellaneous-kernel-statistics-in-proc-stat
	c, err := readLabel("/proc/stat")
	if err != nil {
		log("cpu stat read error:", err)
	}
	for _, l := range strings.Split(c, "\n") {
		f := strings.Fields(l)
		h := f[0]  // head
		t := f[1:] // tail
		switch {
		case strings.HasPrefix(h, "cpu"):
			for i, mx := range []string{"user", "nice", "system", "idle", "iowait", "irq", "softirq", "steal", "guest", "guest_nice"} {
				v := t[i]
				f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
				if err != nil {
					log("SKIP cpu metric error:", h, mx, v, err)
					continue
				}
				safeGaugeSet(prometheus.GaugeOpts{
					Name:        "app_cpu",
					ConstLabels: safeLabels(map[string]string{"cpu": h, "mx": mx}),
				}, float64(f))
			}
		case h == "processes" || h == "procs_running" || h == "procs_blocked" || h == "ctxt":
			v := t[0]
			f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
			if err != nil {
				log("SKIP proc metric error:", h, v, err)
				break
			}
			safeGaugeSet(prometheus.GaugeOpts{
				Name:        "app_proc",
				ConstLabels: safeLabels(map[string]string{"mx": h}),
			}, float64(f))
		}
	}
}

func updateLoadAvg() {
	f, err := readLabel("/proc/loadavg")
	if err != nil {
		log("cpu stat read error:", err)
	}
	c := strings.Fields(f)
	for i, m := range []string{"1", "5", "10"} {
		n, err := strconv.ParseFloat(c[i], 64)
		if err != nil {
			log("SKIP block device metric error:", m, c[i], err)
		}
		safeGaugeSet(prometheus.GaugeOpts{
			Name:        "app_loadavg",
			ConstLabels: safeLabels(map[string]string{"loadavg": m}),
		}, float64(n))
	}
}

func readSnmp(root string) map[string]any {
	out, err := gosnmp.Default.BulkWalkAll(root)
	if err != nil {
		log("snmp err:", err)
		return nil
	}
	res := map[string]any{}
	for _, v := range out {
		res[strings.Trim(v.Name[len(root):], ".")] = v.Value
	}
	return res
}

func updateRouterNetwork() {
	names := readSnmp(".1.3.6.1.2.1.2.2.1.2")
	ins := readSnmp(".1.3.6.1.2.1.2.2.1.10")
	outs := readSnmp(".1.3.6.1.2.1.2.2.1.16")
	dir := []string{"in", "out"}
	for i, source := range []map[string]any{ins, outs} {
		for k, v := range source {
			n, ok := v.(uint)
			if !ok {
				log(fmt.Sprintf("SKIP cannot cast %[1]s: %[2]T: %[2]v", k, v))
				continue
			}
			b, ok := names[k].([]byte)
			if !ok {
				log(fmt.Sprintf("SKIP cannot cast %[1]s: %[2]T: %[2]v", k, names[k]))
				continue
			}
			f := float64(n)
			if f == 0 {
				continue
			}
			safeGaugeSet(prometheus.GaugeOpts{
				Name:        "app_router",
				ConstLabels: safeLabels(map[string]string{"if": string(b), "octets": dir[i]}),
			}, f)
		}
	}
}

func checkTCP(addr string) float64 {
	conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
	if err != nil {
		log("Connecting error:", addr, err)
		return 10
	}
	if conn == nil {
		log("Connection is nil:", addr)
		return 11
	}
	err = conn.Close() // CAUTION! no timeout? goroutine leaking risk
	if err != nil {
		log("Close error:", addr, err)
		return 12
	}
	return 1 // ok
}

func updateTCP() {
	type p struct {
		label string
		value float64
	}
	ch := make(chan p)
	time.AfterFunc(300*time.Millisecond, func() {})
	for _, v := range [][2]string{
		{"10.72.0.2:179", "hop_2"},
		{"192.168.126.218:179", "hop_3"},
	} {
		v := v
		go func() {
			ch <- p{
				label: v[1],
				value: checkTCP(v[0]),
			}
		}()
	}
	for i := 0; i < 2; i++ {
		x := <-ch
		safeGaugeSet(prometheus.GaugeOpts{
			Name:        "app_network_reachable",
			ConstLabels: safeLabels(map[string]string{"point": x.label}),
		}, x.value)
	}
}

func bindAddrArg() string {
	if len(os.Args) >= 2 {
		return os.Args[1]
	}
	return ":9190"
}

func bindAddr() string {
	addr := bindAddrArg()
	log("Going to listen on", addr)
	return addr
}

func main() {
	gosnmp.Default.Target = "192.168.199.1" // ugly, but recommended by author in docs; and good enough for such small project
	err := gosnmp.Default.Connect()
	if err != nil {
		panic(err)
	}
	defer gosnmp.Default.Conn.Close()

	metricHandler := promhttp.Handler()
	http.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // TODO will be nice to pass r.Context everywhere
		log("Update metrics")
		updateLocalNetwork()
		updateSensors()
		updateLoadAvg()
		updateCPU()
		updateBlocks()
		updateRouterNetwork()
		updateTCP()
		metricHandler.ServeHTTP(w, r)
	}))

	err = http.ListenAndServe(bindAddr(), nil)
	if err != nil {
		log("listener error:", err)
	}
}
