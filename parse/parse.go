package parse

import (
	"bytes"
	"errors"
	"io"
	"math"
	"os"
)

type classFileReader bytes.Reader

func (r *classFileReader) Len() int {
	return ((*bytes.Reader)(r)).Len()
}

func (r *classFileReader) Read(p []byte) (n int, err error) {
	return ((*bytes.Reader)(r)).Read(p)
}

func (r *classFileReader) ReadByte() (byte, error) {
	return ((*bytes.Reader)(r)).ReadByte()
}

func (r *classFileReader) ReadU4() (res int, err error) {
	x := make([]byte, 4)

	n, err := r.Read(x)

	if n != 4 {
		err = errors.Join(err, errors.New("couldn't read 4 bytes"))
	}

	res = decodeBigEndian(x)

	return
}

func (r *classFileReader) ReadU2() (res int, err error) {
	x := make([]byte, 2)

	n, err := r.Read(x)

	if n != 2 {
		err = errors.Join(err, errors.New("couldn't read 2 bytes"))
	}

	res = decodeBigEndian(x)

	return
}

func decodeBigEndian(b []byte) (o int) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	for i, x := range b {
		o += int(x) * int(math.Pow(256, float64(i)))
	}
	return
}

type ClassFile struct {
	Magic             int
	MinorVersion      int
	MajorVersion      int
	ConstantPoolCount int
	ConstantPool      []cpInfo
	AccessFlags       int
	ThisClass         int
	SuperClass        int
	InterfacesCount   int
	Interfaces        []byte
	FieldsCount       int
	Fields            interface{}
	MethodsCount      int
	Methods           []methodInfo
	AttributesCount   int
	Attributes        []attributeInfo
}

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
	reader := (*classFileReader)(bytes.NewReader(content))

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

	if reader.Len() != 0 {
		return cF, errors.New("couldn't fully read the .class file")
	}

	return
}
