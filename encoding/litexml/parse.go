package litexml

import (
	"fmt"
	"io"
	"strconv"
	"unicode"
)

// TODO: This should probably work at the byte level.
// TODO: Support for encodings other than UTF-8.

type syntaxerr struct {
	ln  int
	ch  int
	msg string
}

func (err syntaxerr) Error() string {
	return fmt.Sprintf("%d:%d: %s", err.ln+1, err.ch+1, err.msg)
}

type runematcher interface {
	match(r rune) bool
}

type matchone rune

func (m matchone) match(r rune) bool {
	return r == rune(m)
}

type matchmultiple []rune

func (m matchmultiple) match(r rune) bool {
	for _, n := range m {
		if r == rune(n) {
			return true
		}
	}
	return false
}

type matchfn func(r rune) bool

func (m matchfn) match(r rune) bool {
	return m(r)
}

var matchspace = matchmultiple{' ', '\t', '\n', '\v', '\f', '\r'}

var matchdigit = matchfn(func(r rune) bool {
	return r >= '0' && r <= '9'
})

var matchhex = matchfn(func(r rune) bool {
	return matchdigit(r) || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
})

var matchalpha = matchfn(func(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
})

var matchnamestart = matchfn(func(r rune) bool {
	return unicode.IsLetter(r) || r == '_' || r == ':'
})

var matchname = matchfn(func(r rune) bool {
	return matchnamestart(r) || unicode.IsDigit(r) || r == '.' || r == '-'
})

type parser struct {
	s io.RuneScanner

	// File position information.
	// Note: these are zero based internally.
	ln, ch int

	// Saved previous file position for unreading.
	ln0, ch0 int

	// Records unaccepted runes for debugging. Cleared whenever something is accepted.
	rejects []rune
}

func newparser(s io.RuneScanner) *parser {
	return &parser{
		s: s,

		ln: 0,
		ch: 0,
	}
}

func (p *parser) readrune() rune {
	r, _, err := p.s.ReadRune()

	p.ln0 = p.ln
	p.ch0 = p.ch

	if r == '\n' {
		p.ln++
		p.ch = 0
	} else {
		p.ch++
	}

	if err != nil {
		// TODO: wrap?
		panic(err)
	}

	return r
}

// this is only allowed to be called once in a row between readrune calls
func (p *parser) unreadrune() {
	err := p.s.UnreadRune()
	if err != nil {
		// TODO: wrap?
		panic(err)
	}

	p.ln = p.ln0
	p.ch = p.ch0
}

func (p *parser) badsyntax(format string, a ...interface{}) error {
	return syntaxerr{
		ln:  p.ln,
		ch:  p.ch,
		msg: fmt.Sprintf(format, a...),
	}
}

func (p *parser) expect(m runematcher, describe string) rune {
	r := p.readrune()
	if m.match(r) {
		p.rejects = p.rejects[0:0]
		return r
	}
	panic(p.badsyntax("expected %s, got %c", describe, r))
}

func (p *parser) expectstr(s string) {
	err := p.badsyntax("expected %q", s)
	p.rejects = p.rejects[0:0]
	for _, c := range s {
		r := p.readrune()
		if r != c {
			panic(err)
		}
	}
}

func (p *parser) accept(m runematcher, out *rune) bool {
	r := p.readrune()
	if out != nil {
		*out = r
	}
	if m.match(r) {
		p.rejects = p.rejects[0:0]
		return true
	}
	p.unreadrune()
	p.rejects = append(p.rejects, r)
	return false
}

func (p *parser) acceptall(m runematcher, limit int) []rune {
	result := []rune{}
	r := rune(0)
	for p.accept(m, &r) {
		result = append(result, r)
		if limit > 0 && len(result) > limit {
			return result
		}
	}
	return result
}

func (p *parser) eatspace() {
	r := rune(0)
	for p.accept(matchspace, &r) {
	}
}

func (p *parser) expectident() string {
	identch := []rune{p.expect(matchnamestart, "valid identifier")}

	r := rune(0)
	for p.accept(matchname, &r) {
		identch = append(identch, r)
	}

	return string(identch)
}

func (p *parser) expectstrlit() string {
	strlit := []rune{}

	p.expect(matchone('"'), "string")
	for {
		switch {
		case p.accept(matchone('&'), nil):
			switch {
			case p.accept(matchone('#'), nil):
				if p.accept(matchone('x'), nil) {
					hexesc := string(p.acceptall(matchhex, 6))
					v, err := strconv.ParseUint(hexesc, 16, 32)
					if err != nil {
						panic(err)
					}
					strlit = append(strlit, rune(v))
				} else {
					decesc := string(p.acceptall(matchhex, 7))
					v, err := strconv.ParseUint(decesc, 16, 32)
					if err != nil {
						panic(err)
					}
					strlit = append(strlit, rune(v))
				}
			default:
				entity := p.acceptall(matchalpha, 4)
				if len(entity) > 0 {
					if p.accept(matchone(';'), nil) {
						switch string(entity) {
						case "quot":
							strlit = append(strlit, rune('"'))
						case "apos":
							strlit = append(strlit, rune('\''))
						case "lt":
							strlit = append(strlit, rune('<'))
						case "gt":
							strlit = append(strlit, rune('>'))
						case "amp":
							strlit = append(strlit, rune('&'))
						default:
							// Treat as unescaped &.
							strlit = append(strlit, rune('&'))
							strlit = append(strlit, entity...)
						}
					} else {
						// Treat as unescaped &.
						strlit = append(strlit, rune('&'))
						strlit = append(strlit, entity...)
					}
				}
			}
		case p.accept(matchone('"'), nil):
			return string(strlit)
		default:
			strlit = append(strlit, p.readrune())
		}
	}
}

func (p *parser) expectcontent() string {
	contentstr := []rune{}

	for {
		switch {
		case p.accept(matchone('&'), nil):
			switch {
			case p.accept(matchone('#'), nil):
				if p.accept(matchone('x'), nil) {
					hexesc := string(p.acceptall(matchhex, 6))
					v, err := strconv.ParseUint(hexesc, 16, 32)
					if err != nil {
						panic(err)
					}
					contentstr = append(contentstr, rune(v))
				} else {
					decesc := string(p.acceptall(matchhex, 7))
					v, err := strconv.ParseUint(decesc, 16, 32)
					if err != nil {
						panic(err)
					}
					contentstr = append(contentstr, rune(v))
				}
			default:
				entity := p.acceptall(matchalpha, 4)
				if len(entity) > 0 {
					if p.accept(matchone(';'), nil) {
						switch string(entity) {
						case "quot":
							contentstr = append(contentstr, rune('"'))
						case "apos":
							contentstr = append(contentstr, rune('\''))
						case "lt":
							contentstr = append(contentstr, rune('<'))
						case "gt":
							contentstr = append(contentstr, rune('>'))
						case "amp":
							contentstr = append(contentstr, rune('&'))
						default:
							// Treat as unescaped &.
							contentstr = append(contentstr, rune('&'))
							contentstr = append(contentstr, entity...)
						}
					} else {
						// Treat as unescaped &.
						contentstr = append(contentstr, rune('&'))
						contentstr = append(contentstr, entity...)
					}
				}
			}
		case p.accept(matchone('<'), nil):
			p.unreadrune()
			return string(contentstr)
		default:
			contentstr = append(contentstr, p.readrune())
		}
	}
}

func (p *parser) expectattr() (string, string) {
	key := p.expectident()
	if p.accept(matchone('='), nil) {
		return key, p.expectstrlit()
	}
	return key, ""
}

func (p *parser) expectdtpart() DocumentInfo {
	var dt DocumentInfo
	p.eatspace()
	p.expectstr("xml")
	p.eatspace()
	for {
		switch {
		case p.accept(matchone('?'), nil):
			p.expect(matchone('>'), "doctype close")
			return dt
		default:
			key, value := p.expectattr()
			switch key {
			case "version":
				dt.Version = value
			case "encoding":
				dt.Encoding = value
			case "standalone":
				dt.Standalone = value
			default:
				continue
			}
			p.eatspace()
		}
	}
}

func (p *parser) expectstarttagpart(reader xmlreader) {
	p.eatspace()
	tag := p.expectident()
	reader.opentag(tag)
	p.eatspace()
	for {
		switch {
		case p.accept(matchone('>'), nil):
			return
		case p.accept(matchone('/'), nil):
			p.expect(matchone('>'), "end of tag")
			reader.closetag(tag)
			return
		default:
			key, value := p.expectattr()
			reader.attr(key, value)
			p.eatspace()
		}
	}
}

func (p *parser) expectendtagpart() string {
	p.eatspace()
	tag := p.expectident()
	p.eatspace()
	p.expect(matchone('>'), ">")
	return tag
}

func (p *parser) expectcommentpart() {
	// Start of comment
	p.expect(matchone('-'), "start of comment")
	p.expect(matchone('-'), "start of comment")

	wnd := [3]rune{p.readrune(), p.readrune(), p.readrune()}
	for {
		if wnd[0] == '-' && wnd[1] == '-' && wnd[2] == '>' {
			return
		}

		// Slide window.
		wnd[0] = wnd[1]
		wnd[1] = wnd[2]
		wnd[2] = p.readrune()
	}
}

type xmlreader interface {
	doctype(dt DocumentInfo)
	opentag(tag string)
	closetag(tag string)
	attr(key, value string)
	content(data string)
}

func (p *parser) parsenext(reader xmlreader) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	p.eatspace()

	switch {
	case p.accept(matchone('<'), nil):
		switch {
		case p.accept(matchone('!'), nil):
			p.expectcommentpart()
		case p.accept(matchone('/'), nil):
			reader.closetag(p.expectendtagpart())
		case p.accept(matchone('?'), nil):
			reader.doctype(p.expectdtpart())
		default:
			p.expectstarttagpart(reader)
		}
	default:
		reader.content(p.expectcontent())
	}

	return err
}
