package stack

import (
	"errors"

	"github.com/robertkrimen/otto"
	"github.com/yamakiller/magicLibs/files"
)

var (
	// ErrJSNotFindFile :
	ErrJSNotFindFile = errors.New("script file does not exist")
	// ErrJSNotFileData :
	ErrJSNotFileData = errors.New("did not get file data")
)

//JSStack desc
//@struct JSStack desc: javascirpt virtual machine
type JSStack struct {
	_state *otto.Otto
}

//MakeJSStack desc
//@method MakeJSStack desc: create javascript virtual machine
func MakeJSStack() *JSStack {
	return &JSStack{otto.New()}
}

//SetInt desc
//@method SetInt desc: Set the Int variable to the JS script
//@param (string) name
//@param (int)    value
func (slf *JSStack) SetInt(name string, val int) {
	slf._state.Set(name, val)
}

//SetFloat desc
//@method SetFloat desc: Set the Float 32 variable to the JS script
//@param (string)     name
//@param (float32)    value
func (slf *JSStack) SetFloat(name string, val float32) {
	slf._state.Set(name, val)
}

//SetDouble desc
//@method SetDouble desc: Set the Float 64 variable to the JS script
//@param (string)     name
//@param (float64)    value
func (slf *JSStack) SetDouble(name string, val float64) {
	slf._state.Set(name, val)
}

//SetBoolean desc
//@method SetBoolean desc: Set Bool variables for JS scripts
//@param (string)     name
//@param (bool)       value
func (slf *JSStack) SetBoolean(name string, val bool) {
	slf._state.Set(name, val)
}

//SetString desc
//@method SetString desc: Set String variables for JS scripts
//@param (string)     name
//@param (string)       value
func (slf *JSStack) SetString(name string, val string) {
	slf._state.Set(name, val)
}

//SetFunc desc
//@method SetFunc desc: Set the js script to call Go's function
//@param (string)       name
//@param (interface{})  value
func (slf *JSStack) SetFunc(name string, fun interface{}) {
	slf._state.Set(name, fun)
}

//ExecuteScriptFile desc
//@method ExecuteScriptFile desc: Execution script file
//@param   (string) scirpt file path
//@return  (otto.Value) javascript result
//@return  (error) javascript execution error result
func (slf *JSStack) ExecuteScriptFile(filename string) (otto.Value, error) {
	tmpFileName := files.Instance().GetFullPathForFilename(filename)
	if !files.Instance().IsFileExist(tmpFileName) {
		return otto.Value{}, ErrJSNotFindFile
	}

	data, err := files.Instance().GetDataFromCacheFile(tmpFileName)
	if err != nil {
		return otto.Value{}, ErrJSNotFileData
	}

	return slf._state.Run(string(data.GetBytes()))
}
