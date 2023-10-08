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

func ExecRule(code string, env any) (any, error) {
	p, o := buildCache.Load(code)
	if !o {
		return nil, errWrapper(NotFound)
	}
	return expr.Run(p.(*vm.Program), env)
}

func EvalRule(code string, env any) (any, error) {
	program, _ := buildCache.LoadOrStore(code, func() *vm.Program {
		_program, compileErr := expr.Compile(code)
		if compileErr != nil {
			panic(compileErr)
		}
		return _program
	}())
	return expr.Run(program.(*vm.Program), env)
}
