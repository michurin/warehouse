package minlog

import (
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"
)

func FieldMessage() FieldFunc {
	return func(r Record) string {
		return r.Message
	}
}

func FieldCaller(pfx string) FieldFunc {
	return func(r Record) string {
		return fmt.Sprintf("%s:%d", strings.TrimPrefix(r.Caller.File, pfx), r.Caller.Line)
	}
}

func FieldErrorCaller(pfx string) FieldFunc {
	return func(r Record) string {
		if r.IsError {
			if r.ErrorCaller.File == "" {
				return ""
			}
			return fmt.Sprintf("%s:%d", strings.TrimPrefix(r.ErrorCaller.File, pfx), r.ErrorCaller.Line)
		}
		return ""
	}
}

func FieldLevel(info, errr string) FieldFunc {
	return func(r Record) string {
		if r.IsError {
			return errr
		}
		return info
	}
}

func FieldFallbackKV(exclude ...string) FieldFunc {
	exc := map[string]struct{}{}
	for _, v := range exclude {
		exc[v] = struct{}{}
	}
	return func(r Record) string {
		fs := []string(nil)
		for k := range r.Context {
			if _, ok := exc[k]; ok {
				continue
			}
			fs = append(fs, k)
		}
		if fs == nil {
			return ""
		}
		sort.Strings(fs)
		pts := make([]string, len(fs))
		for i, k := range fs {
			pts[i] = fmt.Sprintf("%s=%v", k, r.Context[k])
		}
		return strings.Join(pts, " ")
	}
}

func FieldNamed(fieldName string) FieldFunc {
	return func(r Record) string {
		if x, ok := r.Context[fieldName]; ok {
			return fmt.Sprintf("%v", x)
		}
		return ""
	}
}

func addCaller(kv map[string]any, k string, rc RecordCaller) {
	if rc.File == "" && rc.Line == 0 {
		return
	}
	kv[k] = path.Base(rc.File) + ":" + strconv.Itoa(rc.Line)
}

func FieldJSON() FieldFunc {
	return func(r Record) string {
		level := "info"
		if r.IsError {
			level = "error"
		}
		kv := map[string]any{
			"level":   level,
			"message": r.Message,
			"context": r.Context,
		}
		addCaller(kv, "caller", r.Caller)
		addCaller(kv, "error_caller", r.ErrorCaller)
		b, err := json.Marshal(kv)
		if err != nil {
			b, _ = json.Marshal(map[string]string{
				"marshaller_error": "FieldJSON error: " + err.Error(),
			})
		}
		return string(b)
	}
}
