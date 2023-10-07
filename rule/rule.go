package rule

import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"sync"
)

var buildCache sync.Map

func ExecRule(code string, env any) (any, error) {
	program, _ := buildCache.LoadOrStore(code, func() *vm.Program {
		_program, compileErr := expr.Compile(code)
		if compileErr != nil {
			panic(compileErr)
		}
		return _program
	}())
	return expr.Run(program.(*vm.Program), env)
}
