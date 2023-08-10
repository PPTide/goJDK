package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PPTide/gojdk/parse"
	"io"
	"reflect"
	"strconv"
)

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

			switch t := f.file.ConstantPool[idx-1].(type) {
			case parse.ConstantStringInfo:
				*f.operandStack = append(*f.operandStack, (*t.String).(parse.ConstantUtf8Info).Text)
				return nil
			default:
				_ = t
				return fmt.Errorf("ldc (18) not implemented for %s", reflect.TypeOf(f.file.ConstantPool[idx-1]))
			}
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
	case 89: // dup
		return func(s *state, f frame) error {
			*f.operandStack = append(*f.operandStack, (*f.operandStack)[len(*f.operandStack)-1])
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
	case 182: // invokevirtual
		return func(s *state, f frame) error {
			// TODO: Invoke instance method; dispatch based on class
			// in my test case this will always be System.out.println(int)
			address, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}

			method := f.file.ConstantPool[address-1].(parse.ConstantMethodrefInfo)
			methodNameAndType := (*method.NameAndType).(parse.ConstantNameAndTypeInfo)
			name := (*methodNameAndType.Name).(parse.ConstantUtf8Info).Text
			descriptor := (*methodNameAndType.Descriptor).(parse.ConstantUtf8Info).Text
			className := (*(*method.Class).(parse.ConstantClassInfo).Name).(parse.ConstantUtf8Info).Text
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
				_ = (*f.operandStack)[:len(*f.operandStack)-1] // objectref
				*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

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
				_ = (*f.operandStack)[:len(*f.operandStack)-1] // objectref
				*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

				return fmt.Errorf("unimplementet type for print")
			case "append":
				if className != "java/lang/StringBuilder" {
					return fmt.Errorf("unknow function append in class " + className)
				}
				if !((descriptor == "(Ljava/lang/String;)Ljava/lang/StringBuilder;") ||
					(descriptor == "(I)Ljava/lang/StringBuilder;")) {
					return fmt.Errorf(`unexpected descriptor "%s" for append`, (*methodNameAndType.Descriptor).(parse.ConstantUtf8Info).Text)
				}

				arg1 := (*f.operandStack)[len(*f.operandStack)-1]
				*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
				objectRef := (*f.operandStack)[len(*f.operandStack)-1]
				*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

				_, _ = arg1, objectRef // TODO: add more type checking

				toAppend, err := toString(arg1)
				if err != nil {
					return err
				}

				(objectRef.(class)).vars["string"] = (objectRef.(class)).vars["string"].(string) + toAppend

				*f.operandStack = append(*f.operandStack, objectRef)

				return nil
			case "toString":
				if className != "java/lang/StringBuilder" {
					return fmt.Errorf("unknow function toString in class " + className)
				}
				if !(descriptor == "()Ljava/lang/String;") {
					return fmt.Errorf(`unexpected descriptor "%s" for toString`, (*methodNameAndType.Descriptor).(parse.ConstantUtf8Info).Text)
				}

				objectRef := (*f.operandStack)[len(*f.operandStack)-1]
				*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

				*f.operandStack = append(*f.operandStack, objectRef.(class).vars["string"].(string))

				return nil
			default:
				return fmt.Errorf(`function "%s" not implemented (invokevirtual)`, name)
			}
		}, nil
	case 183: // invokespecial
		return func(s *state, f frame) error {
			// TODO: Invoke instance method; special handling for superclass, private, and instance initialization method invocations
			address, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}

			method := f.file.ConstantPool[address-1].(parse.ConstantMethodrefInfo)
			methodNameAndType := (*method.NameAndType).(parse.ConstantNameAndTypeInfo)
			name := (*methodNameAndType.Name).(parse.ConstantUtf8Info).Text
			switch name { // TODO: find a nicer place to put this
			case "<init>":
				return nil // FIXME: THIS WILL NEVER WORK
			default:
				return fmt.Errorf(`function "%s" not implemented (invokespecial)`, name)
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
			currentClass := f.file.ConstantPool[f.file.ThisClass-1].(parse.ConstantClassInfo)
			currentClassName := (*currentClass.Name).(parse.ConstantUtf8Info).Text
			if className != currentClassName {
				//return fmt.Errorf("support for different classes not implemented")
				for _, file := range s.files {
					fileClass := file.ConstantPool[f.file.ThisClass-1].(parse.ConstantClassInfo)
					fileClassName := (*fileClass.Name).(parse.ConstantUtf8Info).Text

					if fileClassName == className {
						goto fileFound
					}
				}
				// File is not loaded yet
				file, err := parse.Parse(className + ".class")
				if err != nil {
					return err
				}
				s.files = append(s.files, file)
				goto fileFound
			}

		fileFound:
			methodNameAndType := (*method.NameAndType).(parse.ConstantNameAndTypeInfo)
			name := (*methodNameAndType.Name).(parse.ConstantUtf8Info).Text
			descriptor := (*methodNameAndType.Descriptor).(parse.ConstantUtf8Info).Text

			if descriptor != "(I)I" {
				return fmt.Errorf("not supported descriptor")
			}
			argCount := 1 // TODO: this only works for the example
			_ = argCount
			args := make([]variable, 0)
			lastVal := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
			args = append(args, lastVal)

			err = runMethod(name, descriptor, s, args) //FIXME: This won't work for methods with the same name in different classes
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
			// FIXME: no worky :(
			lastVal := strconv.Itoa((*f.operandStack)[len(*f.operandStack)-1].(int))
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			string2Index := bootstrapMethod.bootstrapArguments[0]
			string2 := (*f.file.ConstantPool[string2Index-1].(parse.ConstantStringInfo).String).(parse.ConstantUtf8Info).Text

			*f.operandStack = append(*f.operandStack, string2+lastVal)

			return nil
		}, nil
	case 187: // new
		return func(s *state, f frame) error {
			address, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}
			name := (*f.file.ConstantPool[address-1].(parse.ConstantClassInfo).Name).(parse.ConstantUtf8Info).Text
			if name == "java/lang/StringBuilder" {
				classInst := class{
					isVirtual: true,
					name:      name,
					vars:      make(map[string]variable),
				}
				classInst.vars["initialCapacity"] = 16
				classInst.vars["string"] = ""
				*f.operandStack = append(*f.operandStack, classInst)
				return nil
			}
			return fmt.Errorf(`_new_ not implemented for "%s"`, name)
		}, nil
	}
	return nil, fmt.Errorf(`unknown instruction "%d"`, instruction)
}

func toString(from variable) (string, error) {
	to := ""
	switch t := from.(type) {
	case string:
		to = t
	case int:
		to = strconv.Itoa(t)
	default:
		return "", fmt.Errorf("type %s not implemented for toString", reflect.TypeOf(t))
	}
	return to, nil
}
