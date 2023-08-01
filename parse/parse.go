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
	magic             int
	minorVersion      int
	majorVersion      int
	constantPoolCount int
	constantPool      []cpInfo
	accessFlags       int
	thisClass         int
	superClass        int
	interfacesCount   int
	interfaces        []byte
	fieldsCount       int
	fields            interface{}
	methodsCount      int
	methods           []methodInfo
	attributesCount   int
	attributes        []attributeInfo
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

	cF.magic, err = reader.ReadU4()
	if err != nil {
		return
	}
	if cF.magic != 0xCAFEBABE {
		err = errors.New("incorrect magic")
		return
	}

	cF.minorVersion, err = reader.ReadU2()
	if err != nil {
		return
	}

	cF.majorVersion, err = reader.ReadU2()
	if err != nil {
		return
	}

	cF.constantPoolCount, cF.constantPool, err = reader.ReadConstantPool()
	if err != nil {
		return
	}

	cF.accessFlags, err = reader.ReadU2()
	if err != nil {
		return
	}

	cF.thisClass, err = reader.ReadU2()
	if err != nil {
		return
	}

	cF.superClass, err = reader.ReadU2()
	if err != nil {
		return
	}

	cF.interfacesCount, err = reader.ReadU2()
	if err != nil {
		return
	}

	if cF.interfacesCount > 0 {
		return cF, errors.New("interfaces not implemented")
	}

	cF.fieldsCount, err = reader.ReadU2()
	if err != nil {
		return
	}

	if cF.fieldsCount > 0 {
		return cF, errors.New("fields not implemented")
	}

	cF.methodsCount, cF.methods, err = reader.ReadMethods()
	if err != nil {
		return
	}

	cF.attributesCount, cF.attributes, err = reader.ReadAttributes()
	if err != nil {
		return
	}

	if reader.Len() != 0 {
		return cF, errors.New("couldn't fully read the .class file")
	}

	return
}
