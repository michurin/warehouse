package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ReadSensor(hwpath string) (map[[2]string]float64, error) {
	lst, err := os.ReadDir(hwpath)
	if err != nil {
		return nil, fmt.Errorf("cannot read hw path: %w", err)
	}
	name := ""
	input := map[string]string{}
	label := map[string]string{}
	all := map[string]struct{}{}
	for _, e := range lst {
		n := e.Name()
		p := path.Join(hwpath, n)
		switch {
		case n == "name":
			v, err := os.ReadFile(p)
			if err != nil {
				return nil, err
			}
			name = strings.TrimSpace(string(v))
		case strings.HasPrefix(n, "temp") && strings.HasSuffix(n, "_input"):
			k := n[4 : len(n)-6]
			all[k] = struct{}{}
			v, err := os.ReadFile(p)
			if err != nil {
				return nil, err
			}
			input[k] = strings.TrimSpace(string(v))
			fmt.Println("term k", k, input[k]) // TODO
		case strings.HasPrefix(n, "temp") && strings.HasSuffix(n, "_label"):
			k := n[4 : len(n)-6]
			all[k] = struct{}{}
			v, err := os.ReadFile(p)
			if err != nil {
				return nil, err
			}
			label[k] = strings.TrimSpace(string(v))
		}
	}
	res := map[[2]string]float64{}
	for k := range all {
		v, err := strconv.Atoi(input[k])
		if err != nil {
			return nil, err
		}
		fmt.Println("term n k", name, label[k], v) // TODO
		res[[2]string{name, label[k]}] = float64(v) / 1000
	}
	return res, nil
}

func ReadSensors() (map[[2]string]float64, error) {
	spath := "/sys/class/hwmon"
	lst, err := os.ReadDir(spath)
	if err != nil {
		return nil, fmt.Errorf("cannot read sensors path: %w", err)
	}
	res := map[[2]string]float64{}
	for _, e := range lst {
		r, err := ReadSensor(path.Join(spath, e.Name()))
		if err != nil {
			return nil, err
		}
		for k, v := range r {
			res[k] = v // todo check dups
		}
	}
	return res, nil
}

//

func ReadInterface(iface, netpath string) (map[[2]string]float64, error) {
	lst, err := os.ReadDir(netpath)
	if err != nil {
		return nil, fmt.Errorf("cannot read path: %w", err)
	}
	res := map[[2]string]float64{}
	for _, e := range lst {
		n := e.Name()
		p := path.Join(netpath, n)
		b, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		v, err := strconv.Atoi(strings.TrimSpace(string(b)))
		if err != nil {
			return nil, err
		}
		res[[2]string{iface, n}] = float64(v)
	}
	return res, nil
}

func ReadInterfaces() (map[[2]string]float64, error) {
	spath := "/sys/class/net"
	lst, err := os.ReadDir(spath)
	if err != nil {
		return nil, fmt.Errorf("cannot read interfaces path: %w", err)
	}
	res := map[[2]string]float64{}
	for _, e := range lst {
		n := e.Name()
		r, err := ReadInterface(n, path.Join(spath, n, "statistics"))
		if err != nil {
			return nil, err
		}
		for k, v := range r {
			res[k] = v // todo check dups
		}
	}
	return res, nil
}

//

func safe(x string) string {
	if x == "" {
		return "undef"
	}
	return strings.ReplaceAll(x, " ", "_")
}

type Metrics struct {
	m map[[2]string]float64
	t int64
	u []func() (map[[2]string]float64, error)
}

func (m *Metrics) lazyUpdate() {
	t := time.Now().Unix()
	if m.t < t {
		m.t = t + 10
		m.m = map[[2]string]float64{}
		for _, f := range m.u {
			fmt.Println("call F")
			s, err := f()
			fmt.Printf("s: %#v %#v\n", s, err)
			if err != nil {
				fmt.Println(err)
				return
			}
			for k, v := range s {
				m.m[k] = v
			}
		}
	}
}

func (m *Metrics) Get(k [2]string) float64 {
	m.lazyUpdate()
	return m.m[k]
}

func (m *Metrics) All() [][2]string {
	m.lazyUpdate()
	x := make([][2]string, 0, len(m.m))
	for k := range m.m {
		x = append(x, k)
	}
	return x
}

func main() {
	http.Handle("/metrics", promhttp.Handler())

	{
		mx := new(Metrics)
		mx.u = []func() (map[[2]string]float64, error){ReadSensors}
		for _, k := range mx.All() {
			k := k
			fmt.Println("init gauge", k)
			promauto.NewGaugeFunc(prometheus.GaugeOpts{
				Namespace: "app",
				Subsystem: safe(k[0]),
				Name:      safe(k[1]),
			}, func() float64 {
				return mx.Get(k)
			})
		}
	}

	{
		mx := new(Metrics)
		mx.u = []func() (map[[2]string]float64, error){ReadInterfaces}
		for _, k := range mx.All() {
			k := k
			fmt.Println("init counter", k)
			promauto.NewCounterFunc(prometheus.CounterOpts{
				Namespace: "app",
				Subsystem: safe(k[0]),
				Name:      safe(k[1]),
			}, func() float64 {
				return mx.Get(k)
			})
		}
	}

	err := http.ListenAndServe(":9190", nil)
	if err != nil {
		fmt.Printf("listener error: %s\n", err)
	}
}
