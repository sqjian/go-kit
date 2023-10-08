package rule

import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"sync"
)

var buildCache sync.Map

func Compile(code string) (*vm.Program, error) {
	return expr.Compile(code)
}

func Run(program *vm.Program, env any) (any, error) {
	return expr.Run(program, env)
}

func ExecRule(code string, env any) (any, error) {
	program, _ := buildCache.LoadOrStore(code, func() *vm.Program {
		_program, compileErr := Compile(code)
		if compileErr != nil {
			panic(compileErr)
		}
		return _program
	}())
	return Run(program.(*vm.Program), env)
}
