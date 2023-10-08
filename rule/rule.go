package rule

import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"reflect"
	"sync"
)

var buildCache sync.Map

type Option = expr.Option

func AsString() Option {
	return expr.AsKind(reflect.String)
}

func CompileRule(code string, ops ...Option) error {
	p, e := expr.Compile(code, ops...)
	if e != nil {
		return e
	}
	buildCache.Store(code, p)
	return nil
}

func ExecRule(code string, env any, ops ...Option) (rst any, err error) {
	program, _ := buildCache.LoadOrStore(code, func() *vm.Program {
		_program, compileErr := expr.Compile(code, ops...)
		if compileErr != nil {
			err = compileErr
		}
		return _program
	}())

	if program == nil || err != nil {
		return nil, err
	}

	return expr.Run(program.(*vm.Program), env)
}
