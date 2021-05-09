package dap

type Parser struct {
	imports map[string]*Import
	classes map[string]*Class
	lx      Lexer
}

func (p *Parser) parse() {
	p.imprt()
	p.class()
	p.method()
}
func (p *Parser) method() {}

func (p *Parser) imprt() *Import {
	tok := p.lx.peek()
	if tok.typ != ttImport {
		return nil
	}
	tok = p.lx.fetch()
	if tok.typ != ttConstString {
		panic("should have import path")
	}
	return &Import{
		path: tok.val,
	}
}

func (p *Parser) class() (c *Class) {
	p.lx.trans(func(l Lexer) bool {
		tok := l.get()
		if tok.typ != ttSymbol {
			return false
		}
		c = &Class{}
		c.name = tok.val
		tok = l.get()
		if tok.typ != ttLeftBrace {
			return false
		}
		c.fields = p.classBody()
		if tok.typ != ttRightBrace {
			return false
		}
		return true
	})
	return
}

func (p *Parser) classField() (f *Field) {
	p.lx.trans(func(l Lexer) bool {
		tok1 := l.get()
		if tok1.typ != ttSymbol {
			return false
		}
		tok2 := l.get()
		if tok2.typ == ttLineEnd {
			f = &Field{
				name:  tok1.val,
				class: "",
			}
			return true
		}
		if tok2.typ == ttSymbol {
			f = &Field{
				name:  tok1.val,
				class: tok2.val,
			}
			return true
		}
		return false
	})
	return
}

func (p *Parser) classBody() (fields map[string]*Field) {
	fields = make(map[string]*Field)
	p.lx.trans(func(l Lexer) bool {
		p.classField()
		return true
	})
	return
}
