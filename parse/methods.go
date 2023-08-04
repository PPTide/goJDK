package parse

type MethodInfo struct {
	AccessFlags     int
	NameIndex       int
	DescriptorIndex int
	AttributesCount int
	Attributes      []AttributeInfo
}

func (r *ClassFileReader) ReadMethods() (size int, entries []MethodInfo, err error) {
	size, err = r.ReadU2()
	if err != nil {
		return
	}

	for i := 0; i < size; i++ {
		var entry MethodInfo
		entry, err = r.ReadMethodInfo()
		entries = append(entries, entry)
		if err != nil {
			return
		}
	}

	return
}

func (r *ClassFileReader) ReadMethodInfo() (info MethodInfo, err error) {
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
