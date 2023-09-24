package icon_handlers

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func IconFromPng() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(1 >> 16)
		layers := []iconLayer(nil)
		for idx := 1; idx <= 10; idx++ {
			v := r.MultipartForm.File[fmt.Sprintf("file%d", idx)]
			if v == nil {
				continue
			}
			layer, err := parsePNG(v[0])
			if err != nil {
				panic(err) // TODO
			}
			layers = append(layers, layer)
		}
		iconHeaders := []byte(nil)
		iconBodies := []byte(nil)
		bodyShift := 6 + 16*len(layers)
		for _, layer := range layers {
			iconHeaders = append(iconHeaders, header2(layer, bodyShift)...)
			iconBodies = append(iconBodies, layer.data...)
			bodyShift += len(layer.data)
		}
		body := append(header1(len(layers)), append(iconHeaders, iconBodies...)...)
		w.Header().Add("content-type", `application/download; name="favicon.ico"`)
		w.Header().Add("content-disposition", `attachment; filename="favicon.ico"`)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func header1(n int) []byte {
	res := []byte{0, 0, 1, 0, 0, 0}
	binary.LittleEndian.PutUint16(res[4:], uint16(n))
	return res
}

func header2(layer iconLayer, shift int) []byte {
	res := make([]byte, 16)
	res[0] = uint8(layer.widht)
	res[1] = uint8(layer.height)
	res[4] = 1 // has to be checked
	binary.LittleEndian.PutUint16(res[6:], uint16(layer.bits))
	binary.LittleEndian.PutUint32(res[8:], uint32(len(layer.data)))
	binary.LittleEndian.PutUint32(res[12:], uint32(shift))
	return res
}

type iconLayer struct {
	data   []byte
	widht  int
	height int
	bits   int
}

/* TODO check w, h
   if ((images == 1 and (width != 16 or height != 16)) or
       (width > 255 or height > 255) or
       (width < 16 or height < 16)):
       raise ErrorSize('%dx%d' % (width, height))
*/

func parsePNG(body *multipart.FileHeader) (iconLayer, error) {
	layer := iconLayer{}
	fh, err := body.Open()
	if err != nil {
		return layer, fmt.Errorf("open: %w", err)
	}
	val, err := io.ReadAll(fh)
	if err != nil {
		return layer, fmt.Errorf("read: %w", err)
	}
	// TODO check len
	if !bytes.Equal(val[:8], []byte("\x89PNG\r\n\x1a\n")) {
		return layer, fmt.Errorf("invalid signature: %q", val[:8])
	}
	shift := 8
	for shift <= len(val)-12 {
		chankLen := int(binary.BigEndian.Uint32(val[shift : shift+4]))
		chankType := string(val[shift+4 : shift+8])
		if chankType == "IEND" {
			return layer, fmt.Errorf("found IEND without IHDR: shift=%d", shift)
		}
		if chankType == "IHDR" {
			if shift+18 >= len(val) {
				return layer, fmt.Errorf("hmm... broken IHDR")
			}
			layer.widht = int(binary.BigEndian.Uint32(val[shift+8 : shift+12]))
			layer.height = int(binary.BigEndian.Uint32(val[shift+12 : shift+16]))
			bits := int(val[shift+16])
			ctype := int(val[shift+17])
			switch ctype { // https://www.w3.org/TR/png/#11IHDR
			case 2:
				bits *= 3 // rgb
			case 4:
				bits *= 2 // grayscale+alpha
			case 6:
				bits *= 4 // rgba
			}
			layer.bits = bits
			layer.data = val
			return layer, nil
		}
		shift += chankLen + 12
	}
	return layer, fmt.Errorf("IHDR not found in whole file")
}

func IconFromData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		image, err := makeIcon(r.Body)
		if err != nil {
			// TODO log
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("content-type", `application/download; name="favicon.ico"`)
		w.Header().Add("content-disposition", `attachment; filename="favicon.ico"`)
		w.WriteHeader(http.StatusOK)
		w.Write(image)
	}
}

func makeIcon(body io.Reader) ([]byte, error) {
	reqData, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	values, err := url.ParseQuery(string(reqData))
	if err != nil {
		return nil, fmt.Errorf("parse qs: %w", err)
	}
	data := values.Get("data")
	if len(data) != 352 { // 16*16+16*6
		return nil, fmt.Errorf("invalid length: %d", len(data))
	}
	if strings.Trim(data, "0123456789abcdefgABCDEFG") != "" {
		return nil, fmt.Errorf("invalid chars: %q", data) // length is checked above
	}
	image := make([]byte, 318)
	copy(image, []byte{
		0, 0, 0x1, 0, 0x1, 0, 0x10, 0x10, 0x10, 0, 0x1, 0, 0x4, 0, 0x28, 0x1, 0, 0,
		0x16, 0, 0, 0, 0x28, 0, 0, 0, 0x10, 0, 0, 0, 0x20, 0, 0, 0, 0x1, 0, 0x4, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	})
	k := 62                   // length of header above
	j := 256                  // shift of colors table in data
	for i := 0; i < 16; i++ { // 16 colors
		for m := 4; m >= 0; m -= 2 { // revert r, g, b
			v, err := strconv.ParseInt(data[j+m:j+m+2], 16, 32)
			if err != nil {
				return nil, fmt.Errorf("color: i=%d, m=%d, v=%d: %w", i, m, v, err)
			}
			image[k] = uint8(v)
			k++
		}
		j += 6
		k++ // image[k]=0
	}
	for j := 0; j < 16; j++ { // lines
		for i := 0; i < 16; i++ {
			k = (15-j)*16 + i
			v, err := strconv.ParseInt(data[k:k+1], 17, 32)
			if err != nil {
				return nil, fmt.Errorf("pixels: j=%d, i=%d, k=%d, v=%d: %w", j, i, k, v, err)
			}
			if v < 16 {
				w := uint8(v)
				if i%2 == 0 {
					w <<= 4
				}
				image[126+j*8+i/2] |= w
			} else {
				image[254+j*4+i/8] |= 128 >> (i % 8)
			}
		}
	}
	return image, nil
}
