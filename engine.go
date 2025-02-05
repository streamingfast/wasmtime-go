package wasmtime

// #include <wasmtime.h>
import "C"
import (
	"runtime"
	"sync"
)

// Engine is an instance of a wasmtime engine which is used to create a `Store`.
//
// Engines are a form of global configuration for wasm compilations and modules
// and such.
type Engine struct {
	_ptr *C.wasm_engine_t
	once sync.Once
}

// NewEngine creates a new `Engine` with default configuration.
func NewEngine() *Engine {
	engine := &Engine{_ptr: C.wasm_engine_new()}
	runtime.SetFinalizer(engine, func(engine *Engine) {
		engine.FreeMem()
	})
	return engine
}

// NewEngineWithConfig creates a new `Engine` with the `Config` provided
//
// Note that once a `Config` is passed to this method it cannot be used again.
func NewEngineWithConfig(config *Config) *Engine {
	if config.ptr() == nil {
		panic("config already used")
	}
	engine := &Engine{_ptr: C.wasm_engine_new_with_config(config.ptr())}
	runtime.SetFinalizer(config, nil)
	config._ptr = nil
	runtime.SetFinalizer(engine, func(engine *Engine) {
		C.wasm_engine_delete(engine._ptr)
	})
	return engine
}

func (engine *Engine) ptr() *C.wasm_engine_t {
	ret := engine._ptr
	maybeGC()
	return ret
}

func (engine *Engine) FreeMem() {
	engine.once.Do(func() {
		C.wasm_engine_delete(engine._ptr)
	})
}

// IncrementEpoch will increase the current epoch number by 1 within the
// current engine which will cause any connected stores with their epoch
// deadline exceeded to now be interrupted.
//
// This method is safe to call from any goroutine.
func (engine *Engine) IncrementEpoch() {
	C.wasmtime_engine_increment_epoch(engine.ptr())
	runtime.KeepAlive(engine)
}
