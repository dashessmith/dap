package dap

import "fmt"

type Parser struct {
	Lexer
}

func (this *Parser) parse() (
	imports map[string]*Import,
	classes map[string]*Class,
	methods map[string]*Method,
	err error,
) {
	imports = map[string]*Import{}
	classes = map[string]*Class{}
	methods = map[string]*Method{}
	for tok := this.peek(); tok.Typ != ttEOF; tok = this.peek() {
		var imp *Import
		imp, err = this.imprt()
		if err != nil {
			return
		}
		if imp != nil {
			imports[imp.name] = imp
			continue
		}

		var cls *Class
		cls, err = this.class()
		if err != nil {
			return
		}
		if cls != nil {
			classes[cls.name] = cls
			continue
		}

		var mthd *Method
		mthd, err = this.method()
		if err != nil {
			return
		}
		if mthd != nil {
			methods[mthd.name] = mthd
			continue
		}
		err = fmt.Errorf("unknown token %s", tok)
		return
	}
	return
}

func (this *Parser) trans(f func(p *Parser) bool) {
	backup := this.Lexer
	defer func() {
		this.Lexer = backup
	}()
	this.Lexer = this.begin()
	defer this.done()
	if f(this) {
		this.commit()
	}
}

func (this *Parser) arg() (res *Arg, err error) {
	this.trans(func(p *Parser) bool {
		ntok := p.get()
		if ntok.Typ != ttSymbol {
			return false
		}
		tok2 := p.peek()
		if tok2.Typ != ttSymbol {
			res = &Arg{
				name: ntok.Val,
			}
			return true
		}
		p.get()
		res = &Arg{
			name:  ntok.Val,
			class: tok2.Val,
		}
		return true
	})
	return
}

func (this *Parser) args() (res Args, err error) {
	this.trans(func(p *Parser) bool {
		for tok := p.peek(); tok.Typ != ttRightParenthese; tok = p.peek() {
			arg, _ := p.arg()
			if arg == nil {
				tok = p.peek()
				if tok.Typ == ttRightParenthese {
					p.get()
					if res == nil {
						res = Args{}
					}
					return true
				}
				res = nil
				return false
			}
			res = append(res, arg)
			tok = p.peek()
			if tok.Typ == ttComma {
				p.get()
				continue
			}
			if tok.Typ == ttRightParenthese {
				return true
			}
			err = fmt.Errorf("parse args failed")
			return false
		}
		return true
	})
	return
}

func (this *Parser) argsWithParenthese() (res Args, err error) {
	this.trans(func(p *Parser) bool {
		if tok := p.get(); tok.Typ != ttLeftParenthese {
			return false
		}
		res, err = p.args()
		if err != nil {
			return false
		}
		if tok := p.get(); tok.Typ != ttRightParenthese {
			err = fmt.Errorf("missing )")
			return false
		}
		return true
	})
	return
}

func (this *Parser) method() (res *Method, err error) {
	this.trans(func(p *Parser) bool {
		ctok := p.get()
		if ctok.Typ != ttSymbol {
			return false
		}
		if dtok := p.get(); dtok.Typ != ttDot {
			return false
		}
		ftok := p.get()
		if ftok.Typ != ttSymbol {
			err = fmt.Errorf("missing method name")
			return false
		}
		if tok := p.get(); tok.Typ != ttLeftParenthese {
			err = fmt.Errorf("missing (")
			return false
		}
		var args Args
		args, err = p.args()
		if err != nil {
			return false
		}
		if tok := p.get(); tok.Typ != ttRightParenthese {
			err = fmt.Errorf("missing )")
			return false
		}
		if tok := p.get(); tok.Typ != ttLeftCurve {
			err = fmt.Errorf("missing {")
			return false
		}
		if tok := p.get(); tok.Typ != ttRightCurve {
			err = fmt.Errorf("missing }")
			return false
		}
		res = &Method{
			class: ctok.Val,
			Function: Function{
				args: args,
				name: ftok.Val,
			},
		}
		return true
	})
	return
}

func (this *Parser) imprt() (res *Import, err error) {
	this.trans(func(p *Parser) bool {
		tok := p.get()
		if tok.Typ != ttImport {
			return false
		}
		tok = p.get()
		if tok.Typ != ttConstString {
			panic("should have import path")
		}
		res = &Import{
			path: tok.Val,
		}
		return true
	})
	return
}

func (this *Parser) class() (c *Class, err error) {
	this.trans(func(p *Parser) bool {
		nametok := p.get()
		if nametok.Typ != ttSymbol {
			return false
		}
		leftblock := p.get()
		if leftblock.Typ != ttLeftCurve {
			return false
		}
		var fields map[string]*Field
		fields, err = p.classBody()
		if err != nil {
			return false
		}
		rightblock := p.get()
		if rightblock.Typ != ttRightCurve {
			err = fmt.Errorf("need right }")
			return false
		}
		c = &Class{
			name:   nametok.Val,
			fields: fields,
		}
		return true
	})
	return
}

func (this *Parser) classField() (f *Field, err error) {
	this.trans(func(p *Parser) bool {
		tok1 := p.get()
		if tok1.Typ != ttSymbol {
			err = fmt.Errorf("need field name")
			return false
		}
		tok2 := p.get()
		if tok2.Typ == ttLineEnd {
			f = &Field{
				name:  tok1.Val,
				class: "",
			}
			return true
		}
		if tok2.Typ == ttSymbol {
			f = &Field{
				name:  tok1.Val,
				class: tok2.Val,
			}
			return true
		}
		err = fmt.Errorf("invalid field format")
		return false
	})
	return
}

func (this *Parser) classBody() (fields map[string]*Field, err error) {
	fields = map[string]*Field{}
	this.trans(func(p *Parser) bool {
		for tok := p.peek(); tok.Typ != ttRightCurve; tok = p.peek() {
			var f *Field
			f, err = p.classField()
			if err != nil {
				return false
			}
			if f == nil {
				return true
			}
			fields[f.name] = f
		}
		return true
	})
	return
}

func (this *Parser) exprs() (res []Express, err error) {
	return
}

func (this *Parser) lambda() (res Express, err error) {
	this.trans(func(p *Parser) bool {
		var args Args
		args, err = p.argsWithParenthese()
		if err != nil {
			return false
		}
		if args == nil {
			return false
		}

		var ret Args
		ret, err = p.argsWithParenthese()
		if err != nil {
			return false
		}

		if tok := p.get(); tok.Typ != ttLeftCurve {
			err = fmt.Errorf("missing {")
			return false
		}
		var exprs []Express
		exprs, err = p.exprs()
		if err != nil {
			return false
		}
		if tok := p.get(); tok.Typ != ttRightCurve {
			err = fmt.Errorf("missing }")
			return false
		}
		res = &Lambda{
			args:  args,
			ret:   ret,
			exprs: exprs,
		}
		return true
	})
	return
}

func (this *Parser) define() (res Express, err error) {
	this.trans(func(p *Parser) bool {
		if tok := p.get(); tok.Typ != ttVar {
			return false
		}
		ntok := p.get()
		if ntok.Typ != ttSymbol {
			err = fmt.Errorf("missing variable name")
			return false
		}
		tok2 := p.peek()
		if tok2.Typ == ttSymbol {
			p.get()
			res = &ExprDefine{
				name:  ntok.Val,
				class: tok2.Val,
			}
			return true
		}
		res = &ExprDefine{
			name: ntok.Val,
		}
		return true
	})
	return
}

func (this *Parser) ref() (res Express, err error) {
	this.trans(func(p *Parser) bool {
		names := []string{}
		tok := p.get()
		if tok.Typ != ttSymbol {
			return false
		}
		names = append(names, tok.Val)

		for {
			if dtok := p.peek(); dtok.Typ != ttDot {
				res = &ExprRef{
					names: names,
				}
				return true
			}
			p.get()
			tok = p.get()
			if tok.Typ != ttSymbol {
				err = fmt.Errorf("invalid reference")
				return false
			}
			names = append(names, tok.Val)
		}
	})
	return
}

func (this *Parser) defineOrRefOrReturn() (define Express, ref Express, ret *Token, err error) {
	this.trans(func(p *Parser) bool {
		define, err = p.define()
		if err != nil {
			return false
		}
		if define != nil {
			return true
		}
		ref, err = p.ref()
		if err != nil {
			return false
		}
		if ref != nil {
			return true
		}
		if tok := p.get(); tok.Typ == ttReturn {
			ret = tok
			return true
		}
		return false
	})
	return
}

func (this *Parser) assign() (res Express, err error) {
	this.trans(func(p *Parser) bool {
		target := ExprAssignTarget{}
		for {
			var d Express
			var r Express
			var ret *Token
			d, r, ret, err = p.defineOrRefOrReturn()
			if err != nil {
				return false
			}
			if r != nil {
				target = append(target, r)
				continue
			}
			if d != nil {
				target = append(target, d)
				continue
			}
			if ret != nil {
				target = append(target, ret)
				continue
			}
			break
		}

		if tok := p.get(); tok.Typ != ttAssign {
			return false
		}

		var expr Express
		expr, err = p.expr()
		if err != nil {
			return false
		}
		res = &ExprAssign{
			target: target,
			src:    expr,
		}
		return true
	})
	return
}

func (this *Parser) ifexpr() (res Express, err error) {
	this.trans(func(p *Parser) bool {
		if tok := p.get(); tok.Typ != ttIf {
			return false
		}
		var expr1, expr2 Express
		expr1, err = p.expr()
		if err != nil {
			return false
		}
		if expr1 == nil {
			return false
		}
		tok := p.get()
		if tok.Typ == ttSemi {
			expr2, err = p.expr()
			if err != nil {
				return false
			}
		}
		tok = p.get()
		if tok.Typ != ttLeftCurve {
			err = fmt.Errorf("missing {")
			return false
		}
		var exprs []Express
		exprs, err = p.exprs()
		if err != nil {
			return false
		}
		if tok.Typ != ttRightCurve {
			err = fmt.Errorf("missing }")
			return false
		}
		tok = p.peek()
		if tok.Typ != ttElse {
			if expr2 != nil {
				res = &ExprIf{
					prepare:   expr1,
					condition: expr2,
					exprs:     exprs,
				}
				return true
			}
			res = &ExprIf{
				condition: expr1,
				exprs:     exprs,
			}
			return true
		}
		p.get()
		tok = p.peek()
		if tok.Typ == ttIf {
			var el Express
			el, err = p.ifexpr()
			if err != nil {
				return false
			}
			if expr2 != nil {
				res = &ExprIf{
					prepare:   expr1,
					condition: expr2,
					exprs:     exprs,
					el:        el,
				}
				return true
			}
			res = &ExprIf{
				condition: expr1,
				exprs:     exprs,
				el:        el,
			}
			return true
		}
		tok = p.get()
		if tok.Typ != ttLeftCurve {
			err = fmt.Errorf("missing {")
			return false
		}
		exprs, err = p.exprs()
		if err != nil {
			return false
		}
		if tok.Typ != ttRightCurve {
			err = fmt.Errorf("missing }")
			return false
		}
		el := &ExprIf{
			exprs: exprs,
		}
		if expr2 != nil {
			res = &ExprIf{
				prepare:   expr1,
				condition: expr2,
				exprs:     exprs,
				el:        el,
			}
			return true
		}
		res = &ExprIf{
			condition: expr1,
			exprs:     exprs,
			el:        el,
		}
		return true
	})
	return
}

func (this *Parser) expradd() (res Express, err error) {
	this.trans(func(p *Parser) bool {
		var first, second Express
		first, err = p.expr()
		if err != nil {
			return false
		}
		if first == nil {
			return false
		}
		if tok := p.get(); tok.Typ != ttAdd {
			return false
		}
		second, err = p.expr()
		if err != nil {
			return false
		}
		if second == nil {
			return false
		}
		res = &ExprAdd{
			first:  first,
			second: second,
		}
		return true
	})
	return
}

func (this *Parser) exprsub() (res Express, err error) {
	this.trans(func(p *Parser) bool {
		var first, second Express
		first, err = p.expr()
		if err != nil {
			return false
		}
		if first == nil {
			return false
		}
		if tok := p.get(); tok.Typ != ttSub {
			return false
		}
		second, err = p.expr()
		if err != nil {
			return false
		}
		if second == nil {
			return false
		}
		res = &ExprSub{
			first:  first,
			second: second,
		}
		return true
	})
	return
}

func (this *Parser) exprmulti() (res Express, err error) {
	this.trans(func(p *Parser) bool {
		var first, second Express
		first, err = p.expr()
		if err != nil {
			return false
		}
		if first == nil {
			return false
		}
		if tok := p.get(); tok.Typ != ttMulti {
			return false
		}
		second, err = p.expr()
		if err != nil {
			return false
		}
		if second == nil {
			return false
		}
		res = &ExprMulti{
			first:  first,
			second: second,
		}
		return true
	})
	return
}

func (this *Parser) exprdiv() (res Express, err error) {
	this.trans(func(p *Parser) bool {
		var first, second Express
		first, err = p.expr()
		if err != nil {
			return false
		}
		if first == nil {
			return false
		}
		if tok := p.get(); tok.Typ != ttDiv {
			return false
		}
		second, err = p.expr()
		if err != nil {
			return false
		}
		if second == nil {
			return false
		}
		res = &ExprDiv{
			first:  first,
			second: second,
		}
		return true
	})
	return
}

var _exprfuncs []func(*Parser) (Express, error)

func init() {
	_exprfuncs = []func(*Parser) (Express, error){
		(*Parser).lambda,
		(*Parser).assign,
		(*Parser).ifexpr,
		(*Parser).exprmulti,
		(*Parser).exprdiv,
		(*Parser).expradd,
		(*Parser).exprsub,
		(*Parser).ref,
		(*Parser).define,
	}
}

func (this *Parser) expr() (res Express, err error) {
	this.trans(func(p *Parser) bool {
		for _, f := range _exprfuncs {
			res, err = f(p)
			if err != nil {
				return false
			}
			if res != nil {
				return true
			}
		}
		return false
	})
	return
}
