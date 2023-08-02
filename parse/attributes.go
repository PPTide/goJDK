package parse

type AttributeInfo struct {
	AttributeNameIndex int
	attributeLength    int // u4
	Info               []byte
}

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

func (r *ClassFileReader) ReadAttributeInfo() (info AttributeInfo, err error) {
	info.AttributeNameIndex, err = r.ReadU2()
	if err != nil {
		return
	}

	info.attributeLength, err = r.ReadU4()
	if err != nil {
		return
	}

	for i := 0; i < info.attributeLength; i++ { // TODO: can be made faster by using r.Read() (also see the other instances of me doing this)
		var b byte
		b, err = r.ReadByte()
		info.Info = append(info.Info, b)
	}

	return
}
