package readcloserwatcher_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/michurin/warehouse/go/readcloserwatcher"
)

func ExampleGolderFlow() {
	someReadClower := ioutil.NopCloser(bytes.NewBufferString("data"))
	someReadClower, watcher := readcloserwatcher.Watcher(someReadClower, time.Minute)
	output, err := ioutil.ReadAll(someReadClower)
	if err != nil {
		panic(err)
	}
	someReadClower.Close()
	fmt.Printf("We have read %q\n", output)
	result := <-watcher
	fmt.Printf("And caught by watcher %q\n", result.Octets)
	// OUTPUT:
	// We have read "data"
	// And caught by watcher "data"
}
