package opcodes

import (
	"fmt"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
)

// RegisterPackage registers a package on the CX standard library. This does not create a `CXPackage` structure,
// it only tells the CX runtime that `pkgName` will exist by the time a CX program is run.
//func RegisterPackage(pkgName string) {
//	constants.CorePackages = append(constants.CorePackages, pkgName)
//}

// GetOpCodeCount returns an op code that is available for usage on the CX standard library.
/*
func GetOpCodeCount() int {
	return len(OpcodeHandlers)
}
*/

// RegisterOperator ...
func RegisterOperator(code int, name string, handler ast.OpcodeHandler, inputs []*ast.CXArgument, outputs []*ast.CXArgument, atomicType int, operator int) {
	RegisterOpCode(code, name, handler, inputs, outputs)
	native := ast.Natives[code]
	ast.Operators[ast.GetTypedOperatorOffset(atomicType, operator)] = native
}

// MakeNativeFunction ...
func MakeNativeFunction(opCode int, inputs []*ast.CXArgument, outputs []*ast.CXArgument) *ast.CXFunction {
	fn := &ast.CXFunction{
		IsBuiltin: true,
		OpCode:   opCode,
	}

	offset := 0
	for _, inp := range inputs {
		inp.Offset = offset
		offset += ast.GetSize(inp)
		fn.Inputs = append(fn.Inputs, inp)
	}
	for _, out := range outputs {
		fn.Outputs = append(fn.Outputs, out)
		out.Offset = offset
		offset += ast.GetSize(out)
	}

	return fn
}

// RegisterOpCode ...
func RegisterOpCode(code int, name string, handler ast.OpcodeHandler, inputs []*ast.CXArgument, outputs []*ast.CXArgument) {
	if code >= len(ast.OpcodeHandlers) {
		ast.OpcodeHandlers = append(ast.OpcodeHandlers, make([]ast.OpcodeHandler, code+1)...)
	}
	if ast.OpcodeHandlers[code] != nil {
		panic(fmt.Sprintf("duplicate opcode %d : '%s' width '%s'.\n", code, name, ast.OpNames[code]))
	}
	ast.OpcodeHandlers[code] = handler

	ast.OpNames[code] = name
	ast.OpCodes[name] = code
	//OpVersions[code] = 2

	if inputs == nil {
		inputs = []*ast.CXArgument{}
	}
	if outputs == nil {
		outputs = []*ast.CXArgument{}
	}
	ast.Natives[code] = MakeNativeFunction(code, inputs, outputs)
}

/*
// Debug helper function used to find opcodes when they are not registered
func dumpOpCodes(opCode int) {
	var keys []int
	for k := range ast.OpNames {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Printf("%5d : %s\n", k, ast.OpNames[k])
	}

	fmt.Printf("opCode : %d\n", opCode)
}*/

// Struct helper for creating a struct parameter. It creates a
// `ast.CXArgument` named `argName`, that represents a structure instance of
// `strctName`, from package `pkgName`.
func Struct(pkgName, strctName, argName string) *ast.CXArgument {
	pkg, err := ast.PROGRAM.GetPackage(pkgName)
	if err != nil {
		panic(err)
	}

	strct, err := pkg.GetStruct(strctName)
	if err != nil {
		panic(err)
	}

	arg := ast.MakeArgument(argName, "", -1).AddType(constants.TypeNames[constants.TYPE_CUSTOM])
	arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_STRUCT)
	arg.Size = strct.Size
	arg.TotalSize = strct.Size
	arg.CustomType = strct

	return arg
}

//TODO: Rename OP_DEBUG, OP_DEBUG_PRINT_STACK
func opDebug([]ast.CXValue, []ast.CXValue) {
	ast.PROGRAM.PrintStack()
}
