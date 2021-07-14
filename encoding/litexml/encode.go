package litexml

import (
	"io"
	"reflect"
	"strconv"
)

// Encoder is an encoder for the lite XML format.
type Encoder struct {
	e      emitter
	intag  bool
	indent int
}

// NewEncoder creates a new encoder.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		e:      *newemitter(w),
		intag:  false,
		indent: 0,
	}
}

func (e *Encoder) tagattr(attr, value string) error {
	if err := e.e.emitattr(attr, value); err != nil {
		return err
	}
	return nil
}

func (e *Encoder) unarytag(tag, attr, value string) error {
	if err := e.e.indent(e.indent); err != nil {
		return err
	}
	if err := e.e.emittagopenpart(tag); err != nil {
		return err
	}
	if err := e.e.emitattr(attr, value); err != nil {
		return err
	}
	if err := e.e.emittagcloseendpart(); err != nil {
		return err
	}
	return nil
}

func (e *Encoder) scalarattrval(i interface{}) (string, bool) {
	switch t := i.(type) {
	case int:
		return strconv.FormatInt(int64(t), 10), true
	case int8:
		return strconv.FormatInt(int64(t), 10), true
	case int16:
		return strconv.FormatInt(int64(t), 10), true
	case int32:
		return strconv.FormatInt(int64(t), 10), true
	case int64:
		return strconv.FormatInt(int64(t), 10), true
	case uint:
		return strconv.FormatUint(uint64(t), 10), true
	case uint8:
		return strconv.FormatUint(uint64(t), 10), true
	case uint16:
		return strconv.FormatUint(uint64(t), 10), true
	case uint32:
		return strconv.FormatUint(uint64(t), 10), true
	case uint64:
		return strconv.FormatUint(uint64(t), 10), true
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 32), true
	case float64:
		return strconv.FormatFloat(float64(t), 'f', -1, 64), true
	case string:
		return t, true
	}
	return "", false
}

func (e *Encoder) encode(tag string, value interface{}) error {
	rt := reflect.TypeOf(value)
	if rt.Kind() != reflect.Struct {
		panic("encode type must be struct")
	}

	if tag != "" {
		if e.intag {
			if err := e.e.emittagendpart(); err != nil {
				return err
			}
		}
		if err := e.e.indent(e.indent); err != nil {
			return err
		}
		if err := e.e.emittagopenpart(tag); err != nil {
			return err
		}
		e.intag = true
		e.indent++
	}

	rv := reflect.ValueOf(value)
	for i, n := 0, rv.NumField(); i < n; i++ {
		rfv, rft := rv.Field(i), rt.Field(i)

		// Exclude non-public fields.
		if rft.Name == "" || rft.Name[0] == '_' || (rft.Name[0] >= 'a' && rft.Name[0] <= 'z') {
			continue
		}

		tag := rft.Tag.Get("tag")
		attr := rft.Tag.Get("attr")
		content := rft.Tag.Get("content")

		ifv := rfv.Interface()
		switch content {
		case "":
			if tag == "" {
				// DocType
				if fv, ok := ifv.(DocumentInfo); ok {
					if err := e.e.emitdt(fv); err != nil {
						return err
					}
					continue
				}
				// Attr in tag
				if !e.intag {
					panic("no tag to write attr to")
				}
				if attr == "" {
					panic("empty attr for standalone tag")
				}
				if scalar, ok := e.scalarattrval(ifv); ok {
					if err := e.tagattr(attr, scalar); err != nil {
						return err
					}
					continue
				}
			} else {
				if e.intag {
					e.intag = false
					if err := e.e.emittagendpart(); err != nil {
						return err
					}
				}
				if scalar, ok := e.scalarattrval(ifv); ok {
					if err := e.unarytag(tag, attr, scalar); err != nil {
						return err
					}
					continue
				}
			}
			for rft.Type.Kind() == reflect.Ptr {
				rft.Type = rft.Type.Elem()
				rfv = rfv.Elem()
			}
			if rft.Type.Kind() == reflect.Slice || rft.Type.Kind() == reflect.Array {
				for i, n := 0, rfv.Len(); i < n; i++ {
					if err := e.encode(tag, rfv.Index(i).Interface()); err != nil {
						return err
					}
				}
				continue
			} else if rft.Type.Kind() == reflect.Struct {
				if err := e.encode(tag, ifv); err != nil {
					return err
				}
				continue
			}
			panic("don't know what to do with field: " + rft.Name)

		case "inner":
			if e.intag {
				e.intag = false
				if err := e.e.emittagendpart(); err != nil {
					return err
				}
			}
			switch fv := ifv.(type) {
			case string:
				if err := e.e.indent(e.indent); err != nil {
					return err
				}
				if err := e.e.emitcontent(fv); err != nil {
					return err
				}
				continue
			default:
				panic("invalid content type")
			}
		}
	}

	if tag != "" {
		e.indent--
		if e.intag {
			e.intag = false
			if err := e.e.emittagcloseendpart(); err != nil {
				return err
			}
		} else {
			if err := e.e.indent(e.indent); err != nil {
				return err
			}
			if err := e.e.emittagclosepart(tag); err != nil {
				return err
			}
			if err := e.e.emittagendpart(); err != nil {
				return err
			}
		}
	}

	return nil
}

// Encode encodes a value as an XML document.
func (e *Encoder) Encode(value interface{}) error {
	err1 := e.encode("", value)
	err2 := e.e.w.Flush()
	if err1 != nil {
		return err1
	}
	return err2
}
