package minlog

import "runtime"

func caller(level int) RecordCaller {
	_, file, no, ok := runtime.Caller(level)
	if !ok {
		return RecordCaller{
			File: "nofile",
			Line: 0,
		}
	}
	return RecordCaller{
		File: file,
		Line: no,
	}
}
