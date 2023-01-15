package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/michurin/warehouse/go/smartjson/pkg/smartjson"
)

func consoleSize() (int, int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	s := string(out)
	s = strings.TrimSpace(s)
	sArr := strings.Split(s, " ")

	heigth, err := strconv.Atoi(sArr[0])
	if err != nil {
		return 0, 0, err
	}

	width, err := strconv.Atoi(sArr[1])
	if err != nil {
		return 0, 0, err
	}
	return heigth, width, nil
}

// echo '{"1": [1, {"x":"xx"}, 4444]}' | go run main.go
// echo '{"1": [[1, 2], {"x":{"znn":111111}}, true, null]}' | go run main.go
func main() {
	buff, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	target := any(nil)
	err = json.Unmarshal(buff, &target)
	if err != nil {
		panic(err)
	}
	fmt.Println(target)
	fmt.Println("====")
	r := smartjson.Marshal(target, &smartjson.Opts{
		Width:  30,
		Indent: 2,
	})
	fmt.Println(r)
	w, h, err := consoleSize()
	fmt.Println(w, h, err)
}
