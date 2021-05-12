package dap

import (
	"fmt"
)

type Parser struct {
	Lexer
}

func (this *Parser) Parse() (
	imports map[string]*Import,
	classes map[string]*Class,
	methods map[string]*Method,
	functions map[string]*Function,
	err error,
) {
	imports = map[string]*Import{}
	classes = map[string]*Class{}
	methods = map[string]*Method{}
	functions = map[string]*Function{}

	for tok := this.peek(); tok.Typ != ttEOF; tok = this.peek() {
		var imp *Import
		imp, err = this.imprt()
		if err != nil {
			return
		}
		if imp != nil {
			imports[imp.Name] = imp
			continue
		}

		var cls *Class
		cls, err = this.class()
		if err != nil {
			return
		}
		if cls != nil {
			classes[cls.Name] = cls
			continue
		}

		var mthd *Method
		mthd, err = this.method()
		if err != nil {
			return
		}
		if mthd != nil {
			methods[mthd.Name] = mthd
			continue
		}

		var f *Function
		f, err = this.function()
		if err != nil {
			return
		}
		if f != nil {
			functions[f.Name] = f
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

func (this *Parser) get() (tok *Token) {
	for tok = this.Lexer.get(); ; tok = this.Lexer.get() {
		if tok.Typ == ttBlank {
			continue
		}
		if tok.Typ == ttLineEnd {
			continue
		}
		break
	}
	return
}

func (this *Parser) getex(linend bool) (tok *Token) {
	for tok = this.Lexer.get(); ; tok = this.Lexer.get() {
		if tok.Typ == ttBlank {
			continue
		}
		if !linend && tok.Typ == ttLineEnd {
			continue
		}
		break
	}
	return
}

func (this *Parser) peekex(linend bool) (tok *Token) {
	for tok = this.Lexer.peek(); ; _, tok = this.Lexer.get(), this.Lexer.peek() {
		if tok.Typ == ttBlank {
			continue
		}
		if !linend && tok.Typ == ttLineEnd {
			continue
		}
		break
	}
	return
}

func (this *Parser) peek() (tok *Token) {
	for tok = this.Lexer.peek(); ; _, tok = this.Lexer.get(), this.Lexer.peek() {
		if tok.Typ == ttBlank {
			continue
		}
		if tok.Typ == ttLineEnd {
			continue
		}
		break
	}
	return
}

func (this *Parser) classRef() (res *ClassRef, err error) {
	this.trans(func(p *Parser) bool {
		tok1 := p.getex(true)
		if tok1.Typ != ttSymbol {
			return false
		}
		tok2 := p.getex(true)
		if tok2.Typ != ttDot {
			res = &ClassRef{
				Name: tok1.Val,
			}
			return true
		}
		tok3 := p.getex(true)
		if tok3.Typ != ttSymbol {
			return false
		}
		res = &ClassRef{
			Pkg:  tok1.Val,
			Name: tok3.Val,
		}
		return true
	})
	return
}

func (this *Parser) arg() (res *Arg, err error) {
	this.trans(func(p *Parser) bool {
		ntok := p.get()
		if ntok.Typ != ttSymbol {
			return false
		}
		var cr *ClassRef
		cr, err = p.classRef()
		if err != nil {
			return false
		}
		res = &Arg{
			Name:  ntok.Val,
			Class: cr,
		}
		return true
	})
	return
}

func (this *Parser) args() (res Args, err error) {
	this.trans(func(p *Parser) (ok bool) {
		defer func() {
			if ok && res == nil {
				res = Args{}
			}
		}()
		for tok := p.peek(); tok.Typ != ttRightParenthese; tok = p.peek() {
			arg, _ := p.arg()
			if arg == nil {
				tok = p.peek()
				if tok.Typ == ttRightParenthese {
					p.get()
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
		var exprs []Express
		exprs, err = p.exprs(0)
		if err != nil {
			return false
		}
		if tok := p.get(); tok.Typ != ttRightCurve {
			err = fmt.Errorf("missing } after method")
			return false
		}
		res = &Method{
			Class: ctok.Val,
			Function: Function{
				Args:  args,
				Name:  ftok.Val,
				Exprs: exprs,
			},
		}
		return true
	})
	return
}

func (this *Parser) function() (res *Function, err error) {
	this.trans(func(p *Parser) bool {
		ftok := p.get()
		if ftok.Typ != ttSymbol {
			return false
		}
		if tok := p.get(); tok.Typ != ttLeftParenthese {
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
		var exprs []Express
		exprs, err = p.exprs(0)
		if err != nil {
			return false
		}
		if tok := p.get(); tok.Typ != ttRightCurve {
			err = fmt.Errorf("missing function }")
			return false
		}
		res = &Function{
			Name:  ftok.Val,
			Args:  args,
			Exprs: exprs,
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
			Path: tok.Val,
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
			Name:   nametok.Val,
			Fields: fields,
		}
		return true
	})
	return
}

func (this *Parser) classField() (f *Field, err error) {
	this.trans(func(p *Parser) bool {
		tok1 := p.getex(false)
		if tok1.Typ != ttSymbol {
			return false
		}
		var cr *ClassRef
		cr, err = p.classRef()
		if err != nil {
			return false
		}
		f = &Field{
			Name:  tok1.Val,
			Class: cr,
		}
		return true
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
			fields[f.Name] = f
		}
		return true
	})
	return
}

func (this *Parser) exprs(priority int) (res []Express, err error) {
	for {
		var expr Express
		expr, err = this.expr(priority)
		if expr == nil || err != nil {
			return
		}
		res = append(res, expr)
	}
}

func (this *Parser) lambda(priority int) (res Express, err error) {
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
		exprs, err = p.exprs(0)
		if err != nil {
			return false
		}
		if tok := p.get(); tok.Typ != ttRightCurve {
			err = fmt.Errorf("missing }")
			return false
		}
		res = &Lambda{
			Args:  args,
			Rets:  ret,
			Exprs: exprs,
		}
		return true
	})
	return
}

func (this *Parser) define(priority int) (res Express, err error) {
	this.trans(func(p *Parser) bool {
		if tok := p.get(); tok.Typ != ttVar {
			return false
		}
		ntok := p.get()
		if ntok.Typ != ttSymbol {
			err = fmt.Errorf("missing variable name")
			return false
		}
		var cr *ClassRef
		cr, err = p.classRef()
		if err != nil {
			return false
		}
		res = &ExprDefine{
			Name:  ntok.Val,
			Class: cr,
		}
		return true
	})
	return
}

func (this *Parser) ref(priority int) (res Express, err error) {
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
					Names: names,
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

func (this *Parser) defineOrRefOrReturn(priority int) (define Express, ref Express, ret *Token, err error) {
	this.trans(func(p *Parser) bool {
		define, err = p.define(priority)
		if err != nil {
			return false
		}
		if define != nil {
			return true
		}
		ref, err = p.ref(priority)
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

func (this *Parser) assign(priority int) (res Express, err error) {
	this.trans(func(p *Parser) bool {
		target := ExprAssignTarget{}
		for {
			var d Express
			var r Express
			var ret *Token
			d, r, ret, err = p.defineOrRefOrReturn(priority + 1)
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
		expr, err = p.expr(priority)
		if err != nil {
			return false
		}
		res = &ExprAssign{
			Target: target,
			Src:    expr,
		}
		return true
	})
	return
}

func (this *Parser) ifexpr(priority int) (res Express, err error) {
	this.trans(func(p *Parser) bool {
		priority++
		if tok := p.get(); tok.Typ != ttIf {
			return false
		}
		var expr1, expr2 Express
		expr1, err = p.expr(priority)
		if err != nil {
			return false
		}
		if expr1 == nil {
			return false
		}
		tok := p.get()
		if tok.Typ == ttSemi {
			expr2, err = p.expr(priority)
			if err != nil {
				return false
			}
			tok = p.get()
		}
		if tok.Typ != ttLeftCurve {
			err = fmt.Errorf("missing { after if condition")
			return false
		}
		var exprs []Express
		exprs, err = p.exprs(0)
		if err != nil {
			return false
		}
		tok = p.get()
		if tok.Typ != ttRightCurve {
			err = fmt.Errorf("missing } after if statements")
			return false
		}
		tok = p.peek()
		if tok.Typ != ttElse {
			if expr2 != nil {
				res = &ExprIf{
					Prepare: expr1,
					Cond:    expr2,
					Exprs:   exprs,
				}
				return true
			}
			res = &ExprIf{
				Cond:  expr1,
				Exprs: exprs,
			}
			return true
		}
		p.get()
		tok = p.peek()
		if tok.Typ == ttIf {
			var el Express
			el, err = p.ifexpr(priority)
			if err != nil {
				return false
			}
			if expr2 != nil {
				res = &ExprIf{
					Prepare: expr1,
					Cond:    expr2,
					Exprs:   exprs,
					Else:    el,
				}
				return true
			}
			res = &ExprIf{
				Cond:  expr1,
				Exprs: exprs,
				Else:  el,
			}
			return true
		}
		tok = p.get()
		if tok.Typ != ttLeftCurve {
			err = fmt.Errorf("missing {")
			return false
		}
		exprs, err = p.exprs(0)
		if err != nil {
			return false
		}
		if tok.Typ != ttRightCurve {
			err = fmt.Errorf("missing }")
			return false
		}
		el := &ExprIf{
			Exprs: exprs,
		}
		if expr2 != nil {
			res = &ExprIf{
				Prepare: expr1,
				Cond:    expr2,
				Exprs:   exprs,
				Else:    el,
			}
			return true
		}
		res = &ExprIf{
			Cond:  expr1,
			Exprs: exprs,
			Else:  el,
		}
		return true
	})
	return
}

func (this *Parser) expradd(priority int) (res Express, err error) {
	this.trans(func(p *Parser) bool {
		var first, second Express
		first, err = p.expr(priority + 1)
		if err != nil {
			return false
		}
		if first == nil {
			return false
		}
		if tok := p.get(); tok.Typ != ttAdd {
			return false
		}
		second, err = p.expr(priority)
		if err != nil {
			return false
		}
		if second == nil {
			return false
		}
		res = &ExprAdd{
			First:  first,
			Second: second,
		}
		return true
	})
	return
}

func (this *Parser) exprsub(priority int) (res Express, err error) {
	this.trans(func(p *Parser) bool {
		var first, second Express
		first, err = p.expr(priority + 1)
		if err != nil {
			return false
		}
		if first == nil {
			return false
		}
		if tok := p.get(); tok.Typ != ttSub {
			return false
		}
		second, err = p.expr(priority)
		if err != nil {
			return false
		}
		if second == nil {
			return false
		}
		res = &ExprSub{
			First:  first,
			Second: second,
		}
		return true
	})
	return
}

func (this *Parser) exprmulti(priority int) (res Express, err error) {
	this.trans(func(p *Parser) bool {
		var first, second Express
		first, err = p.expr(priority + 1)
		if err != nil {
			return false
		}
		if first == nil {
			return false
		}
		if tok := p.get(); tok.Typ != ttMulti {
			return false
		}
		second, err = p.expr(priority)
		if err != nil {
			return false
		}
		if second == nil {
			return false
		}
		res = &ExprMulti{
			First:  first,
			Second: second,
		}
		return true
	})
	return
}

func (this *Parser) exprdiv(priority int) (res Express, err error) {
	this.trans(func(p *Parser) bool {
		var first, second Express
		first, err = p.expr(priority + 1)
		if err != nil {
			return false
		}
		if first == nil {
			return false
		}
		if tok := p.get(); tok.Typ != ttDiv {
			return false
		}
		second, err = p.expr(priority)
		if err != nil {
			return false
		}
		if second == nil {
			return false
		}
		res = &ExprDiv{
			First:  first,
			Second: second,
		}
		return true
	})
	return
}

func (this *Parser) instnum(priority int) (res Express, err error) {
	this.trans(func(p *Parser) bool {
		tok := p.get()
		if tok.Typ != ttConstNumber {
			return false
		}
		res = &ExpressLiteralNumber{
			Val: tok.Val,
		}
		return true
	})
	return
}

type PriExpr struct {
	pri int
	f   func(p *Parser, priority int) (Express, error)
}

var _exprfuncs []*PriExpr

func init() {
	_exprfuncs = []*PriExpr{
		{0, (*Parser).assign},
		{5, (*Parser).define},
		{10, (*Parser).ifexpr},
		{15, (*Parser).exprCmp},
		{20, (*Parser).lambda},
		{30, (*Parser).expradd},
		{30, (*Parser).exprsub},
		{40, (*Parser).exprmulti},
		{40, (*Parser).exprdiv},
		{90, (*Parser).instnum},
		{100, (*Parser).ref},
	}
}

func (this *Parser) exprCmp(priority int) (res Express, err error) {
	this.trans(func(p *Parser) bool {
		var expr1 Express
		expr1, err = p.expr(priority + 1)
		if err != nil {
			return false
		}
		if expr1 == nil {
			return false
		}
		var expr2 Express
		switch tok := p.get(); tok.Typ {
		case ttGT, ttGTE, ttLT, ttLTE:
			if expr2, err = p.expr(priority); err == nil && expr2 != nil {
				res = &ExprCmp{
					OP:     tok.Typ,
					First:  expr1,
					Second: expr2,
				}
				return true
			}
		}
		return false
	})
	return
}

func (this *Parser) expr(priority int) (res Express, err error) {
	this.trans(func(p *Parser) bool {
		for _, pf := range _exprfuncs {
			if pf.pri < priority {
				continue
			}
			res, err = pf.f(p, pf.pri)
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
