package litexml

import (
	"bufio"
	"io"
	"reflect"
	"strconv"
	"strings"
)

type decstate struct {
	val  reflect.Value
	tag  string
	attr string

	tags    map[string]int
	attrs   map[string]int
	content int
	doctype int
}

// Decoder is a decoder for the lite XML format.
type Decoder struct {
	p     parser
	intag bool
	stack []decstate
	cur   decstate
}

// NewDecoder creates a new decoder.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		p:     *newparser(bufio.NewReader(r)),
		intag: false,
	}
}

func (d *Decoder) push() {
	d.stack = append(d.stack, d.cur)
}

func (d *Decoder) pop() {
	d.cur = d.stack[len(d.stack)-1]
	d.stack = d.stack[:len(d.stack)-1]
}

func (d *Decoder) doctype(dt DocumentInfo) {
	if d.cur.doctype != -1 {
		d.cur.val.Field(d.cur.doctype).Set(reflect.ValueOf(dt))
	}
}

func (d *Decoder) opentag(tag string) {
	v := d.cur.val
	i, ok := d.cur.tags[tag]
	if !ok {
		return
	}

	fv := v.Field(i)
	ft := v.Type().Field(i)

	switch ft.Type.Kind() {
	case reflect.Struct:
		d.setval(fv, true, tag, "")
	case reflect.Slice:
		fv.Set(reflect.Append(fv, reflect.New(ft.Type.Elem()).Elem()))
		d.setval(fv.Index(fv.Len()-1), true, tag, "")
	default:
		d.setval(fv, true, tag, ft.Tag.Get("attr"))
	}
}

func (d *Decoder) closetag(tag string) {
	d.pop()
}

func (d *Decoder) setscalarattr(v reflect.Value, value string) {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, _ := strconv.ParseInt(value, 10, v.Type().Bits())
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, _ := strconv.ParseUint(value, 10, v.Type().Bits())
		v.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, _ := strconv.ParseFloat(value, v.Type().Bits())
		v.SetFloat(f)
	default:
		panic("don't know how to handle field type")
	}
}

func (d *Decoder) attr(key, value string) {
	v := d.cur.val
	switch v.Kind() {
	case reflect.Struct:
		if i, ok := d.cur.attrs[key]; ok {
			d.setscalarattr(v.Field(i), value)
		}
	default:
		if d.cur.attr == key {
			d.setscalarattr(v, value)
		}
	}
}

func (d *Decoder) content(data string) {
	if d.cur.content != -1 {
		// TODO: trimming space is not safe (will change whitespace behavior)
		d.setscalarattr(d.cur.val.Field(d.cur.content), strings.TrimSpace(data))
	}
}

func (d *Decoder) setval(val reflect.Value, push bool, tag string, attr string) {
	if push {
		d.push()
	}

	d.cur = decstate{
		val:     val,
		tag:     tag,
		attr:    attr,
		tags:    map[string]int{},
		attrs:   map[string]int{},
		content: -1,
		doctype: -1,
	}

	switch val.Kind() {
	case reflect.Struct:
		typ := val.Type()
		for i, n := 0, typ.NumField(); i < n; i++ {
			rft := typ.Field(i)
			if rft.Name == "" || rft.Name[0] == '_' || (rft.Name[0] >= 'a' && rft.Name[0] <= 'z') {
				continue
			}
			tag := rft.Tag.Get("tag")
			if tag != "" {
				d.cur.tags[tag] = i
			}
			attr := rft.Tag.Get("attr")
			if attr != "" {
				d.cur.attrs[attr] = i
			}
			if rft.Tag.Get("content") == "inner" {
				d.cur.content = i
			}
			if rft.Type == reflect.TypeOf(DocumentInfo{}) {
				d.cur.doctype = i
			}
		}
	case reflect.Array:
	case reflect.Slice:
	}
}

func (d *Decoder) decode(rv reflect.Value) error {
	if rv.Kind() != reflect.Struct {
		panic("decode type must be struct")
	}

	d.setval(rv, false, "", "")

	for {
		err := d.p.parsenext(d)
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}
}

// Decode decodes an XML document to the value.
func (d *Decoder) Decode(value interface{}) error {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr {
		panic("decode type must be pointer to struct")
	}
	return d.decode(rv.Elem())
}
