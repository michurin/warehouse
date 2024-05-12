package aw

// Application-wide global variables.
// Not sure it is brilliant idea, however, injections look like overkill for such small project.
// I considered slog.SetDefault, it seems too wide.
// So, I decided just to keep all such dirty stuff in one place (here).

import "context"

var L = func(context.Context, any) {} // atomic.Value etc would be more safe, however, seems overkill
