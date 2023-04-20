// Exercise: rot13Reader
package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (r rot13Reader) Read(buf []byte) (n int, err error) {
	if n, err = r.r.Read(buf); err != nil {
		return 0, err
	}
	for i, r := range buf {
		switch {
		case r >= 'A' && r <= 'Z':
			buf[i] = 'A' + (r-'A'+13)%26
		case r >= 'a' && r <= 'z':
			buf[i] = 'a' + (r-'a'+13)%26
		}
	}
	return n, nil
}

func main() {
	r := rot13Reader{strings.NewReader("Lbh penpxrq gur pbqr!")}
	io.Copy(os.Stdout, &r)
}
