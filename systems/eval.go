/// eval implements the rules engine for the mainland server. Any code being ran
/// is first hashed and looked up within a cache table. This prevents the excessive overhead
/// of recompiling the code every time it's ran. (Hopefully enough for real-time use.)
package systems

import (
	"github.com/d5/tengo/v2"
)

type EvalSystem struct {
	compiledCode map[string]*tengo.Compiled
	moduleMap    *tengo.ModuleMap
}

func NewEvalSystem() *EvalSystem {
	engine := new(EvalSystem)
	engine.compiledCode = make(map[string]*tengo.Compiled)
	engine.moduleMap = tengo.NewModuleMap()
	engine.moduleMap.AddBuiltinModule("engine", EngineModule())
	return engine
}

func (engine *EvalSystem) CacheCode(code string) (*tengo.Compiled, error) {
	script := tengo.NewScript([]byte(code))
	script.SetImports(engine.moduleMap)
	compiled, err := script.Compile()
	if err != nil {
		return nil, err
	}
	engine.compiledCode[code] = compiled
	return compiled, nil
}

func (engine *EvalSystem) RunCode(code string) error {
	compiled, ok := engine.compiledCode[code]
	if !ok {
		new_compiled, err := engine.CacheCode(code)
		compiled = new_compiled
		if err != nil {
			return err
		}
	}
	return compiled.Run()
}
