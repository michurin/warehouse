package minlog

// WithStdLogger sets output stream accepting std Go logger.
func WithStdLogger(printer interface{ Print(v ...any) }) Option {
	return func(l *Logger) {
		l.printer = func(x string) {
			printer.Print(x)
		}
	}
}

func WithFields(fields ...FieldFunc) Option {
	return func(l *Logger) {
		l.fields = fields
	}
}

func WithArgFormatter(formatter func(any) string) Option {
	return func(l *Logger) {
		l.argFormatter = formatter
	}
}

func WithPersistFields(kv ...string) Option {
	p := map[string]string{}
	for i := 0; i < len(kv)-1; i++ {
		p[kv[i]] = kv[i+1]
	}
	return func(l *Logger) {
		l.persistFields = p
	}
}
