package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PPTide/gojdk/parse"
	"io"
	"strconv"
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

func getInstruction(instruction byte) (func(s *state, f frame) error, error) {
	switch instruction {
	case 2:
		return func(s *state, f frame) error { // iconst_m1
			*f.operandStack = append(*f.operandStack, -1)
			return nil
		}, nil
	case 3:
		return func(s *state, f frame) error { // iconst_0
			*f.operandStack = append(*f.operandStack, 0)
			return nil
		}, nil
	case 4:
		return func(s *state, f frame) error { // iconst_1
			*f.operandStack = append(*f.operandStack, 1)
			return nil
		}, nil
	case 5:
		return func(s *state, f frame) error { // iconst_2
			*f.operandStack = append(*f.operandStack, 2)
			return nil
		}, nil
	case 6:
		return func(s *state, f frame) error { // iconst_3
			*f.operandStack = append(*f.operandStack, 3)
			return nil
		}, nil
	case 7:
		return func(s *state, f frame) error { // iconst_4
			*f.operandStack = append(*f.operandStack, 4)
			return nil
		}, nil
	case 16:
		return func(s *state, f frame) error { // bipush
			b, err := f.codeReader.ReadByte()
			if err != nil {
				return err
			}

			*f.operandStack = append(*f.operandStack, int(b))
			return nil
		}, nil
	case 18:
		return func(s *state, f frame) error { // ldc
			idx, err := f.codeReader.ReadByte()
			if err != nil {
				return err
			}

			if utf8, ok := s.file.ConstantPool[idx].(parse.ConstantUtf8Info); ok {
				*f.operandStack = append(*f.operandStack, utf8.Text)
				return nil
			}
			return fmt.Errorf("ldc (18) only implemented for string")
		}, nil
	case 26:
		return func(s *state, f frame) error { // iload_0
			*f.operandStack = append(*f.operandStack, (*f.localVariable)[0]) // TODO: catch errors
			return nil
		}, nil
	case 27:
		return func(s *state, f frame) error { // iload_1
			*f.operandStack = append(*f.operandStack, (*f.localVariable)[1]) // TODO: catch errors
			return nil
		}, nil
	case 28:
		return func(s *state, f frame) error { // iload_2
			*f.operandStack = append(*f.operandStack, (*f.localVariable)[2]) // TODO: catch errors
			return nil
		}, nil
	case 29:
		return func(s *state, f frame) error { // iload_3
			*f.operandStack = append(*f.operandStack, (*f.localVariable)[3]) // TODO: catch errors
			return nil
		}, nil
	case 59:
		return func(s *state, f frame) error { // istore_0
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			(*f.localVariable)[0] = lastVal
			return nil
		}, nil
	case 60:
		return func(s *state, f frame) error { // istore_1
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			(*f.localVariable)[1] = lastVal
			return nil
		}, nil
	case 61:
		return func(s *state, f frame) error { // istore_2
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			(*f.localVariable)[2] = lastVal
			return nil
		}, nil
	case 62:
		return func(s *state, f frame) error { // istore_3
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			(*f.localVariable)[3] = lastVal
			return nil
		}, nil
	case 96:
		return func(s *state, f frame) error { // iadd
			value2 := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
			value1 := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			*f.operandStack = append(*f.operandStack, value1+value2)
			return nil
		}, nil
	case 100:
		return func(s *state, f frame) error { // isub
			value2 := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
			value1 := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			*f.operandStack = append(*f.operandStack, value1-value2)
			return nil
		}, nil
	case 104:
		return func(s *state, f frame) error { // imul
			lastVal := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
			lastVal2 := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			*f.operandStack = append(*f.operandStack, lastVal*lastVal2)
			return nil
		}, nil
	case 132:
		return func(s *state, f frame) error { // iinc
			index, err := f.codeReader.ReadByte()
			if err != nil {
				return err
			}
			constVal, err := f.codeReader.ReadByte()
			if err != nil {
				return err
			}

			(*f.localVariable)[index] = (*f.localVariable)[index].(int) + int(constVal)
			return nil
		}, nil
	case 162:
		return func(s *state, f frame) error { // if_icmpge <- >=
			value2 := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
			value1 := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			branchOffset, err := f.codeReader.ReadU2() //FIXME: might be negative
			if err != nil {
				return err
			}
			branchOffset = branchOffset - 3 // get the offset without the branch bytes

			if value1 >= value2 {
				_, err := f.codeReader.Seek(int64(branchOffset), io.SeekCurrent)
				if err != nil {
					return err
				}
			}
			return nil
		}, nil
	case 163:
		return func(s *state, f frame) error { // if_icmpgt <- >
			value2 := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
			value1 := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			branchOffset, err := f.codeReader.ReadU2() //FIXME: might be negative
			if err != nil {
				return err
			}
			branchOffset = branchOffset - 3 // get the offset without the branch bytes

			if value1 > value2 {
				_, err := f.codeReader.Seek(int64(branchOffset), io.SeekCurrent)
				if err != nil {
					return err
				}
			}
			return nil
		}, nil
	case 167:
		return func(s *state, f frame) error {
			branchOffsetTmp, err := f.codeReader.ReadU2() //FIXME: might be negative
			if err != nil {
				return err
			}

			branchOffset := int16(branchOffsetTmp)

			branchOffset = branchOffset - 3 // get the offset without the branch bytes

			_, err = f.codeReader.Seek(int64(branchOffset), io.SeekCurrent)
			return err
		}, nil
	case 172:
		return func(s *state, f frame) error { // ireturn
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]

			*s.frames[len(s.frames)-2].operandStack = append(*s.frames[len(s.frames)-2].operandStack, lastVal)

			s.frames = s.frames[:len(s.frames)-1]

			_, err := f.codeReader.Seek(0, io.SeekEnd)
			if err != nil {
				return err
			}

			return nil // TODO: check errors
		}, nil
	case 177:
		return func(s *state, f frame) error { // return
			// TODO: return void and empty operandStack

			_, err := f.codeReader.Seek(0, io.SeekEnd)
			if err != nil {
				return err
			}

			return nil
		}, nil
	case 178:
		return func(s *state, f frame) error { // getstatic
			// TODO: get a static field
			_, err := f.codeReader.ReadByte()
			_, err = f.codeReader.ReadByte()
			return err
		}, nil
	case 182:
		return func(s *state, f frame) error { // invokevirtual
			// TODO: Invoke instance method; dispatch based on class
			// in my test case this will always be System.out.println(int)
			address, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}

			method := f.file.ConstantPool[address-1].(parse.ConstantMethodrefInfo)
			methodNameAndType := (*method.NameAndType).(parse.ConstantNameAndTypeInfo)
			name := (*methodNameAndType.Name).(parse.ConstantUtf8Info).Text
			switch name { // TODO: find a nicer place to put this
			case "println":
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
			case "print":
				lastVal := (*f.operandStack)[len(*f.operandStack)-1]
				*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
				if out, ok := lastVal.(int); ok {
					print(out)
					return nil
				}
				if out, ok := lastVal.(string); ok {
					print(out)
					return nil
				}

				return fmt.Errorf("unimplementet type for print")
			default:
				return fmt.Errorf(`function "%s" not implemented`, name)
			}
		}, nil
	case 184:
		return func(s *state, f frame) error { // invokestatic
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
		}, nil
	case 186:
		return func(s *state, f frame) error { // invokedynamic
			address, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}
			check, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}
			if check != 0 {
				return errors.New("unexpected value in execution of invokedynamic")
			}

			DynamicInfo := f.file.ConstantPool[address-1].(parse.ConstantInvokeDynamicInfo)

			var bootstrapMethodAttribute parse.AttributeInfo

			for _, attribute := range f.file.Attributes { //TODO: extract BootstrapMethods in extra function
				attributeName := f.file.ConstantPool[attribute.AttributeNameIndex-1].(parse.ConstantUtf8Info).Text

				if attributeName == "BootstrapMethods" {
					bootstrapMethodAttribute = attribute
					goto foundAttribute
				}
			}
			return errors.New("BootstrapMethod attribute not found")

		foundAttribute:
			r := (*parse.ClassFileReader)(bytes.NewReader(bootstrapMethodAttribute.Info))
			numBootstrapMethods, err := r.ReadU2()
			if err != nil {
				return err
			}

			bootstrapMethods := make([]struct {
				bootstrapMethodRef    int
				numBootstrapArguments int
				bootstrapArguments    []int // index to constant pool
			}, 0)

			for i := 0; i < numBootstrapMethods; i++ {
				x := struct {
					bootstrapMethodRef    int
					numBootstrapArguments int
					bootstrapArguments    []int
				}{}
				x.bootstrapMethodRef, err = r.ReadU2()
				if err != nil {
					return err
				}
				x.numBootstrapArguments, err = r.ReadU2()
				if err != nil {
					return err
				}

				for j := 0; j < x.numBootstrapArguments; j++ {
					argument, err := r.ReadU2()
					if err != nil {
						return err
					}
					x.bootstrapArguments = append(x.bootstrapArguments, argument)
				}

				bootstrapMethods = append(bootstrapMethods, x)
			}
			// ------------ End Search for BootstrapMethods ------------

			bootstrapMethod := bootstrapMethods[DynamicInfo.BootstrapMethodAttrIndex]

			if (*(*DynamicInfo.NameAndType).(parse.ConstantNameAndTypeInfo).Name).(parse.ConstantUtf8Info).Text != "makeConcatWithConstants" {
				return fmt.Errorf(`unknown Dynamic funtion "%s"`, (*(*DynamicInfo.NameAndType).(parse.ConstantNameAndTypeInfo).Name).(parse.ConstantUtf8Info).Text)
			}

			// FIXME: expects only makeConcatWithConstants
			lastVal := strconv.Itoa((*f.operandStack)[len(*f.operandStack)-1].(int))
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			string2Index := bootstrapMethod.bootstrapArguments[0]
			string2 := (*f.file.ConstantPool[string2Index-1].(parse.ConstantStringInfo).String).(parse.ConstantUtf8Info).Text

			*f.operandStack = append(*f.operandStack, lastVal+string2)

			return nil
		}, nil
	}
	return nil, fmt.Errorf(`unknown instruction "%s"`, instruction)
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

func runMethod(methodName string, methodDescriptor string, s *state, args []interface{}) error { // TODO: cashing
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

	_, _, err = reader.ReadAttributes() // TODO: Parse these while creating the class file to get quicker :)
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
