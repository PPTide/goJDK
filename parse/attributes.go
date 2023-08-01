package parse

type attributeInfo struct {
	attributeNameIndex int
	attributeLength    int // u4
	info               []byte
}

func (r *classFileReader) ReadAttributes() (size int, entries []attributeInfo, err error) {
	size, err = r.ReadU2()
	if err != nil {
		return
	}

	for i := 0; i < size; i++ {
		for i := 0; i < size; i++ {
			var entry attributeInfo
			entry, err = r.ReadAttributeInfo()
			entries = append(entries, entry)
			if err != nil {
				return
			}
		}
	}

	return
}

func (r *classFileReader) ReadAttributeInfo() (info attributeInfo, err error) {
	info.attributeNameIndex, err = r.ReadU2()
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
		info.info = append(info.info, b)
	}

	return
}
