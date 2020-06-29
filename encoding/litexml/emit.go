package litexml

import (
	"bufio"
	"encoding/xml"
	"io"
)

// TODO: Support for encodings other than UTF-8.
// Consider writing directly to io.Writer and managing encoding.

type emitter struct {
	w bufio.Writer
}

func newemitter(w io.Writer) *emitter {
	return &emitter{
		w: *bufio.NewWriter(w),
	}
}

func (e *emitter) emitattr(key, value string) {
	e.w.WriteString(" ")
	e.w.WriteString(key)
	e.w.WriteString(`="`)
	xml.EscapeText(&e.w, []byte(value))
	e.w.WriteString(`"`)
}

func (e *emitter) emitdt(dt DocumentInfo) {
	e.w.WriteString("<?xml")
	e.emitattr("version", dt.Version)
	if dt.Encoding == "" {
		dt.Encoding = "utf-8"
	}
	e.emitattr("encoding", dt.Encoding)
	if dt.Standalone != "" {
		e.emitattr("standalone", dt.Standalone)
	}
	e.w.WriteString(" ?>\n")
}

func (e *emitter) emittagopenpart(tag string) {
	e.w.WriteString("<")
	e.w.WriteString(tag)
}

func (e *emitter) emittagclosepart(tag string) {
	e.w.WriteString("</")
	e.w.WriteString(tag)
}

func (e *emitter) emittagendpart() {
	e.w.WriteString(">\n")
}

func (e *emitter) emittagcloseendpart() {
	e.w.WriteString(" />\n")
}

func (e *emitter) emitcontent(content string) {
	// TODO: EscapeText is overzealous, should write our own.
	xml.EscapeText(&e.w, []byte(content))
	// TODO: probably should try to detect if we can safely do this without messing up whitespace collapse
	e.w.WriteString("\n")
}

func (e *emitter) indent(level int) {
	for i := 0; i < level; i++ {
		e.w.WriteString("        ")
	}
}

func (e *emitter) flush() error {
	return e.w.Flush()
}
