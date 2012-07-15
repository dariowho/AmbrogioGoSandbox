package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (self rot13Reader) Read(p []byte) (n int, err error) {
	// I have no idea how to write into p directly...
	// and this syntax sucks
	n,err = self.r.Read(p)
	
	for i:=0; i<len(p); i++ {
		c := p[i]
		switch {
			case c >= 'A' && c<= 'Z':
				c=(((c-'A')+13)%('Z'-'A'+1))+'A'
			case c >= 'a' && c<= 'z':
				c=(((c-'a')+13)%('z'-'a'+1))+'a'
		}
		
		p[i]=c
	}
	
	return
}

func main() {
	s := strings.NewReader(
		"Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
