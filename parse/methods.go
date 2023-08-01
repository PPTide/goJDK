package parse

type methodInfo struct {
	accessFlags     int
	nameIndex       int
	descriptorIndex int
	attributesCount int
	attributes      []attributeInfo
}

func (r *classFileReader) ReadMethods() (size int, entries []methodInfo, err error) {
	size, err = r.ReadU2()
	if err != nil {
		return
	}

	for i := 0; i < size; i++ {
		for i := 0; i < size-1; i++ {
			var entry methodInfo
			entry, err = r.ReadMethodInfo()
			entries = append(entries, entry)
			if err != nil {
				return
			}
		}
	}

	return
}

func (r *classFileReader) ReadMethodInfo() (info methodInfo, err error) {
	info.accessFlags, err = r.ReadU2()
	if err != nil {
		return
	}

	info.nameIndex, err = r.ReadU2()
	if err != nil {
		return
	}

	info.descriptorIndex, err = r.ReadU2()
	if err != nil {
		return
	}

	info.attributesCount, info.attributes, err = r.ReadAttributes()
	if err != nil {
		return
	}

	return
}
