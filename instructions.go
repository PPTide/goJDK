package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PPTide/gojdk/parse"
	"io"
	"reflect"
	"strconv"
	"strings"
)

type field struct {
	fieldRefInfo parse.ConstantFieldrefInfo
	value        interface{}
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

			switch t := f.file.ConstantPool[idx-1].(type) {
			case parse.ConstantStringInfo:
				*f.operandStack = append(*f.operandStack, (*t.String).(parse.ConstantUtf8Info).Text)
				return nil
			default:
				_ = t
				return fmt.Errorf("ldc (18) not implemented for %s", reflect.TypeOf(f.file.ConstantPool[idx-1]))
			}
		}, nil
	case 25: // aload
		return func(s *state, f frame) error {
			b, err := f.codeReader.ReadByte()
			if err != nil {
				return err
			}

			*f.operandStack = append(*f.operandStack, (*f.localVariable)[b])
			return nil
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
	case 42: // aload_0
		return func(s *state, f frame) error {
			*f.operandStack = append(*f.operandStack, (*f.localVariable)[0])
			return nil
		}, nil
	case 43: // aload_1
		return func(s *state, f frame) error {
			*f.operandStack = append(*f.operandStack, (*f.localVariable)[1])
			return nil
		}, nil
	case 44: // aload_2
		return func(s *state, f frame) error {
			*f.operandStack = append(*f.operandStack, (*f.localVariable)[2])
			return nil
		}, nil
	case 45: // aload_3
		return func(s *state, f frame) error {
			*f.operandStack = append(*f.operandStack, (*f.localVariable)[3])
			return nil
		}, nil
	case 58: // astore
		// TODO: type checking
		return func(s *state, f frame) error {
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			b, err := f.codeReader.ReadByte()
			if err != nil {
				return err
			}

			(*f.localVariable)[b] = lastVal
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
	case 75: // astore_0
		// TODO: type checking
		return func(s *state, f frame) error {
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			(*f.localVariable)[0] = lastVal
			return nil
		}, nil
	case 76: // astore_1
		// TODO: type checking
		return func(s *state, f frame) error {
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			(*f.localVariable)[1] = lastVal
			return nil
		}, nil
	case 77: // astore_2
		// TODO: type checking
		return func(s *state, f frame) error {
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			(*f.localVariable)[2] = lastVal
			return nil
		}, nil
	case 78: // astore_3
		// TODO: type checking
		return func(s *state, f frame) error {
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			(*f.localVariable)[3] = lastVal
			return nil
		}, nil
	case 87: // pop
		return func(s *state, f frame) error {
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
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
	case 112: // irem
		return func(s *state, f frame) error {
			value2 := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
			value1 := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			*f.operandStack = append(*f.operandStack, value1-(value1/value2)*value2)
			return nil
		}, nil
	case 122: // ishr
		return func(s *state, f frame) error {
			value2 := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
			value1 := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			*f.operandStack = append(*f.operandStack, value1>>value2)
			return nil
		}, nil
	case 132: // iinc
		return func(s *state, f frame) error {
			index, err := f.codeReader.ReadByte()
			if err != nil {
				return err
			}
			constVal, err := f.codeReader.ReadByte()
			if err != nil {
				return err
			}

			(*f.localVariable)[index] = (*f.localVariable)[index].(int) + int(int8(constVal))
			return nil
		}, nil
	case 158: // ifle -> <=
		return func(s *state, f frame) error {
			lastVal := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			branchOffset, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}
			branchOffset = branchOffset - 3 // get the offset without the branch bytes

			if lastVal <= 0 {
				_, err := f.codeReader.Seek(int64(branchOffset), io.SeekCurrent)
				if err != nil {
					return err
				}
			}
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
	case 180: // getfield
		return func(s *state, f frame) error {
			// FIXME: This is not implemented
			address, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}

			fieldref := f.file.ConstantPool[address-1].(parse.ConstantFieldrefInfo)

			x := field{
				fieldRefInfo: fieldref,
				value:        []rune{}, // FIXME: i am making this work for exactly one case xD
			}

			*f.operandStack = append(*f.operandStack, x)
			return nil
		}, nil
	case 182: // invokevirtual
		return func(s *state, f frame) error {
			/*// TODO: Invoke instance method; dispatch based on class
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
				if out, ok := lastVal.(class); ok && out.name == "java/lang/StringBuilder" {
					println(out.vars["string"].(string))
					return nil
				}
				_ = (*f.operandStack)[:len(*f.operandStack)-1] // objectref
				*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

				return fmt.Errorf(`unimplementet type "%v" for println`, reflect.TypeOf(lastVal))
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

				return fmt.Errorf(`unimplementet type "%v" for print`, reflect.TypeOf(lastVal))
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
				return fmt.Errorf(`function "%s" in class "%s" not implemented (invokevirtual)`, name, className)
			}*/

			// FIXME: this is just a copy of invokestatic (184)
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

			des, err := parseDescriptor(descriptor)
			args := make([]variable, 0)
			pTypes := des.parameterTypes
			for i, j := 0, len(pTypes)-1; i < j; i, j = i+1, j-1 {
				pTypes[i], pTypes[j] = pTypes[j], pTypes[i]
			}
			for _, parameterType := range des.parameterTypes {
				var arg variable
				lastVal := (*f.operandStack)[len(*f.operandStack)-1]
				*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
				switch parameterType {
				case "I":
					arg = lastVal.(int)
				case "Ljava/lang/String":
					arg = lastVal.(string)
				default:
					return fmt.Errorf("unknown parameter type %s", parameterType)
				}
				args = append(args, arg)
			}
			for i, j := 0, len(args)-1; i < j; i, j = i+1, j-1 {
				args[i], args[j] = args[j], args[i]
			}

			err = runMethod(name, descriptor, s, args) //FIXME: This won't work for methods with the same name in different classes
			if err != nil {
				return err
			}

			return err
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

			des, err := parseDescriptor(descriptor)
			args := make([]variable, 0)
			pTypes := des.parameterTypes
			for i, j := 0, len(pTypes)-1; i < j; i, j = i+1, j-1 {
				pTypes[i], pTypes[j] = pTypes[j], pTypes[i]
			}
			for _, parameterType := range des.parameterTypes {
				var arg variable
				lastVal := (*f.operandStack)[len(*f.operandStack)-1]
				*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
				switch parameterType {
				case "I":
					arg = lastVal.(int)
				case "Ljava/lang/String":
					arg = lastVal.(string)
				default:
					return fmt.Errorf("unknown parameter type %s", parameterType)
				}
				args = append(args, arg)
			}
			for i, j := 0, len(args)-1; i < j; i, j = i+1, j-1 {
				args[i], args[j] = args[j], args[i]
			}

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
	case 188: // newarray
		return func(s *state, f frame) error {
			count := (*f.operandStack)[len(*f.operandStack)-1].(int)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			atype, err := f.codeReader.ReadByte()
			if err != nil {
				return err
			}

			switch atype {
			case 5: // T_CHAR
				*f.operandStack = append(*f.operandStack, make([]rune, count))
				return nil
			default:
				return fmt.Errorf("creating new array atype %d not implemented", atype)
			}
		}, nil
	case 190: // arraylength
		return func(s *state, f frame) error {
			lastVal := (*f.operandStack)[len(*f.operandStack)-1].(field)
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			arr := lastVal.value.([]rune) // FIXME: i am making this work for exactly one case xD

			*f.operandStack = append(*f.operandStack, len(arr))
			return nil
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

type descriptor struct {
	parameterTypes []string
	returnType     string
}

func parseDescriptor(descriptorString string) (des descriptor, err error) {
	r := strings.NewReader(descriptorString)
	if c, _, err := r.ReadRune(); c != '(' {
		if err != nil {
			return des, err
		}
		err = fmt.Errorf("error parsing first letter of descriptor")
		return des, err
	}

	for {
		c, err := parseType(r)
		if err != nil {
			return des, err
		}
		if c == ")" {
			break
		}

		des.parameterTypes = append(des.parameterTypes, c)
	}

	c, err := parseType(r)
	if err != nil {
		return des, err
	}

	des.returnType = c

	return
}

func parseType(r *strings.Reader) (string, error) {
	c, _, err := r.ReadRune()
	if err != nil {
		return "", err
	}

	if c == 'L' {
		out := []rune{c}
		for {
			c, _, err := r.ReadRune()
			if c == ';' {
				break
			}
			if err != nil {
				return "", err
			}

			out = append(out, c)
		}
		return string(out), nil
	}
	if c == '[' {
		c2, err := parseType(r)
		if err != nil {
			return "", err
		}

		return "[" + c2, nil
	}

	return string([]rune{c}), nil
}
