package minlog

import "context"

type Patch struct {
	kv map[string]any
}

func TakePatch(ctx context.Context) Patch {
	return Patch{kv: ctxKv(ctx)}
}

func ApplyPatch(ctx context.Context, patch Patch) context.Context {
	kv := ctxKv(ctx)
	for k, v := range patch.kv {
		kv[k] = v
	}
	return context.WithValue(ctx, ctxKey, kv)
}
