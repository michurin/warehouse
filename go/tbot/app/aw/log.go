package aw

// Application-wide global variables
// Not sure it is brilliant idea, however, injections look like overkill for such small project

import "context"

var Log = func(context.Context, ...any) {}
