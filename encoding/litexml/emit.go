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

func (e *emitter) emitattr(key, value string) error {
	if _, err := e.w.WriteString(" "); err != nil {
		return err
	}
	if _, err := e.w.WriteString(key); err != nil {
		return err
	}
	if _, err := e.w.WriteString(`="`); err != nil {
		return err
	}
	if err := xml.EscapeText(&e.w, []byte(value)); err != nil {
		return err
	}
	if _, err := e.w.WriteString(`"`); err != nil {
		return err
	}
	return nil
}

func (e *emitter) emitdt(dt DocumentInfo) error {
	if _, err := e.w.WriteString("<?xml"); err != nil {
		return err
	}
	if err := e.emitattr("version", dt.Version); err != nil {
		return err
	}
	if dt.Encoding == "" {
		dt.Encoding = "utf-8"
	}
	if err := e.emitattr("encoding", dt.Encoding); err != nil {
		return err
	}
	if dt.Standalone != "" {
		if err := e.emitattr("standalone", dt.Standalone); err != nil {
			return err
		}
	}
	if _, err := e.w.WriteString(" ?>\n"); err != nil {
		return err
	}
	return nil
}

func (e *emitter) emittagopenpart(tag string) error {
	if _, err := e.w.WriteString("<"); err != nil {
		return err
	}
	if _, err := e.w.WriteString(tag); err != nil {
		return err
	}
	return nil
}

func (e *emitter) emittagclosepart(tag string) error {
	if _, err := e.w.WriteString("</"); err != nil {
		return err
	}
	if _, err := e.w.WriteString(tag); err != nil {
		return err
	}
	return nil
}

func (e *emitter) emittagendpart() error {
	_, err := e.w.WriteString(">\n")
	return err
}

func (e *emitter) emittagcloseendpart() error {
	_, err := e.w.WriteString(" />\n")
	return err
}

func (e *emitter) emitcontent(content string) error {
	// TODO: EscapeText is overzealous, should write our own.
	if err := xml.EscapeText(&e.w, []byte(content)); err != nil {
		return err
	}
	// TODO: probably should try to detect if we can safely do this without messing up whitespace collapse
	if _, err := e.w.WriteString("\n"); err != nil {
		return err
	}
	return nil
}

func (e *emitter) indent(level int) error {
	for i := 0; i < level; i++ {
		if _, err := e.w.WriteString("        "); err != nil {
			return err
		}
	}
	return nil
}
