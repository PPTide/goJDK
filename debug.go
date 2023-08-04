package main

import (
	"bytes"
	"fmt"
	"github.com/PPTide/gojdk/parse"
	"github.com/k0kubun/pp"
	"strings"
)

func niceShow(file parse.ClassFile) {
	Constants := make([]string, 0)
	for i, info := range file.ConstantPool { // TODO: switch this to a switch case
		if utf8, ok := info.(parse.ConstantUtf8Info); ok {
			Constants = append(Constants, fmt.Sprintf("Utf8\t\t%s", utf8.Text))
			continue
		}
		if methRef, ok := info.(parse.ConstantMethodrefInfo); ok {
			Constants = append(Constants, fmt.Sprintf("Methodref\t\t#%d.#%d\t\t// %s.%s:%s",
				methRef.ClassIndex,
				methRef.NameAndTypeIndex,
				(*(*methRef.Class).(parse.ConstantClassInfo).Name).(parse.ConstantUtf8Info).Text,
				(*(*methRef.NameAndType).(parse.ConstantNameAndTypeInfo).Name).(parse.ConstantUtf8Info).Text,
				(*(*methRef.NameAndType).(parse.ConstantNameAndTypeInfo).Descriptor).(parse.ConstantUtf8Info).Text,
			))
			continue
		}
		if fieldRef, ok := info.(parse.ConstantFieldrefInfo); ok {
			Constants = append(Constants, fmt.Sprintf("Fieldref\t\t#%d.#%d\t\t// %s.%s:%s",
				fieldRef.ClassIndex,
				fieldRef.NameAndTypeIndex,
				(*(*fieldRef.Class).(parse.ConstantClassInfo).Name).(parse.ConstantUtf8Info).Text,
				(*(*fieldRef.NameAndType).(parse.ConstantNameAndTypeInfo).Name).(parse.ConstantUtf8Info).Text,
				(*(*fieldRef.NameAndType).(parse.ConstantNameAndTypeInfo).Descriptor).(parse.ConstantUtf8Info).Text,
			))
			continue
		}
		if class, ok := info.(parse.ConstantClassInfo); ok {
			Constants = append(Constants, fmt.Sprintf("Class\t\t#%d\t\t// %s",
				class.NameIndex,
				(*class.Name).(parse.ConstantUtf8Info).Text,
			))
			continue
		}
		if nat, ok := info.(parse.ConstantNameAndTypeInfo); ok {
			Constants = append(Constants, fmt.Sprintf("NameAndType\t#%d:%d\t\t// %s:%s",
				nat.NameIndex,
				nat.DescriptorIndex,
				(*nat.Name).(parse.ConstantUtf8Info).Text,
				(*nat.Name).(parse.ConstantUtf8Info).Text,
			))
			continue
		}
		Constants = append(Constants, pp.Sprintf("%v: %v\n", i+1, info))
	}
	fmt.Printf("Constant pool:\n")
	for i, constant := range Constants {
		fmt.Printf("#%d = %s\n", i+1, constant)
	}

	println("Methods: ")
	for _, method := range file.Methods {
		pp.Printf("  Name: %s\n", getAsUtf8String(method.NameIndex, file))
		pp.Printf("  Descriptor: %s\n", getAsUtf8String(method.DescriptorIndex, file))
		for _, attribute := range method.Attributes {
			println("  Attribute: ")
			print(showAttribute(attribute, file, 4))
		}
		println()
	}
	println("Attributes: ")
	for _, attribute := range file.Attributes {
		print(showAttribute(attribute, file, 2))
	}
}

func showAttribute(info parse.AttributeInfo, file parse.ClassFile, pad int) (out string) {
	attributeName := getAsUtf8String(info.AttributeNameIndex, file)
	out += fmt.Sprintf("%sName: %s\n", strings.Repeat(" ", pad), attributeName)

	reader := (*parse.ClassFileReader)(bytes.NewReader(info.Info))

	switch attributeName {
	case "LineNumberTable":
		lineNumberLen, _ := reader.ReadU2()
		for i := 0; i < lineNumberLen; i++ {
			startPc, _ := reader.ReadU2()
			out += fmt.Sprintf("%sStart pc: %d\n", strings.Repeat(" ", pad), startPc)
			lineNumber, _ := reader.ReadU2()
			out += fmt.Sprintf("%sLine Number: %d\n", strings.Repeat(" ", pad), lineNumber)
		}
		if reader.Len() != 0 {
			//panic("couldn't read attributes LineNumberTable")
		}
	case "Code":
		maxStack, _ := reader.ReadU2()
		maxLocals, _ := reader.ReadU2()

		codeLength, _ := reader.ReadU4()
		code := make([]byte, codeLength)
		reader.Read(code)

		exceptionTableLength, _ := reader.ReadU2()
		exceptionTable := make([]byte, exceptionTableLength)
		reader.Read(exceptionTable) //TODO: Parse the exception table

		_, attributes, _ := reader.ReadAttributes()

		out += fmt.Sprintf("%sMax Stack: %d\n", strings.Repeat(" ", pad), maxStack)
		out += fmt.Sprintf("%sMax Locals: %d\n", strings.Repeat(" ", pad), maxLocals)
		out += fmt.Sprintf("%sExeption Table Len: %d\n", strings.Repeat(" ", pad), exceptionTableLength)
		out += fmt.Sprintf("%sCode: %v\n", strings.Repeat(" ", pad), code)
		out += fmt.Sprintf("%sAttributes: \n", strings.Repeat(" ", pad))
		for _, attribute := range attributes {
			out += showAttribute(attribute, file, pad+2)
		}
	case "SourceFile":
		sourcefileIndex, _ := reader.ReadU2()
		out += fmt.Sprintf("%sSource File: %s\n", strings.Repeat(" ", pad), getAsUtf8String(sourcefileIndex, file))
	default:
		out += strings.Repeat(" ", pad) + "Unknown Attribute\n"
	}

	return
}

func getAsUtf8String(index int, file parse.ClassFile) string {
	return string(file.ConstantPool[index-1].(parse.ConstantUtf8Info).Content)
}
