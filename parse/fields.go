package parse

type FieldInfo struct {
	AccessFlags     int
	NameIndex       int
	DescriptorIndex int
	AttributesCount int
	Attributes      []AttributeInfo
}

func (r *ClassFileReader) ReadFields() (size int, entries []FieldInfo, err error) {
	size, err = r.ReadU2()
	if err != nil {
		return
	}

	for i := 0; i < size; i++ {
		var entry FieldInfo
		entry, err = r.ReadFieldInfo()
		entries = append(entries, entry)
		if err != nil {
			return
		}
	}

	return
}

func (r *ClassFileReader) ReadFieldInfo() (info FieldInfo, err error) {
	info.AccessFlags, err = r.ReadU2()
	if err != nil {
		return
	}

	info.NameIndex, err = r.ReadU2()
	if err != nil {
		return
	}

	info.DescriptorIndex, err = r.ReadU2()
	if err != nil {
		return
	}

	info.AttributesCount, info.Attributes, err = r.ReadAttributes()
	if err != nil {
		return
	}

	return
}
