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
	frames []frame
	file   parse.ClassFile
}

type frame struct {
	codeReader    *parse.ClassFileReader
	operandStack  *[]interface{} // FIXME: there are also different types of variables xD
	localVariable *[]interface{} // FIXME: see above ;)
	file          parse.ClassFile
}

var instructionSet map[byte]func(s *state, f frame) error

func init() {
	instructionSet = map[byte]func(s *state, f frame) error{
		4: func(s *state, f frame) error { // iconst_1
			*f.operandStack = append(*f.operandStack, 1)
			return nil
		},
		5: func(s *state, f frame) error { // iconst_2
			*f.operandStack = append(*f.operandStack, 2)
			return nil
		},
		6: func(s *state, f frame) error { // iconst_3
			*f.operandStack = append(*f.operandStack, 3)
			return nil
		},
		7: func(s *state, f frame) error { // iconst_4
			*f.operandStack = append(*f.operandStack, 4)
			return nil
		},
		18: func(s *state, f frame) error { // ldc
			idx, err := f.codeReader.ReadByte()
			if err != nil {
				return err
			}

			if utf8, ok := s.file.ConstantPool[idx].(parse.ConstantUtf8Info); ok {
				*f.operandStack = append(*f.operandStack, utf8.Text)
				return nil
			}
			return fmt.Errorf("ldc (18) only implemented for string")
		},
		26: func(s *state, f frame) error { // iload_0
			*f.operandStack = append(*f.operandStack, (*f.localVariable)[0]) // TODO: catch errors
			return nil
		},
		27: func(s *state, f frame) error { // iload_1
			*f.operandStack = append(*f.operandStack, (*f.localVariable)[1]) // TODO: catch errors
			return nil
		},
		28: func(s *state, f frame) error { // iload_2
			*f.operandStack = append(*f.operandStack, (*f.localVariable)[2]) // TODO: catch errors
			return nil
		},
		29: func(s *state, f frame) error { // iload_3
			*f.operandStack = append(*f.operandStack, (*f.localVariable)[3]) // TODO: catch errors
			return nil
		},
		59: func(s *state, f frame) error { // istore_0
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			(*f.localVariable)[0] = lastVal
			return nil
		},
		60: func(s *state, f frame) error { // istore_1
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			(*f.localVariable)[1] = lastVal
			return nil
		},
		61: func(s *state, f frame) error { // istore_2
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			(*f.localVariable)[2] = lastVal
			return nil
		},
		62: func(s *state, f frame) error { // istore_3
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			(*f.localVariable)[3] = lastVal
			return nil
		},
		104: func(s *state, f frame) error { // imul
			lastVal := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
			lastVal2 := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			*f.operandStack = append(*f.operandStack, lastVal*lastVal2)
			return nil
		},
		172: func(s *state, f frame) error { // ireturn
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]

			*s.frames[len(s.frames)-2].operandStack = append(*s.frames[len(s.frames)-2].operandStack, lastVal)

			s.frames = s.frames[:len(s.frames)-1]
			return nil // TODO: check errors
		},
		177: func(s *state, f frame) error { // return
			// TODO: return void and empty operandStack
			return nil
		},
		178: func(s *state, f frame) error { // getstatic
			// TODO: get a static field
			_, err := f.codeReader.ReadByte()
			_, err = f.codeReader.ReadByte()
			return err
		},
		182: func(s *state, f frame) error { // invokevirtual
			// TODO: Invoke instance method; dispatch based on class
			// in my test case this will always be System.out.println(int)
			address, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}

			method := f.file.ConstantPool[address-1].(parse.ConstantMethodrefInfo)
			methodNameAndType := (*method.NameAndType).(parse.ConstantNameAndTypeInfo)
			name := (*methodNameAndType.Name).(parse.ConstantUtf8Info).Text
			if name != "println" {
				return fmt.Errorf(`function "%s" not implemented`, name)
			}
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
			if out, ok := lastVal.(int); ok {
				println(out)
				return nil
			}
			if out, ok := lastVal.(string); ok {
				println(out)
				return nil
			}

			return fmt.Errorf("unimplementet type for println")
		},
		184: func(s *state, f frame) error { // invokestatic
			address, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}

			method := f.file.ConstantPool[address-1].(parse.ConstantMethodrefInfo)

			methodClass := (*method.Class).(parse.ConstantClassInfo)
			className := (*methodClass.Name).(parse.ConstantUtf8Info).Text
			currentClass := s.file.ConstantPool[s.file.ThisClass-1].(parse.ConstantClassInfo)
			currentClassName := (*currentClass.Name).(parse.ConstantUtf8Info).Text
			if className != currentClassName {
				return fmt.Errorf("support for different classes not implemented")
			}

			methodNameAndType := (*method.NameAndType).(parse.ConstantNameAndTypeInfo)
			name := (*methodNameAndType.Name).(parse.ConstantUtf8Info).Text
			descriptor := (*methodNameAndType.Descriptor).(parse.ConstantUtf8Info).Text

			if descriptor != "(I)I" {
				return fmt.Errorf("not supported descriptor")
			}
			argCount := 1 // TODO: this only works for the example
			_ = argCount
			args := make([]interface{}, 0)
			lastVal := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
			args = append(args, lastVal)

			err = runMethod(name, descriptor, s, args)
			if err != nil {
				return err
			}

			return err
		},
	}
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
	operandStack := make([]interface{}, 0)
	localVariable := make([]interface{}, maxLocals)
	f := frame{
		codeReader:    (*parse.ClassFileReader)(bytes.NewReader(code)),
		operandStack:  &operandStack,
		localVariable: &localVariable,
		file:          file,
	}
	s := state{
		frames: make([]frame, 0),
		file:   file,
	}
	s.frames = append(s.frames, f)

	for f.codeReader.Len() > 0 {
		b, err := f.codeReader.ReadByte()
		if err != nil {
			return err
		}

		if _, ok := instructionSet[b]; !ok {
			return fmt.Errorf("instruction \"%d\" not implemented", b)
		}

		err = instructionSet[b](&s, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func runMethod(methodName string, methodDescriptor string, s *state, args []interface{}) error {
	var mainMethod parse.MethodInfo
	for _, method := range s.file.Methods {
		if s.file.ConstantPool[method.NameIndex-1].(parse.ConstantUtf8Info).Text == methodName {
			mainMethod = method
			goto methodFound
		}
	}
	return fmt.Errorf(`method "%s" not found`, methodName)
methodFound:
	descriptor := s.file.ConstantPool[mainMethod.DescriptorIndex-1].(parse.ConstantUtf8Info).Text
	if !(descriptor == methodDescriptor) {
		return fmt.Errorf("main method not formated corectly")
	}
	var mainMethodCodeAttribute parse.AttributeInfo
	for _, attribute := range mainMethod.Attributes {
		if s.file.ConstantPool[attribute.AttributeNameIndex-1].(parse.ConstantUtf8Info).Text == "Code" {
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
	operandStack := make([]interface{}, 0)
	localVariable := args
	f := frame{
		codeReader:    (*parse.ClassFileReader)(bytes.NewReader(code)),
		operandStack:  &operandStack,
		localVariable: &localVariable,
		file:          s.file,
	}
	s.frames = append(s.frames, f)

	for f.codeReader.Len() > 0 {
		b, err := f.codeReader.ReadByte()
		if err != nil {
			return err
		}

		if _, ok := instructionSet[b]; !ok {
			return fmt.Errorf("instruction \"%d\" not implemented", b)
		}

		err = instructionSet[b](s, f)
		if err != nil {
			return err
		}
	}

	return nil
}
