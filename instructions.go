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

func getInstruction(instruction byte) (func(s *state, f frame) error, error) {
	switch instruction {
	case 2:
		return func(s *state, f frame) error { // iconst_m1
			*f.operandStack = append(*f.operandStack, variable{
				valType: "int",
				val:     -1,
			})
			return nil
		}, nil
	case 3:
		return func(s *state, f frame) error { // iconst_0
			*f.operandStack = append(*f.operandStack, variable{
				valType: "int",
				val:     0,
			})
			return nil
		}, nil
	case 4:
		return func(s *state, f frame) error { // iconst_1
			*f.operandStack = append(*f.operandStack, variable{
				valType: "int",
				val:     1,
			})
			return nil
		}, nil
	case 5:
		return func(s *state, f frame) error { // iconst_2
			*f.operandStack = append(*f.operandStack, variable{
				valType: "int",
				val:     2,
			})
			return nil
		}, nil
	case 6:
		return func(s *state, f frame) error { // iconst_3
			*f.operandStack = append(*f.operandStack, variable{
				valType: "int",
				val:     3,
			})
			return nil
		}, nil
	case 7:
		return func(s *state, f frame) error { // iconst_4
			*f.operandStack = append(*f.operandStack, variable{
				valType: "int",
				val:     4,
			})
			return nil
		}, nil
	case 16:
		return func(s *state, f frame) error { // bipush
			b, err := f.codeReader.ReadByte()
			if err != nil {
				return err
			}

			*f.operandStack = append(*f.operandStack, variable{
				valType: "int",
				val:     int(b),
			})
			return nil
		}, nil
	case 18: // ldc
		return func(s *state, f frame) error {
			idx, err := f.codeReader.ReadByte()
			if err != nil {
				return err
			}

			switch t := f.file.ConstantPool[idx-1].(type) {
			case parse.ConstantStringInfo:
				*f.heap = append(*f.heap, (*t.String).(parse.ConstantUtf8Info).Text)
				*f.operandStack = append(*f.operandStack, variable{
					valType: "reference",
					val:     &((*f.heap)[len(*f.heap)-1]),
				})
				return nil
			case parse.ConstantIntegerInfo:
				*f.operandStack = append(*f.operandStack, variable{
					valType: "int",
					val:     t.Integer,
				})
				return nil
			default:
				_ = t
				return fmt.Errorf("ldc (18) not implemented for %s", reflect.TypeOf(f.file.ConstantPool[idx-1]))
			}
		}, nil
	case 21: // iload // TODO: type checking
		return func(s *state, f frame) error {
			b, err := f.codeReader.ReadByte()
			if err != nil {
				return err
			}

			*f.operandStack = append(*f.operandStack, (*f.localVariable)[b])
			return nil
		}, nil
	case 25: // aload // TODO: type checking
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
	case 46: // iaload
		return func(s *state, f frame) error {
			index := f.operandStack.pop().expectType("int").(int)

			arrayref := (*f.operandStack.pop().expectReferenceOfType("[I")).([]int)

			*f.operandStack = append(*f.operandStack, variable{
				valType: "int",
				val:     arrayref[index],
			})

			return nil
		}, nil
	case 52: // caload
		return func(s *state, f frame) error {
			index := f.operandStack.pop().expectType("int").(int)

			arrayref := (*f.operandStack.pop().expectReferenceOfType("[C")).([]rune)

			*f.operandStack = append(*f.operandStack, variable{
				valType: "char",
				val:     arrayref[index],
			})

			return nil
		}, nil
	case 54: // istore // TODO: type checking
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
	case 58: // astore // TODO: type checking
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
	case 85: // castore
		return func(s *state, f frame) error {
			value := f.operandStack.pop().expectType("char").(rune)

			index := f.operandStack.pop().expectType("int").(int)

			arrayref := f.operandStack.pop().expectReferenceOfType("[C")
			array := (*arrayref).([]rune)

			array[index] = value
			*arrayref = array

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
			value2 := f.operandStack.pop().expectType("int").(int)
			value1 := f.operandStack.pop().expectType("int").(int)

			*f.operandStack = append(*f.operandStack, variable{
				valType: "int",
				val:     value1 + value2,
			})
			return nil
		}, nil
	case 100:
		return func(s *state, f frame) error { // isub
			value2 := f.operandStack.pop().expectType("int").(int)
			value1 := f.operandStack.pop().expectType("int").(int)

			*f.operandStack = append(*f.operandStack, asIntVariable(value1-value2))
			return nil
		}, nil
	case 104:
		return func(s *state, f frame) error { // imul
			lastVal := f.operandStack.pop().expectType("int").(int)
			lastVal2 := f.operandStack.pop().expectType("int").(int)

			*f.operandStack = append(*f.operandStack, asIntVariable(lastVal*lastVal2))
			return nil
		}, nil
	case 120: // iushl
		return func(s *state, f frame) error {
			value2 := f.operandStack.pop().expectType("int").(int)
			value1 := f.operandStack.pop().expectType("int").(int)

			*f.operandStack = append(*f.operandStack, asIntVariable(value1<<(value2&0b11111)))
			return nil
		}, nil
	case 124: // iushr
		return func(s *state, f frame) error {
			value2 := f.operandStack.pop().expectType("int").(int)
			value1 := f.operandStack.pop().expectType("int").(int)

			*f.operandStack = append(*f.operandStack, asIntVariable(value1>>(value2&0b11111)))
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

			value1 := (*f.localVariable)[index].expectType("int").(int)
			value2 := int(int8(constVal))

			(*f.localVariable)[index] = asIntVariable(value1 + value2)
			return nil
		}, nil
	case 153: // ifeq -> ==
		return func(s *state, f frame) error {
			lastVal := f.operandStack.pop().expectType("int").(int)

			branchOffset, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}
			branchOffset = branchOffset - 3 // get the offset without the branch bytes

			if lastVal == 0 {
				_, err := f.codeReader.Seek(int64(branchOffset), io.SeekCurrent)
				if err != nil {
					return err
				}
			}
			return nil
		}, nil
	case 154: // ifne -> !=
		return func(s *state, f frame) error {
			lastVal := f.operandStack.pop().expectType("int").(int)

			branchOffset, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}
			branchOffset = branchOffset - 3 // get the offset without the branch bytes

			if lastVal != 0 {
				_, err := f.codeReader.Seek(int64(branchOffset), io.SeekCurrent)
				if err != nil {
					return err
				}
			}
			return nil
		}, nil
	case 156: // ifle -> >=
		return func(s *state, f frame) error {
			lastVal := f.operandStack.pop().expectType("int").(int)

			branchOffset, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}
			branchOffset = branchOffset - 3 // get the offset without the branch bytes

			if lastVal >= 0 {
				_, err := f.codeReader.Seek(int64(branchOffset), io.SeekCurrent)
				if err != nil {
					return err
				}
			}
			return nil
		}, nil
	case 158: // ifle -> <=
		return func(s *state, f frame) error {
			lastVal := f.operandStack.pop().expectType("int").(int)

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
	case 160:
		return func(s *state, f frame) error { // if_icmpne <- !=
			value2 := f.operandStack.pop().expectType("int").(int)
			value1 := f.operandStack.pop().expectType("int").(int)

			branchOffset, err := f.codeReader.ReadU2() //FIXME: might be negative
			if err != nil {
				return err
			}
			branchOffset = branchOffset - 3 // get the offset without the branch bytes

			if value1 != value2 {
				_, err := f.codeReader.Seek(int64(branchOffset), io.SeekCurrent)
				if err != nil {
					return err
				}
			}
			return nil
		}, nil
	case 161:
		return func(s *state, f frame) error { // if_icmplt <- <
			value2 := f.operandStack.pop().expectType("int").(int)
			value1 := f.operandStack.pop().expectType("int").(int)

			branchOffset, err := f.codeReader.ReadU2() //FIXME: might be negative
			if err != nil {
				return err
			}
			branchOffset = branchOffset - 3 // get the offset without the branch bytes

			if value1 < value2 {
				_, err := f.codeReader.Seek(int64(branchOffset), io.SeekCurrent)
				if err != nil {
					return err
				}
			}
			return nil
		}, nil
	case 162:
		return func(s *state, f frame) error { // if_icmpge <- >=
			value2 := f.operandStack.pop().expectType("int").(int)
			value1 := f.operandStack.pop().expectType("int").(int)

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
			value2 := f.operandStack.pop().expectType("int").(int)
			value1 := f.operandStack.pop().expectType("int").(int)

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
	case 176:
		return func(s *state, f frame) error { // areturn
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]

			*(*s).frames[len(s.frames)-2].operandStack = append(*(*s).frames[len(s.frames)-2].operandStack, lastVal)

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

			s.frames = s.frames[:len(s.frames)-1]

			return nil
		}, nil
	case 178: // getstatic
		return func(s *state, f frame) error {
			index, err := f.codeReader.ReadU2()

			val := f.file.ConstantPool[index-1]

			/*if t, ok := val.(parse.ConstantFieldrefInfo); ok && (*(*t.NameAndType).(parse.ConstantNameAndTypeInfo).Name).(parse.ConstantUtf8Info).Text == "sizeTable" { // FIXME: jank xD
				*f.heap = append(*f.heap, []int{9, 99, 999, 9999, 99999, 999999, 9999999, 99999999, 999999999, math.MaxInt64})
				*f.operandStack = append(*f.operandStack, variable{
					referenceType: "[I",
					reference:     &((*f.heap)[len(*f.heap)-1]),
				})
				return nil
			}
			if t, ok := val.(parse.ConstantFieldrefInfo); ok && (*(*t.NameAndType).(parse.ConstantNameAndTypeInfo).Name).(parse.ConstantUtf8Info).Text == "digits" { // FIXME: jank xD
				*f.heap = append(*f.heap, []rune{
					'0', '1', '2', '3', '4', '5',
					'6', '7', '8', '9', 'a', 'b',
					'c', 'd', 'e', 'f', 'g', 'h',
					'i', 'j', 'k', 'l', 'm', 'n',
					'o', 'p', 'q', 'r', 's', 't',
					'u', 'v', 'w', 'x', 'y', 'z',
				})
				*f.operandStack = append(*f.operandStack, variable{
					referenceType: "[C",
					reference:     &((*f.heap)[len(*f.heap)-1]),
				})
				return nil
			}
			*/

			t, ok := val.(parse.ConstantFieldrefInfo)
			if !ok {
				return fmt.Errorf("type %s not of type ConstantFieldrefInfo int getstatic", reflect.TypeOf(val))
			}
			className := (*(*t.Class).(parse.ConstantClassInfo).Name).(parse.ConstantUtf8Info).Text
			fieldName := (*(*t.NameAndType).(parse.ConstantNameAndTypeInfo).Name).(parse.ConstantUtf8Info).Text
			fieldDescriptor := (*(*t.NameAndType).(parse.ConstantNameAndTypeInfo).Descriptor).(parse.ConstantUtf8Info).Text

			ref, err := initializeClass(className, s)
			if err != nil {
				return err
			}

			classContent := (*ref.expectReferenceOfType("L" + className)).(class)
			field := classContent.vars[fieldName]

			if !(field.valType == fieldDescriptor ||
				field.referenceType == fieldDescriptor) {
				return fmt.Errorf("fieldDescriptor doesn't match field type in getstatic")
			}

			*f.operandStack = append(*f.operandStack, field)

			return nil
		}, nil
	case 180: // getfield
		return func(s *state, f frame) error {
			index, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}
			field := f.file.ConstantPool[index-1].(parse.ConstantFieldrefInfo)

			fieldClassName := (*(*field.Class).(parse.ConstantClassInfo).Name).(parse.ConstantUtf8Info).Text
			fieldName := (*(*field.NameAndType).(parse.ConstantNameAndTypeInfo).Name).(parse.ConstantUtf8Info).Text
			fieldDescriptor := (*(*field.NameAndType).(parse.ConstantNameAndTypeInfo).Descriptor).(parse.ConstantUtf8Info).Text

			objectref := f.operandStack.pop().expectReferenceOfType("L" + fieldClassName)
			resolvedObjectref := (*objectref).(class)

			value := resolvedObjectref.vars[fieldName]

			if fieldDescriptor != value.referenceType { // or valType but different xD maybe "L"+referenceType?
				return fmt.Errorf("field and value types don't match: %s != %s", fieldDescriptor, value.referenceType)
			}

			*f.operandStack = append(*f.operandStack, value)

			return nil
		}, nil

	case 181: // putfield
		return func(s *state, f frame) error {
			index, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}
			field := f.file.ConstantPool[index-1].(parse.ConstantFieldrefInfo)

			fieldClassName := (*(*field.Class).(parse.ConstantClassInfo).Name).(parse.ConstantUtf8Info).Text
			fieldName := (*(*field.NameAndType).(parse.ConstantNameAndTypeInfo).Name).(parse.ConstantUtf8Info).Text
			fieldDescriptor := (*(*field.NameAndType).(parse.ConstantNameAndTypeInfo).Descriptor).(parse.ConstantUtf8Info).Text

			value := f.operandStack.pop()
			objectref := f.operandStack.pop().expectReferenceOfType("L" + fieldClassName)
			resolvedObjectref := (*objectref).(class)

			if fieldDescriptor != value.referenceType { // or valType but different xD maybe "L"+referenceType?
				return fmt.Errorf("field and value types don't match: %s != %s", fieldDescriptor, value.referenceType)
			}

			resolvedObjectref.vars[fieldName] = value
			*objectref = resolvedObjectref

			return nil
		}, nil
	case 182: // invokevirtual
		return func(s *state, f frame) error {
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
					fileClass := file.ConstantPool[file.ThisClass-1].(parse.ConstantClassInfo)
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
				lastVal := f.operandStack.pop()
				switch parameterType {
				case "I":
					_ = lastVal.expectType("int")
					arg = lastVal
				case "Z": // boolean
					if !(lastVal.valType == "boolean" || lastVal.valType == "int") {
						panic("Expected val of int or boolean but got " + lastVal.valType)
					}
					arg = lastVal
				case "Ljava/lang/String":
					_ = lastVal.expectReferenceOfType("Ljava/lang/String")
					arg = lastVal
				case "[C":
					_ = lastVal.expectReferenceOfType("[C")
					arg = lastVal
				default:
					return fmt.Errorf("unknown parameter type %s in invokevirtual", parameterType)
				}
				args = append(args, arg)
			}

			// invokevirtual needs to pass a objectref as the 0th arg
			lastVal := (*f.operandStack)[len(*f.operandStack)-1]
			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]
			args = append(args, lastVal) // TODO: typecheck

			for i, j := 0, len(args)-1; i < j; i, j = i+1, j-1 {
				args[i], args[j] = args[j], args[i]
			}

			name = className + "." + name

			err = getClassFile(className, s)
			if err != nil {
				return err
			}

			err = runMethod(name, descriptor, s, args) //FIXME: This won't work for methods with the same name in different classes
			if err != nil {
				return err
			}

			return err
		}, nil
	case 183: // invokespecial
		/*return func(s *state, f frame) error {
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
				return nil // FIXME: THIS WILL NEVER WORK // FIXME: YOU FUCKING NEED TO WORK ON THIS *RIGHT NOW*
			default:
				return fmt.Errorf(`function "%s" not implemented (invokespecial)`, name)
			}
		}, nil*/
		return getInstruction(182) // TODO: tmp
	case 184:
		return func(s *state, f frame) error { // invokestatic
			address, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}

			method := f.file.ConstantPool[address-1].(parse.ConstantMethodrefInfo)
			currentClass := f.file.ConstantPool[f.file.ThisClass-1].(parse.ConstantClassInfo)

			methodClass := (*method.Class).(parse.ConstantClassInfo)
			className := (*methodClass.Name).(parse.ConstantUtf8Info).Text
			currentClassName := (*currentClass.Name).(parse.ConstantUtf8Info).Text
			if className != currentClassName {
				err := getClassFile(className, s)
				if err != nil {
					return err
				}
			}

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
					_ = lastVal.expectType("int")
					arg = lastVal
				case "Ljava/lang/String":
					_ = lastVal.expectReferenceOfType("Ljava/lang/String")
					arg = lastVal
				case "[C":
					_ = lastVal.expectReferenceOfType("[C")
					arg = lastVal
				default:
					return fmt.Errorf("unknown parameter type %s in invokestatic", parameterType)
				}
				args = append(args, arg)
			}
			for i, j := 0, len(args)-1; i < j; i, j = i+1, j-1 {
				args[i], args[j] = args[j], args[i]
			}

			name = className + "." + name

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
			// FIXME: no worky :(
			lastVal := strconv.Itoa(f.operandStack.pop().expectType("int").(int))

			string2Index := bootstrapMethod.bootstrapArguments[0]
			string2 := (*f.file.ConstantPool[string2Index-1].(parse.ConstantStringInfo).String).(parse.ConstantUtf8Info).Text

			*f.operandStack = append(*f.operandStack, createAsReferenceAndAddToHeap("Ljava/lang/String", string2+lastVal, f))

			return nil
		}, nil
	case 187: // new
		return func(s *state, f frame) error {
			address, err := f.codeReader.ReadU2()
			if err != nil {
				return err
			}
			name := (*f.file.ConstantPool[address-1].(parse.ConstantClassInfo).Name).(parse.ConstantUtf8Info).Text

			ref, err := initializeClass(name, s)
			if err != nil {
				return err
			}

			*f.operandStack = append(*f.operandStack, ref)

			return nil
		}, nil
	case 188: // newarray
		return func(s *state, f frame) error {
			atype, err := f.codeReader.ReadByte()
			if err != nil {
				return err
			}

			count := f.operandStack.pop().expectType("int").(int)

			switch atype {
			case 5: // T_CHAR
				*f.operandStack = append(*f.operandStack,
					createAsReferenceAndAddToHeap("[C", make([]rune, count), f))
			case 10: // T_INT
				*f.operandStack = append(*f.operandStack,
					createAsReferenceAndAddToHeap("[I", make([]int, count), f))
			default:
				return fmt.Errorf("unkown atype %d in new array", atype)
			}

			return nil
		}, nil
	case 194: // monitorenter
		return func(s *state, f frame) error {
			// FIXME: I don't plan to support multithreading so this is a stub

			*f.operandStack = (*f.operandStack)[:len(*f.operandStack)-1]

			return nil
		}, nil
	}
	return nil, fmt.Errorf(`unknown instruction "%d"`, instruction)
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

// initializeClass takes the name of a class creates a local representation and runs the <init> method
func initializeClass(name string, s *state) (variable, error) {
	heap := make([]interface{}, 0) // FIXME: the heap should be in state anyway lol
	f := frame{heap: &heap}

	ref := createAsReferenceAndAddToHeap("L"+name, class{
		name: "L" + name,
		vars: make(map[string]variable),
	}, f)

	err := getClassFile(name, s)
	if err != nil {
		return variable{}, err
	}

	//err = runMethod(name+".<clinit>", "()V", s, []variable{ref})
	err = runMethod(name+".<clinit>", "()V", s, []variable{})

	return ref, err
}

// getClassFile looks if the class got parsed and added to the list and adds and parses it if it didn't
func getClassFile(className string, s *state) error {
	var file parse.ClassFile

	for _, f := range s.files {
		fileClass := f.ConstantPool[f.ThisClass-1].(parse.ConstantClassInfo)
		fileClassName := (*fileClass.Name).(parse.ConstantUtf8Info).Text

		if fileClassName == className {
			file = f
			return nil
		}
	}
	// File is not loaded yet
	file, err := parse.Parse(className + ".class")
	if err != nil {
		return err
	}
	s.files = append(s.files, file)

	return nil
}
