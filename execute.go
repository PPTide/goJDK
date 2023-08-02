package main

// https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-2.html#jvms-2.6.1
type variable interface {
	implementVariable()
}

type boolean struct {
	val bool
}

func (b boolean) implementVariable() {
	return
}

// https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-2.html#jvms-2.6
type frame struct {
	localVariables []variable
	operandStack   interface{}
	constantPool   *interface{}
}

type runTimeData struct {
	pc                int
	stack             []frame
	heap              interface{}
	methodArea        interface{}
	constantPool      interface{}
	nativeMethodStack interface{}
}
