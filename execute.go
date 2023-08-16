package main

import (
	"bytes"
	"fmt"
	"github.com/PPTide/gojdk/parse"
)

const (
	methodFlagAccPublic = 1
	methodFlagAccStatic = 8
)

type variable interface{}

type state struct {
	frames []frame
	files  []parse.ClassFile
}

type frame struct {
	codeReader    *parse.ClassFileReader
	operandStack  *[]variable // FIXME: there are also different types of variables xD
	localVariable *[]variable // FIXME: see above ;)
	file          parse.ClassFile
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

	_, err := reader.ReadU2() // maxStack
	if err != nil {
		return err
	}
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
	operandStack := make([]variable, 0)
	localVariable := make([]variable, maxLocals)
	f := frame{
		codeReader:    (*parse.ClassFileReader)(bytes.NewReader(code)),
		operandStack:  &operandStack,
		localVariable: &localVariable,
		file:          file,
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

func runMethod(methodName string, methodDescriptor string, s *state, args []variable) error { // TODO: cashing
	var mainMethod parse.MethodInfo
	var file parse.ClassFile
	for _, classFile := range s.files {
		for _, method := range classFile.Methods {
			if classFile.ConstantPool[method.NameIndex-1].(parse.ConstantUtf8Info).Text == methodName {
				mainMethod = method
				file = classFile
				goto methodFound
			}
		}

	}
	return fmt.Errorf(`method "%s" not found`, methodName)
methodFound:
	descriptor := file.ConstantPool[mainMethod.DescriptorIndex-1].(parse.ConstantUtf8Info).Text
	if !(descriptor == methodDescriptor) {
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

	_, err := reader.ReadU2() // maxStack
	if err != nil {
		return err
	}
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

	_, _, err = reader.ReadAttributes() // TODO: Parse these while creating the class file to get quicker :)
	if err != nil {
		return err
	}

	// ------------------- Code Execution ---------------------
	operandStack := make([]variable, 0)
	localVariable := make([]variable, maxLocals)
	for i, arg := range args {
		localVariable[i] = arg
	}
	f := frame{
		codeReader:    (*parse.ClassFileReader)(bytes.NewReader(code)),
		operandStack:  &operandStack,
		localVariable: &localVariable,
		file:          file,
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
