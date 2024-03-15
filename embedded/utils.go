package embedded

import (
	"embed"
	"io"
	"log"
)

const (
	FontPrimeMono     string = "prime-mono.ttf"
	FontOpenGostTypeA string = "open-gost-type-a.ttf"
)

//go:embed prime-mono.ttf
//go:embed open-gost-type-a.ttf
var f embed.FS

func GetFontReader(font string) io.Reader {
	r, e := f.Open(font)
	if e != nil {
		log.Panic(e)
	}

	return r
}
