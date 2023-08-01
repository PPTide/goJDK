package parse

import "fmt"

type cpInfoRaw struct {
	tag  byte   // u1 tag;
	info []byte // u1 info[];
}

type cpInfo interface {
	toCPInfoRaw() cpInfoRaw
}

type ConstantMethodrefInfo struct {
	tag              byte // u1 tag; (always 10)
	classIndex       int  // u2 class_index;
	nameAndTypeIndex int  // u2 name_and_type_index;
}

func (c ConstantMethodrefInfo) toCPInfoRaw() cpInfoRaw {
	//TODO implement me
	panic("implement me")
}

type ConstantClassInfo struct {
	tag       byte // u1 tag;
	nameIndex int  // u2 name_index;
}

func (c ConstantClassInfo) toCPInfoRaw() cpInfoRaw {
	//TODO implement me
	panic("implement me")
}

type ConstantNameAndTypeInfo struct {
	tag             byte // u1 tag;
	nameIndex       int  // u2 name_index;
	descriptorIndex int  // u2 descriptor_index;
}

func (c ConstantNameAndTypeInfo) toCPInfoRaw() cpInfoRaw {
	//TODO implement me
	panic("implement me")
}

type ConstantUtf8Info struct {
	tag     byte   // 1 tag;
	length  int    // u2 length;
	content []byte // u1 bytes[length];
}

func (c ConstantUtf8Info) toCPInfoRaw() cpInfoRaw {
	//TODO implement me
	panic("implement me")
}

func (r *classFileReader) ReadCPInfo() (info cpInfo, err error) {
	tag, err := r.ReadByte()
	if err != nil {
		return
	}

	// TODO: implement fully: https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html#jvms-4.4-140
	switch tag {
	case 7: // CONSTANT_Class
		CInfo := ConstantClassInfo{tag: 7}
		CInfo.nameIndex, err = r.ReadU2()
		if err != nil {
			return
		}

		info = CInfo
	case 10: // CONSTANT_Methodref
		MRInfo := ConstantMethodrefInfo{tag: 10}
		MRInfo.classIndex, err = r.ReadU2()
		if err != nil {
			return
		}
		MRInfo.nameAndTypeIndex, err = r.ReadU2()
		if err != nil {
			return
		}

		info = MRInfo
	case 12: // CONSTANT_NameAndType
		NaTInfo := ConstantNameAndTypeInfo{tag: 12}
		NaTInfo.nameIndex, err = r.ReadU2()
		if err != nil {
			return
		}
		NaTInfo.descriptorIndex, err = r.ReadU2()
		if err != nil {
			return
		}

		info = NaTInfo
	case 1: // CONSTANT_Utf8
		Utf8Info := ConstantUtf8Info{tag: 1}
		Utf8Info.length, err = r.ReadU2()
		if err != nil {
			return
		}
		for i := 0; i < Utf8Info.length; i++ {
			var b byte
			b, err = r.ReadByte()
			if err != nil {
				return
			}
			Utf8Info.content = append(Utf8Info.content, b)
		}
		info = Utf8Info
	default:
		err = fmt.Errorf(`constant pool tag "%d" not implemented`, tag)
		return
	}

	return
}

func (r *classFileReader) ReadConstantPool() (size int, entries []cpInfo, err error) {
	size, err = r.ReadU2()
	if err != nil {
		return
	}

	for i := 0; i < size-1; i++ {
		var entry cpInfo
		entry, err = r.ReadCPInfo()
		if err != nil {
			return
		}
		entries = append(entries, entry)
	}

	return
}
