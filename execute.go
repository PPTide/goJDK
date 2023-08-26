package main

import (
	"bytes"
	"fmt"
	"github.com/PPTide/gojdk/parse"
	"strings"
)

const (
	methodFlagAccPublic = 0x0001
	methodFlagAccStatic = 0x0008
	methodFlagAccNative = 0x0100
)

type variable struct {
	valType       string
	val           interface{}
	referenceType string
	reference     *interface{}
}

func (v variable) expectType(varType string) interface{} {
	if v.valType != varType {
		panic(fmt.Errorf("type %s didn't match expected type %s", v.valType, varType))
	}
	return v.val
}

func asIntVariable(val int) variable {
	return variable{
		valType: "int",
		val:     val,
	}
}

func (v variable) expectReferenceOfType(referenceType string) *interface{} {
	if v.referenceType != referenceType {
		panic(fmt.Errorf("type %s didn't match expected type %s", v.referenceType, referenceType))
	}
	return v.reference
}

func createAsReferenceAndAddToHeap(referenceType string, reference interface{}, f frame) variable {
	*f.heap = append(*f.heap, reference)
	return variable{
		referenceType: referenceType,
		reference:     &((*f.heap)[len(*f.heap)-1]),
	}
}

type state struct {
	frames []frame
	files  []parse.ClassFile
}

type frame struct {
	codeReader    *parse.ClassFileReader
	operandStack  *varSlice
	localVariable *varSlice
	file          parse.ClassFile
	heap          *[]interface{}
}

type varSlice []variable

func (s *varSlice) pop() variable {
	x := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	return x
}

type class struct {
	//methods map[string]func() error
	isVirtual bool
	name      string
	vars      map[string]variable
}

func execute(file parse.ClassFile) error {
	var mainMethod parse.MethodInfo
	for _, method := range file.Methods {
		if file.ConstantPool[method.NameIndex-1].(parse.ConstantUtf8Info).Text == "main" {
			mainMethod = method
			goto mainFound
		}
	}
	return fmt.Errorf("no main method found")
mainFound:
	descriptor := file.ConstantPool[mainMethod.DescriptorIndex-1].(parse.ConstantUtf8Info).Text
	if !(descriptor == "([Ljava/lang/String;)V" && mainMethod.AccessFlags == methodFlagAccPublic|methodFlagAccStatic) {
		return fmt.Errorf("main method not formated corectly")
	}
	var mainMethodCodeAttribute parse.AttributeInfo
	for _, attribute := range mainMethod.Attributes {
		if file.ConstantPool[attribute.AttributeNameIndex-1].(parse.ConstantUtf8Info).Text == "Code" {
			mainMethodCodeAttribute = attribute
			goto codeFound
		}
	}
	return fmt.Errorf("code attribute not found in main")
codeFound:
	reader := (*parse.ClassFileReader)(bytes.NewReader(mainMethodCodeAttribute.Info))

	maxStack, err := reader.ReadU2() // maxStack
	if err != nil {
		return err
	}
	_ = maxStack
	maxLocals, err := reader.ReadU2() // maxLocals
	if err != nil {
		return err
	}

	codeLength, err := reader.ReadU4()
	if err != nil {
		return err
	}
	code := make([]byte, codeLength)
	_, err = reader.Read(code)
	if err != nil {
		return err
	}

	exceptionTableLength, _ := reader.ReadU2()
	exceptionTable := make([]byte, exceptionTableLength)
	_, err = reader.Read(exceptionTable) //TODO: Parse the exception table
	if err != nil {
		return err
	}

	_, _, err = reader.ReadAttributes() // Attributes
	if err != nil {
		return err
	}

	// ------------------- Code Execution ---------------------
	operandStack := make(varSlice, 0)
	localVariable := make(varSlice, maxLocals)
	heap := make([]interface{}, 0)
	f := frame{
		codeReader:    (*parse.ClassFileReader)(bytes.NewReader(code)),
		operandStack:  &operandStack,
		localVariable: &localVariable,
		file:          file,
		heap:          &heap,
	}
	s := state{
		frames: make([]frame, 0),
		files:  make([]parse.ClassFile, 0),
	}
	s.frames = append(s.frames, f)
	s.files = append(s.files, file)

	for f.codeReader.Len() > 0 {
		b, err := f.codeReader.ReadByte()
		if err != nil {
			return err
		}

		inst, err := getInstruction(b)
		if err != nil {
			return err
		}

		err = inst(&s, f)
		if err != nil {
			return err
		}
	}

	return nil
}

// runMethod finds a method and runs it in a new frame
//
// method is formated as *methodClass*.*methodName*
func runMethod(method string, methodDescriptor string, s *state, args []variable) error { // TODO: cashing
	methodSplit := strings.SplitN(method, ".", 2)
	methodClass, methodName := methodSplit[0], methodSplit[1]

	var mainMethod parse.MethodInfo
	var file parse.ClassFile
	for _, classFile := range s.files {
		if (*classFile.ConstantPool[classFile.ThisClass-1].(parse.ConstantClassInfo).Name).(parse.ConstantUtf8Info).Text != methodClass {
			continue
		}
		for _, method := range classFile.Methods {
			if classFile.ConstantPool[method.NameIndex-1].(parse.ConstantUtf8Info).Text == methodName &&
				classFile.ConstantPool[method.DescriptorIndex-1].(parse.ConstantUtf8Info).Text == methodDescriptor {
				mainMethod = method
				file = classFile
				goto methodFound
			}
		}
		break
	}
	return fmt.Errorf(`method "%s" not found`, methodName)
methodFound:
	descriptor := file.ConstantPool[mainMethod.DescriptorIndex-1].(parse.ConstantUtf8Info).Text
	if !(descriptor == methodDescriptor) {
		return fmt.Errorf("method not formated as expected corectly: %s != %s", descriptor, methodDescriptor)
	}

	if mainMethod.AccessFlags&methodFlagAccNative != 0 {
		if methodName == "registerNatives" {
			// TODO: implement native methods correctly
			return nil
		} else {
			return fmt.Errorf("method %s is native", methodName)
		}
	}

	var mainMethodCodeAttribute parse.AttributeInfo
	for _, attribute := range mainMethod.Attributes {
		if file.ConstantPool[attribute.AttributeNameIndex-1].(parse.ConstantUtf8Info).Text == "Code" {
			mainMethodCodeAttribute = attribute
			goto codeFound
		}
	}
	return fmt.Errorf("code attribute not found in main")
codeFound:
	reader := (*parse.ClassFileReader)(bytes.NewReader(mainMethodCodeAttribute.Info))

	maxStack, err := reader.ReadU2() // maxStack
	if err != nil {
		return err
	}
	_ = maxStack
	maxLocals, err := reader.ReadU2() // maxLocals
	if err != nil {
		return err
	}
	_ = maxLocals

	codeLength, err := reader.ReadU4()
	if err != nil {
		return err
	}
	code := make([]byte, codeLength)
	_, err = reader.Read(code)
	if err != nil {
		return err
	}

	exceptionTableLength, _ := reader.ReadU2()             // FIXME: fuck I'm reading exception tables wrong...
	exceptionTable := make([]byte, exceptionTableLength*8) // each entry is 8 bytes long so im just reading them away here... hope it works :)
	_, err = reader.Read(exceptionTable)                   //TODO: Parse the exception table
	if err != nil {
		return err
	}

	_, _, err = reader.ReadAttributes() // TODO: Parse these while creating the class file to get quicker :)
	if err != nil {
		return err
	}

	// ------------------- Code Execution ---------------------
	operandStack := make(varSlice, 0)
	localVariable := make(varSlice, maxLocals)
	heap := make([]interface{}, 0)
	for i, arg := range args {
		localVariable[i] = arg
	}
	f := frame{
		codeReader:    (*parse.ClassFileReader)(bytes.NewReader(code)),
		operandStack:  &operandStack,
		localVariable: &localVariable,
		file:          file,
		heap:          &heap,
	}
	s.frames = append(s.frames, f)

	for f.codeReader.Len() > 0 {
		b, err := f.codeReader.ReadByte()
		if err != nil {
			return err
		}

		inst, err := getInstruction(b)
		if err != nil {
			return err
		}

		err = inst(s, f)
		if err != nil {
			return err
		}
	}

	return nil
}
