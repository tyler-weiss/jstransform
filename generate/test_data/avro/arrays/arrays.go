// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.

package arrays

import (
	"github.com/actgardner/gogen-avro/compiler"
	"github.com/actgardner/gogen-avro/container"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/vm/types"
	"io"
)

type Arrays struct {

	// The timestamp when this avro data is written. Useful for identifying the newest row of data sharing keys.
	AvroWriteTime int64

	// This is set to true when the Avro data is recording a delete in the source data.
	AvroDeleted bool
	Heights     []int64
	Parents     []*Parents_record
}

func NewArraysWriter(writer io.Writer, codec container.Codec, recordsPerBlock int64) (*container.Writer, error) {
	str := &Arrays{}
	return container.NewWriter(writer, codec, recordsPerBlock, str.Schema())
}

func DeserializeArrays(r io.Reader) (*Arrays, error) {
	t := NewArrays()

	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	return t, err
}

func NewArrays() *Arrays {
	return &Arrays{}
}

func (r *Arrays) Schema() string {
	return "{\"fields\":[{\"doc\":\"The timestamp when this avro data is written. Useful for identifying the newest row of data sharing keys.\",\"logicalType\":\"timestamp-millis\",\"name\":\"AvroWriteTime\",\"type\":\"long\"},{\"default\":false,\"doc\":\"This is set to true when the Avro data is recording a delete in the source data.\",\"name\":\"AvroDeleted\",\"type\":\"boolean\"},{\"name\":\"heights\",\"type\":{\"items\":{\"type\":\"long\"},\"type\":\"array\"}},{\"name\":\"parents\",\"type\":{\"items\":{\"fields\":[{\"name\":\"count\",\"namespace\":\"parents\",\"type\":\"long\"},{\"name\":\"children\",\"namespace\":\"parents\",\"type\":{\"items\":{\"type\":\"string\"},\"type\":\"array\"}}],\"name\":\"parents_record\",\"namespace\":\"parents\",\"type\":\"record\"},\"type\":\"array\"}}],\"name\":\"Arrays\",\"type\":\"record\"}"
}

func (r *Arrays) SchemaName() string {
	return "Arrays"
}

func (r *Arrays) Serialize(w io.Writer) error {
	return writeArrays(r, w)
}

func (_ *Arrays) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *Arrays) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *Arrays) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *Arrays) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *Arrays) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *Arrays) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *Arrays) SetString(v string)   { panic("Unsupported operation") }
func (_ *Arrays) SetUnionElem(v int64) { panic("Unsupported operation") }
func (r *Arrays) Get(i int) types.Field {
	switch i {
	case 0:
		return (*types.Long)(&r.AvroWriteTime)
	case 1:
		return (*types.Boolean)(&r.AvroDeleted)
	case 2:
		r.Heights = make([]int64, 0)
		return (*ArrayLongWrapper)(&r.Heights)
	case 3:
		r.Parents = make([]*Parents_record, 0)
		return (*ArrayParents_recordWrapper)(&r.Parents)

	}
	panic("Unknown field index")
}
func (r *Arrays) SetDefault(i int) {
	switch i {
	case 1:
		r.AvroDeleted = false
		return

	}
	panic("Unknown field index")
}
func (_ *Arrays) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *Arrays) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *Arrays) Finalize()                        {}

type ArraysReader struct {
	r io.Reader
	p *vm.Program
}

func NewArraysReader(r io.Reader) (*ArraysReader, error) {
	containerReader, err := container.NewReader(r)
	if err != nil {
		return nil, err
	}

	t := NewArrays()
	deser, err := compiler.CompileSchemaBytes([]byte(containerReader.AvroContainerSchema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	return &ArraysReader{
		r: containerReader,
		p: deser,
	}, nil
}

func (r *ArraysReader) Read() (*Arrays, error) {
	t := NewArrays()
	err := vm.Eval(r.r, r.p, t)
	return t, err
}
