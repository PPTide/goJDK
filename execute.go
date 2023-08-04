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

type state struct {
	codeReader   *bytes.Reader
	operandStack *[]int // FIXME: there are also different types of variables xD
}

var instructionSet = map[byte]func(s state) error{
	4: func(s state) error { // iconst_1
		*s.operandStack = append(*s.operandStack, 1)
		return nil
	},
	5: func(s state) error { // iconst_2
		*s.operandStack = append(*s.operandStack, 2)
		return nil
	},
	6: func(s state) error { // iconst_3
		*s.operandStack = append(*s.operandStack, 3)
		return nil
	},
	7: func(s state) error { // iconst_4
		*s.operandStack = append(*s.operandStack, 4)
		return nil
	},
	177: func(s state) error {
		// TODO: return void and empty operandStack
		return nil
	},
	178: func(s state) error { // getstatic
		// TODO: get a static field
		_, err := s.codeReader.ReadByte()
		_, err = s.codeReader.ReadByte()
		return err
	},
	182: func(s state) error { // invokevirtual
		// TODO: Invoke instance method; dispatch based on class
		// in my test case this will always be System.out.println(int)
		_, err := s.codeReader.ReadByte()
		_, err = s.codeReader.ReadByte()
		lastVal := (*s.operandStack)[len(*s.operandStack)-1]
		*s.operandStack = (*s.operandStack)[:len(*s.operandStack)-1]
		println(lastVal)

		return err
	},
	184: func(s state) error { // invokestatic
		// TODO: Invoke a class (static) method
		// in my test case this will be square?
		_, err := s.codeReader.ReadByte()
		_, err = s.codeReader.ReadByte()
		lastVal := (*s.operandStack)[len(*s.operandStack)-1]
		*s.operandStack = (*s.operandStack)[:len(*s.operandStack)-1]

		*s.operandStack = append(*s.operandStack, lastVal*lastVal)
		return err
	},
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
	_, err = reader.ReadU2() // maxLocals
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
	operandStack := make([]int, 0)
	s := state{
		codeReader:   bytes.NewReader(code),
		operandStack: &operandStack,
	}

	for s.codeReader.Len() > 0 {
		b, err := s.codeReader.ReadByte()
		if err != nil {
			return err
		}

		if _, ok := instructionSet[b]; !ok {
			return fmt.Errorf("instruction \"%d\" not implemented", b)
		}

		instructionSet[b](s)
	}

	return nil
}
