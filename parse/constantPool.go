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
	ClassIndex       int  // u2 class_index;
	NameAndTypeIndex int  // u2 name_and_type_index;
}

func (c ConstantMethodrefInfo) toCPInfoRaw() cpInfoRaw {
	//TODO implement me
	panic("implement me")
}

type ConstantFieldrefInfo struct {
	tag              byte // u1 tag; (always 9)
	ClassIndex       int  // u2 class_index;
	NameAndTypeIndex int  // u2 name_and_type_index;
}

func (c ConstantFieldrefInfo) toCPInfoRaw() cpInfoRaw {
	//TODO implement me
	panic("implement me")
}

type ConstantClassInfo struct {
	tag       byte // u1 tag;
	NameIndex int  // u2 name_index;
}

func (c ConstantClassInfo) toCPInfoRaw() cpInfoRaw {
	//TODO implement me
	panic("implement me")
}

type ConstantNameAndTypeInfo struct {
	tag             byte // u1 tag;
	NameIndex       int  // u2 name_index;
	DescriptorIndex int  // u2 descriptor_index;
}

func (c ConstantNameAndTypeInfo) toCPInfoRaw() cpInfoRaw {
	//TODO implement me
	panic("implement me")
}

type ConstantUtf8Info struct {
	tag     byte   // 1 tag;
	length  int    // u2 length;
	Content []byte // u1 bytes[length];
}

func (c ConstantUtf8Info) toCPInfoRaw() cpInfoRaw {
	//TODO implement me
	panic("implement me")
}

func (r *ClassFileReader) ReadCPInfo() (info cpInfo, err error) {
	tag, err := r.ReadByte()
	if err != nil {
		return
	}

	// TODO: implement fully: https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html#jvms-4.4-140
	switch tag {
	case 7: // CONSTANT_Class
		CInfo := ConstantClassInfo{tag: 7}
		CInfo.NameIndex, err = r.ReadU2()
		if err != nil {
			return
		}

		info = CInfo
	case 9: // CONSTANT_Fieldref
		FieldredInfo := ConstantFieldrefInfo{tag: 9}
		FieldredInfo.ClassIndex, err = r.ReadU2()
		if err != nil {
			return
		}
		FieldredInfo.NameAndTypeIndex, err = r.ReadU2()
		if err != nil {
			return
		}

		info = FieldredInfo
	case 10: // CONSTANT_Methodref
		MRInfo := ConstantMethodrefInfo{tag: 10}
		MRInfo.ClassIndex, err = r.ReadU2()
		if err != nil {
			return
		}
		MRInfo.NameAndTypeIndex, err = r.ReadU2()
		if err != nil {
			return
		}

		info = MRInfo
	case 12: // CONSTANT_NameAndType
		NaTInfo := ConstantNameAndTypeInfo{tag: 12}
		NaTInfo.NameIndex, err = r.ReadU2()
		if err != nil {
			return
		}
		NaTInfo.DescriptorIndex, err = r.ReadU2()
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
			Utf8Info.Content = append(Utf8Info.Content, b)
		}
		info = Utf8Info
	default:
		err = fmt.Errorf(`constant pool tag "%d" not implemented`, tag)
		return
	}

	return
}

func (r *ClassFileReader) ReadConstantPool() (size int, entries []cpInfo, err error) {
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
