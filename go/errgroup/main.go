package main

import (
	"context"
	"errors"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

func mkfunc(ctx context.Context, name string, duration time.Duration, result error) func() error {
	return func() error {
		reason := "n/a"
		log.Printf("[%s] start...", name)
		defer func() { log.Printf("[%s] fin: reason=%s error=%v", name, reason, result) }()
		select {
		case <-time.After(duration):
			reason = "timeout"
		case <-ctx.Done():
			reason = "ctx: " + ctx.Err().Error() // In real world it is nice to return ctx.Err() here
		}
		return result
	}
}

func main() {
	ctx := context.Background()
	gr, ctx := errgroup.WithContext(ctx)
	gr.Go(mkfunc(ctx, "a", 3*time.Second, nil))                         // Just regular run
	gr.Go(mkfunc(ctx, "b", 6*time.Second, errors.New("Err")))           // Finished with error
	gr.Go(mkfunc(ctx, "c", 9*time.Second, nil))                         // Will be canceled by ctx, due to error in "b"
	gr.Go(mkfunc(context.WithoutCancel(ctx), "d", 12*time.Second, nil)) // Won't exit on ctx; WithoutCancel works
	err := gr.Wait()                                                    // Wait for all, but manage ctx as expected: cancel it right after first error
	log.Printf("Waiting result: %v", err)
}
