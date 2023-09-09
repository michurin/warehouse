package slogtotext

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"text/template"
	"time"
)

const maxParsibleLen = 4096

type pplog struct {
	mx         *sync.Mutex
	next       io.Writer
	collecting bool
	buff       []byte
	logline    *template.Template
	errline    *template.Template
	knownKeys  map[string]any
}

func (pp *pplog) Write(p []byte) (int, error) {
	pp.mx.Lock()
	defer pp.mx.Unlock()
	n := len(p)
	for len(p) > 0 {
		s := bytes.IndexByte(p, '\n')
		if s == -1 {
			err := pp.acc(p)
			if err != nil {
				return 0, err // hmm... zero?
			}
			break
		} else {
			err := pp.fin(p[:s+1])
			if err != nil {
				return 0, err
			}
			p = p[s+1:]
		}
	}
	return n, nil
}

func (pp *pplog) fin(p []byte) error {
	err := pp.acc(p)
	if err != nil {
		return err
	}
	if pp.collecting {
		// https://github.com/golang/go/issues/24963
		data := any(nil)
		d := json.NewDecoder(bytes.NewReader(pp.buff))
		d.UseNumber()
		err := d.Decode(&data)
		if err != nil {
			err = pp.errline.Execute(pp.next, string(pp.buff[:len(pp.buff)-1])) // here in fin() we are sure that we have \n and the end of the buffer [fragile]
			pp.buff = pp.buff[:0]
			if err != nil {
				return err
			}
			if err != nil {
				return nil
			}
			return nil
		}
		pp.buff = pp.buff[:0]
		mdata, ok := data.(map[string]any)
		if ok {
			mdata["UNKNOWN"] = unknowPairs("", pp.knownKeys, data)
			data = mdata
		}
		err = pp.logline.Execute(pp.next, data)
		if err != nil {
			return err
		}
	}
	pp.collecting = true
	return nil
}

func (pp *pplog) acc(p []byte) error {
	if len(pp.buff)+len(p) > maxParsibleLen {
		pp.collecting = false
		err := pp.flush()
		if err != nil {
			return err
		}
	}
	if pp.collecting {
		pp.buff = append(pp.buff, p...)
		return nil
	}
	_, err := pp.next.Write(p)
	if err != nil {
		return err
	}
	return nil
}

func (pp *pplog) flush() error { // TODO it is used in acc only
	if len(pp.buff) == 0 {
		return nil
	}
	_, err := pp.next.Write(pp.buff)
	pp.buff = pp.buff[:0]
	if err != nil {
		return err
	}
	return nil
}

func PPLog(writer io.Writer, errlineTemplate, loglineTemplate string, knownKeys map[string]any, funcMap map[string]any) io.Writer {
	// TODO: validate knownKeys
	fm := template.FuncMap{"tmf": func(from, to string, tm any) string {
		ts, ok := tm.(string)
		if !ok {
			return fmt.Sprintf("invalid time type: %[1]T (%[1]v)", tm)
		}
		t, err := time.Parse(from, ts)
		if err != nil {
			return err.Error()
		}
		return t.Format(to)
	}}
	for k, v := range funcMap {
		fm[k] = v
	}
	if len(errlineTemplate) == 0 {
		errlineTemplate = `INVALID JSON: {{. | printf "%q"}}`
	}
	ll := template.Must(template.New("l").Option("missingkey=zero").Funcs(fm).Parse(loglineTemplate + "\n"))
	el := template.Must(template.New("e").Option("missingkey=zero").Funcs(fm).Parse(errlineTemplate + "\n"))
	return &pplog{
		mx:         new(sync.Mutex),
		next:       writer,
		collecting: true,
		buff:       make([]byte, 0, maxParsibleLen),
		logline:    ll,
		errline:    el,
		knownKeys:  knownKeys,
	}
}
