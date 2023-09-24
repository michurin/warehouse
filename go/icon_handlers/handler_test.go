package icon_handlers_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/michurin/warehouse/go/icon_handlers"
)

func TestIconFromPng(t *testing.T) {
	h := icon_handlers.IconFromPng()
	s := httptest.NewServer(h)
	defer s.Close()
	outBuff := new(bytes.Buffer)
	errBuff := new(bytes.Buffer)
	cmd := exec.Command( // I use curl. I'm too lazy to prepare multipart request by myself
		"curl",
		"-q", // must be first
		"-s",
		s.URL+"/bin/png/favicon.ico",
		"-F", "file1=@testdata/16x16.png",
		"-F", "file2=@testdata/32x32.png",
		"-D", "/dev/stderr") // wired way to check headers
	cmd.Stdout = outBuff
	cmd.Stderr = errBuff
	err := cmd.Run()
	require.NoError(t, err)
	headers := errBuff.String()
	assert.Contains(t, headers, "HTTP/1.1 200 OK\r\n")
	assert.Contains(t, headers, "Content-Disposition: attachment; filename=\"favicon.ico\"\r\n")
	assert.Contains(t, headers, "Content-Type: application/download; name=\"favicon.ico\"\r\n")
	assert.Contains(t, headers, "Content-Length: 1690\r\n")
	icon, err := os.ReadFile("testdata/favicon_png_test.ico")
	require.NoError(t, err)
	assert.Equal(t, icon, outBuff.Bytes())
}

func TestIconFromData(t *testing.T) {
	h := icon_handlers.IconFromData()
	r, err := http.NewRequest(http.MethodPost, "/bin/data/favicon.ico", bytes.NewBuffer([]byte("data=00000000000000000ggggggggggggggg0g999999999999990g9ggggggggggggg0g9gaaaaaaaaaaaa0g9gaggggggggggg0g9gagcccccccccc0g9gagcggggggggg0g9gagcggggggggg0g9gagcggggggggg0g9gagcggggggggg0g9gagcggggggggg0g9gagcggggggggg0g9gagcggggggggg0g9gagcggggggggg0g9gagcggggggggg000000800000008000808000000080800080008080c0c0c0808080ff000000ff00ffff000000ffff00ff00ffffffffff")))
	require.NoError(t, err)
	w := httptest.NewRecorder()
	h(w, r)
	resp := w.Result()
	assert.Equal(t, http.Header{
		"Content-Disposition": {`attachment; filename="favicon.ico"`},
		"Content-Type":        {`application/download; name="favicon.ico"`},
	}, resp.Header)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	icon, err := os.ReadFile("testdata/favicon_test.ico")
	require.NoError(t, err)
	assert.Equal(t, icon, body)
}

/*
	// dump image
	fmt.Println(hex.EncodeToString(image[:62]))
	for i := 0; i < 16; i++ { // 64
		fmt.Println(hex.EncodeToString(image[62+i*4 : 62+i*4+4]))
	}
	for i := 0; i < 16; i++ { // 128
		fmt.Println(hex.EncodeToString(image[126+i*8 : 126+i*8+8]))
	}
	for i := 0; i < 16; i++ { // 64
		fmt.Println(hex.EncodeToString(image[254+i*4 : 254+i*4+4]))
	}
	fmt.Println(hex.EncodeToString(image[62+256:]))
*/
