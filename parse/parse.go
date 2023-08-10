package parse

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

// ClassFileReader extends bytes.Reader to allow java specific readers.
type ClassFileReader bytes.Reader

// Len calls *bytes.Reader.Len()
func (r *ClassFileReader) Len() int {
	return ((*bytes.Reader)(r)).Len()
}

// Read calls *bytes.Reader.Read()
func (r *ClassFileReader) Read(p []byte) (n int, err error) {
	return ((*bytes.Reader)(r)).Read(p)
}

// ReadByte calls *bytes.Reader.ReadByte()
func (r *ClassFileReader) ReadByte() (byte, error) {
	return ((*bytes.Reader)(r)).ReadByte()
}

// ReadU4 reads 4 bytes and interprets them as a big endian int.
func (r *ClassFileReader) ReadU4() (res int, err error) {
	x := make([]byte, 4)

	n, err := r.Read(x)

	if n != 4 {
		err = errors.Join(err, errors.New("couldn't read 4 bytes"))
	}

	res = decodeBigEndian(x)

	return
}

// ReadU2 reads 2 bytes and interprets them as a big endian int.
func (r *ClassFileReader) ReadU2() (res int, err error) {
	x := make([]byte, 2)

	n, err := r.Read(x)

	if n != 2 {
		err = errors.Join(err, errors.New("couldn't read 2 bytes"))
	}

	res = decodeBigEndian(x)

	return
}

// Seek calls bytes.Reader.Seek() (implement io.Seeker)
func (r *ClassFileReader) Seek(offset int64, whence int) (int64, error) {
	return ((*bytes.Reader)(r)).Seek(offset, whence)
}

// decodeBigEndian reads multiple bytes as a big endian number and returns an int.
func decodeBigEndian(b []byte) (o int) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	for i, x := range b {
		o += int(x) << (i * 8)
	}
	return
}

// ClassFile stores all the inner parts of a Java .class file.
type ClassFile struct {
	Magic             int
	MinorVersion      int
	MajorVersion      int
	ConstantPoolCount int
	ConstantPool      []CpInfo
	AccessFlags       int
	ThisClass         int
	SuperClass        int
	InterfacesCount   int
	Interfaces        []byte
	FieldsCount       int
	Fields            interface{}
	MethodsCount      int
	Methods           []MethodInfo
	AttributesCount   int
	Attributes        []AttributeInfo
}

func (f *ClassFile) resolveIndexes() error {
	for i, info := range f.ConstantPool {
		switch t := info.(type) {
		case ConstantUtf8Info:
			t.Text = string(t.Content)
			f.ConstantPool[i] = t
		case ConstantNameAndTypeInfo:
			if _, ok := f.ConstantPool[t.NameIndex-1].(ConstantUtf8Info); !ok {
				return fmt.Errorf("name not of type ConstantUtf8Info")
			}
			t.Name = &f.ConstantPool[t.NameIndex-1]
			if _, ok := f.ConstantPool[t.DescriptorIndex-1].(ConstantUtf8Info); !ok {
				return fmt.Errorf("name not of type ConstantUtf8Info")
			}
			t.Descriptor = &f.ConstantPool[t.DescriptorIndex-1]
			f.ConstantPool[i] = t
		case ConstantClassInfo:
			if _, ok := f.ConstantPool[t.NameIndex-1].(ConstantUtf8Info); !ok {
				return fmt.Errorf("name not of type ConstantUtf8Info")
			}
			t.Name = &f.ConstantPool[t.NameIndex-1]
			f.ConstantPool[i] = t
		case ConstantMethodrefInfo:
			if _, ok := f.ConstantPool[t.ClassIndex-1].(ConstantClassInfo); !ok {
				return fmt.Errorf("name not of type ConstantClassInfo")
			}
			t.Class = &f.ConstantPool[t.ClassIndex-1]
			if _, ok := f.ConstantPool[t.NameAndTypeIndex-1].(ConstantNameAndTypeInfo); !ok {
				return fmt.Errorf("name not of type ConstantNameAndTypeInfo")
			}
			t.NameAndType = &f.ConstantPool[t.NameAndTypeIndex-1]
			f.ConstantPool[i] = t
		case ConstantFieldrefInfo:
			if _, ok := f.ConstantPool[t.ClassIndex-1].(ConstantClassInfo); !ok {
				return fmt.Errorf("name not of type ConstantClassInfo")
			}
			t.Class = &f.ConstantPool[t.ClassIndex-1]
			if _, ok := f.ConstantPool[t.NameAndTypeIndex-1].(ConstantNameAndTypeInfo); !ok {
				return fmt.Errorf("name not of type ConstantNameAndTypeInfo")
			}
			t.NameAndType = &f.ConstantPool[t.NameAndTypeIndex-1]
			f.ConstantPool[i] = t
		case ConstantStringInfo:
			if _, ok := f.ConstantPool[t.StringIndex-1].(ConstantUtf8Info); !ok {
				return fmt.Errorf("name not of type ConstantUft8Info")
			}
			t.String = &f.ConstantPool[t.StringIndex-1]
			f.ConstantPool[i] = t
		case ConstantInvokeDynamicInfo:
			if _, ok := f.ConstantPool[t.NameAndTypeIndex-1].(ConstantNameAndTypeInfo); !ok {
				return fmt.Errorf("name not of type ConstantNameAndTypeInfo")
			}
			t.NameAndType = &f.ConstantPool[t.NameAndTypeIndex-1]
			f.ConstantPool[i] = t
		case ConstantMethodHandleInfo:
			//TODO: Check reference using reference kind
			t.Reference = &f.ConstantPool[t.ReferenceIndex]
			f.ConstantPool[i] = t
		default:
			return fmt.Errorf("unkown type trying to resolve indexes")
		}
	}
	//TODO: add the same for interfaces, fields, methods and attributes
	return nil
}

// Parse a file and return a ClassFile.
func Parse(filename string) (cF ClassFile, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	content, err := io.ReadAll(file)
	if err != nil {
		return
	}
	reader := (*ClassFileReader)(bytes.NewReader(content))

	cF.Magic, err = reader.ReadU4()
	if err != nil {
		return
	}
	if cF.Magic != 0xCAFEBABE {
		err = errors.New("incorrect magic")
		return
	}

	cF.MinorVersion, err = reader.ReadU2()
	if err != nil {
		return
	}

	cF.MajorVersion, err = reader.ReadU2()
	if err != nil {
		return
	}

	cF.ConstantPoolCount, cF.ConstantPool, err = reader.ReadConstantPool()
	if err != nil {
		return
	}

	cF.AccessFlags, err = reader.ReadU2()
	if err != nil {
		return
	}

	cF.ThisClass, err = reader.ReadU2()
	if err != nil {
		return
	}

	cF.SuperClass, err = reader.ReadU2()
	if err != nil {
		return
	}

	cF.InterfacesCount, err = reader.ReadU2()
	if err != nil {
		return
	}

	if cF.InterfacesCount > 0 {
		return cF, errors.New("interfaces not implemented")
	}

	cF.FieldsCount, err = reader.ReadU2()
	if err != nil {
		return
	}

	if cF.FieldsCount > 0 {
		return cF, errors.New("fields not implemented")
	}

	cF.MethodsCount, cF.Methods, err = reader.ReadMethods()
	if err != nil {
		return
	}

	cF.AttributesCount, cF.Attributes, err = reader.ReadAttributes()
	if err != nil {
		return
	}

	// There shouldn't be any data left in the file at this point
	if reader.Len() != 0 {
		return cF, errors.New("couldn't fully read the .class file")
	}

	err = cF.resolveIndexes()
	if err != nil {
		return
	}

	return
}
