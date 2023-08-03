package parse

type AttributeInfo struct {
	AttributeNameIndex int
	attributeLength    int // u4
	Info               []byte
}

// ReadAttributes reads attributes of a class/method/... and returns the size and all attributes.
func (r *ClassFileReader) ReadAttributes() (size int, entries []AttributeInfo, err error) {
	size, err = r.ReadU2()
	if err != nil {
		return
	}

	for i := 0; i < size; i++ {
		var entry AttributeInfo
		entry, err = r.ReadAttributeInfo()
		entries = append(entries, entry)
		if err != nil {
			return
		}
	}

	return
}

// ReadAttributeInfo reads all information for one attribute.
func (r *ClassFileReader) ReadAttributeInfo() (info AttributeInfo, err error) { //TODO: differentiate between different Attributes here
	info.AttributeNameIndex, err = r.ReadU2()
	if err != nil {
		return
	}

	info.attributeLength, err = r.ReadU4()
	if err != nil {
		return
	}

	info.Info = make([]byte, info.attributeLength)
	_, err = r.Read(info.Info)
	if err != nil {
		return
	}

	return
}
