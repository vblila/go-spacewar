package graphics

import (
	"github.com/go-gl/gltext"
	"io"
	"log"
	"strconv"
)

type FontsBuffer struct {
	buffer map[string]*gltext.Font
}

func (fb *FontsBuffer) Init() {
	fb.buffer = make(map[string]*gltext.Font, 10)
}

func (fb *FontsBuffer) getKey(name string, scale int) string {
	return name + strconv.Itoa(scale)
}

func (fb *FontsBuffer) Load(name string, scale int, reader io.Reader) {
	var err error
	fb.buffer[fb.getKey(name, scale)], err = gltext.LoadTruetype(reader, int32(scale), 0, 127, gltext.LeftToRight)
	if err != nil {
		log.Panic(err)
	}
}

func (fb *FontsBuffer) Get(name string, scale int) *gltext.Font {
	font, ok := fb.buffer[fb.getKey(name, scale)]
	if !ok {
		return nil
	}
	return font
}
