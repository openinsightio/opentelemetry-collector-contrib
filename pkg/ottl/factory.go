// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"

import "go.opentelemetry.io/collector/component"

// Arguments holds the arguments for an OTTL function, with arguments
// specified as fields on a struct. Argument ordering is defined
type Arguments interface{}

// FunctionContext contains data provided by the Collector
// component to the OTTL for use in functions.
type FunctionContext struct {
	Set component.TelemetrySettings
}

// Factory defines an OTTL function factory that will generate an OTTL
// function to be called within a statement.
type Factory[K any] interface {
	// Name is the canonical name to be used by the user when invocating
	// the function generated by this Factory.
	Name() string

	// CreateDefaultArguments initializes an Arguments struct specific to this
	// Factory containing the arguments for the function.
	CreateDefaultArguments() Arguments

	// CreateFunction creates an OTTL function that will use the given Arguments.
	CreateFunction(fCtx FunctionContext, args Arguments) (ExprFunc[K], error)

	// Disallow implementations outside this package.
	unexportedFactoryFunc()
}

type CreateFunctionFunc[K any] func(fCtx FunctionContext, args Arguments) (ExprFunc[K], error)

type factory[K any] struct {
	name               string
	args               Arguments
	createFunctionFunc CreateFunctionFunc[K]
}

// nolint:unused
func (f *factory[K]) unexportedFactoryFunc() {}

func (f *factory[K]) Name() string {
	return f.name
}

func (f *factory[K]) CreateDefaultArguments() Arguments {
	return f.args
}

func (f *factory[K]) CreateFunction(fCtx FunctionContext, args Arguments) (ExprFunc[K], error) {
	return f.createFunctionFunc(fCtx, args)
}

type FactoryOption[K any] func(factory *factory[K])

func NewFactory[K any](name string, args Arguments, createFunctionFunc CreateFunctionFunc[K], options ...FactoryOption[K]) Factory[K] {
	f := &factory[K]{
		name:               name,
		args:               args,
		createFunctionFunc: createFunctionFunc,
	}

	for _, option := range options {
		option(f)
	}

	return f
}

// CreateFactoryMap takes a list of factories and returns a map of Factories
// keyed on their canonical names.
func CreateFactoryMap[K any](factories ...Factory[K]) map[string]Factory[K] {
	factoryMap := map[string]Factory[K]{}

	for _, fn := range factories {
		factoryMap[fn.Name()] = fn
	}

	return factoryMap
}
