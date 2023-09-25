package aw

// Application-wide global variables
// Not sure it is brilliant idea, however, injections look like overkill for such small project

import "context"

var L = func(context.Context, any) {} // TODO atomic etc
