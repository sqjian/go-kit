package rule

import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"sync"
)

var buildCache sync.Map

func CompileRule(code string) error {
	p, e := expr.Compile(code)
	if e != nil {
		return e
	}
	buildCache.Store(code, p)
	return nil
}

func ExecRule(code string, env any) (rst any, err error) {
	program, _ := buildCache.LoadOrStore(code, func() *vm.Program {
		_program, compileErr := expr.Compile(code)
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
