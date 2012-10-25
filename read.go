package keyvalues

import ("io";"bufio";"fmt")

func FromReader(r io.Reader) (kv *KeyValues, err error) {
	p := parser{r: bufio.NewReader(r)}
	yyParse(&p)
	return p.kv, p.err
}

type parser struct {
	r    *bufio.Reader
	kv   *KeyValues
	err  error
	line int
}

func (p *parser) Error(err string) {
	if p.err == nil {
		p.err = fmt.Errorf("Error on line %d: %s", p.line, err)
	}
}

func (p *parser) Lex(y *yySymType) int {
	r, _, err := p.r.ReadRune()
	if err != nil {
		p.Error(err.Error())
		return yyErrCode
	}

	switch r {
	case ' ', '\t', '\r':
		return p.Lex(y)

	case '\n':
		p.line++
		return p.Lex(y)

	case '{', '}', '"':
		return int(r)

	case '/':
		if r, _, _ := p.r.ReadRune(); r == '/' {
			p.r.ReadString('\n')
			p.line++
			return p.Lex(y)
		}
		p.r.UnreadRune()

		fallthrough
	default:
		runes := []rune{r}
readloop:
		for err == nil {
			r, _, err = p.r.ReadRune()
			switch r {
			case '\n':
				p.line++
				fallthrough
			case ' ', '\t', '\r':
				break readloop
			case '{', '}', '"':
				p.r.UnreadRune()
				break readloop
			default:
				runes = append(runes, r)
			}
		}
		y.s = string(runes)
		return str
	}

	return 0
}
