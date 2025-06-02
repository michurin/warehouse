package argerr_test

import (
	"errors"
	"fmt"

	"github.com/michurin/minlog/argerr"
)

func Example_simplestAlltogether() {
	err := errors.New("message")
	err = argerr.Wrap(err, 1, "a", "b")
	fmt.Println(argerr.Args(err))
	// output:
	// [1 a b]
}

func ExampleWrap_unwrap() {
	origErr := errors.New("message")
	err := origErr
	err = argerr.Wrap(err, 1)
	err = argerr.Wrap(err, 2)            // do nothing, already wrapped
	err = fmt.Errorf("Wrapped: %w", err) // next wrappers won't affect Args
	fmt.Println(argerr.Args(err))
	fmt.Println(errors.Is(err, origErr)) // true, unwrap and all related stuff works
	fmt.Println(err.Error())             // message is unchanged
	// output:
	// [1]
	// true
	// Wrapped: message
}

func ExampleWrap_nilErrorNilArgs() {
	err := error(nil)
	err = argerr.Wrap(err, 1) // result will be <nil>
	fmt.Println(err)
	err = errors.New("message")
	fmt.Println(err == argerr.Wrap(err)) // without args Wrap returns the same err
	// output:
	// <nil>
	// true
}

func ExampleArgs_nilError() {
	err := error(nil)
	fmt.Printf("%#v\n", argerr.Args(err)) // corner case: nil error, nil slice of args
	// output:
	// []interface {}(nil)
}
