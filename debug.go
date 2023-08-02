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
	for i, info := range file.ConstantPool {
		if utf8, ok := info.(parse.ConstantUtf8Info); ok {
			Constants = append(Constants, fmt.Sprintf("Utf8\t\t%s", string(utf8.Content)))
			continue
		}
		if methRef, ok := info.(parse.ConstantMethodrefInfo); ok {
			Constants = append(Constants, fmt.Sprintf("Methodref\t\t#%d.#%d\t\t// %s.%s:%s",
				methRef.ClassIndex,
				methRef.NameAndTypeIndex,
				getAsUtf8String(file.ConstantPool[methRef.ClassIndex-1].(parse.ConstantClassInfo).NameIndex, file),
				getAsUtf8String(file.ConstantPool[methRef.NameAndTypeIndex-1].(parse.ConstantNameAndTypeInfo).NameIndex, file),
				getAsUtf8String(file.ConstantPool[methRef.NameAndTypeIndex-1].(parse.ConstantNameAndTypeInfo).DescriptorIndex, file),
			))
			continue
		}
		if fieldRef, ok := info.(parse.ConstantFieldrefInfo); ok {
			Constants = append(Constants, fmt.Sprintf("Fieldref\t\t#%d.#%d\t\t// %s.%s:%s",
				fieldRef.ClassIndex,
				fieldRef.NameAndTypeIndex,
				getAsUtf8String(file.ConstantPool[fieldRef.ClassIndex-1].(parse.ConstantClassInfo).NameIndex, file),
				getAsUtf8String(file.ConstantPool[fieldRef.NameAndTypeIndex-1].(parse.ConstantNameAndTypeInfo).NameIndex, file),
				getAsUtf8String(file.ConstantPool[fieldRef.NameAndTypeIndex-1].(parse.ConstantNameAndTypeInfo).DescriptorIndex, file),
			))
			continue
		}
		if class, ok := info.(parse.ConstantClassInfo); ok {
			Constants = append(Constants, fmt.Sprintf("Class\t\t#%d\t\t// %s",
				class.NameIndex,
				getAsUtf8String(class.NameIndex, file),
			))
			continue
		}
		if nat, ok := info.(parse.ConstantNameAndTypeInfo); ok {
			Constants = append(Constants, fmt.Sprintf("NameAndType\t#%d:%d\t\t// %s:%s",
				nat.NameIndex,
				nat.DescriptorIndex,
				getAsUtf8String(nat.NameIndex, file),
				getAsUtf8String(nat.DescriptorIndex, file),
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
	attributeName := string(file.ConstantPool[info.AttributeNameIndex].(parse.ConstantUtf8Info).Content)
	//attributeName := getAsUtf8String(info.AttributeNameIndex, file)
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
	default:
		out += strings.Repeat(" ", pad) + "Unknown Attribute\n"
	}

	return
}

func getAsUtf8String(index int, file parse.ClassFile) string {
	return string(file.ConstantPool[index-1].(parse.ConstantUtf8Info).Content)
}
