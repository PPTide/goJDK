package parse

import "fmt"

type CpInfo interface {
	implementCPInfo()
}

type ConstantMethodrefInfo struct {
	tag              byte    // u1 tag; (always 10)
	Class            *CpInfo // ConstantClassInfo
	ClassIndex       int     // u2 class_index;
	NameAndType      *CpInfo // ConstantNameAndTypeInfo
	NameAndTypeIndex int     // u2 name_and_type_index;
}

func (c ConstantMethodrefInfo) implementCPInfo() {
	return
}

type ConstantFieldrefInfo struct {
	tag              byte    // u1 tag; (always 9)
	Class            *CpInfo // ConstantClassInfo
	ClassIndex       int     // u2 class_index;
	NameAndType      *CpInfo // ConstantNameAndTypeInfo
	NameAndTypeIndex int     // u2 name_and_type_index;
}

func (c ConstantFieldrefInfo) implementCPInfo() {
	return
}

type ConstantClassInfo struct {
	tag       byte    // u1 tag;
	Name      *CpInfo // ConstantUtf8Info
	NameIndex int     // u2 name_index;
}

func (c ConstantClassInfo) implementCPInfo() {
	return
}

type ConstantNameAndTypeInfo struct {
	tag             byte    // u1 tag;
	Name            *CpInfo // ConstantUtf8Info
	NameIndex       int     // u2 name_index;
	Descriptor      *CpInfo // ConstantUtf8Info
	DescriptorIndex int     // u2 descriptor_index;
}

func (c ConstantNameAndTypeInfo) implementCPInfo() {
	return
}

type ConstantStringInfo struct {
	tag         byte    // u1 tag;
	String      *CpInfo // ConstantUtf8Info
	StringIndex int     // u2 string_index;

}

func (c ConstantStringInfo) implementCPInfo() {
	return
}

type ConstantUtf8Info struct {
	tag     byte // 1 tag;
	length  int  // u2 length;
	Text    string
	Content []byte // u1 bytes[length];
}

func (c ConstantUtf8Info) implementCPInfo() {
	return
}

type ConstantInvokeDynamicInfo struct {
	tag byte
	//bootstrapMethodAttr *CpInfo //
	BootstrapMethodAttrIndex int     // u2
	NameAndType              *CpInfo // ConstantNameAndTypeInfo, ConstantMethodrefInfo or ConstantInterfaceMethodrefInfo
	NameAndTypeIndex         int     // u2
}

func (c ConstantInvokeDynamicInfo) implementCPInfo() {
	return
}

type ConstantMethodHandleInfo struct {
	tag            byte
	ReferenceKind  byte    // u1
	Reference      *CpInfo // one of ConstantFieldrefInfo,
	ReferenceIndex int     // u2
}

func (c ConstantMethodHandleInfo) implementCPInfo() {
	return
}

// ReadCPInfo reads one constant pool entry.
func (r *ClassFileReader) ReadCPInfo() (info CpInfo, err error) {
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
	case 8: // CONSTANT_String
		StringInfo := ConstantStringInfo{tag: 8}
		StringInfo.StringIndex, err = r.ReadU2()
		if err != nil {
			return
		}

		info = StringInfo
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
		Utf8Info.Content = make([]byte, Utf8Info.length)
		_, err = r.Read(Utf8Info.Content)
		if err != nil {
			return
		}
		info = Utf8Info
	case 18: // CONSTANT_InvokeDynamic
		InvokeDynamixInfo := ConstantInvokeDynamicInfo{tag: 18}
		InvokeDynamixInfo.BootstrapMethodAttrIndex, err = r.ReadU2()
		if err != nil {
			return
		}
		InvokeDynamixInfo.NameAndTypeIndex, err = r.ReadU2()
		info = InvokeDynamixInfo
	case 15: // CONSTANT_MethodHandle
		MethodHandleInfo := ConstantMethodHandleInfo{tag: 15}
		MethodHandleInfo.ReferenceKind, err = r.ReadByte()
		if err != nil {
			return
		}
		MethodHandleInfo.ReferenceIndex, err = r.ReadU2()
		if err != nil {
			return
		}
		info = MethodHandleInfo
	default:
		err = fmt.Errorf(`constant pool tag "%d" not implemented`, tag)
		return
	}

	return
}

// ReadConstantPool reads all entries in a constant pool and return the count and entries.
func (r *ClassFileReader) ReadConstantPool() (size int, entries []CpInfo, err error) {
	size, err = r.ReadU2()
	if err != nil {
		return
	}

	// For some reason this is the one place where the length is one bigger then it needs to be?
	for i := 0; i < size-1; i++ {
		var entry CpInfo
		entry, err = r.ReadCPInfo()
		if err != nil {
			return
		}
		entries = append(entries, entry)
	}

	return
}
