// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.

package simple

import (
	"github.com/actgardner/gogen-avro/compiler"
	"github.com/actgardner/gogen-avro/container"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/vm/types"
	"io"
)

type Simple struct {

	// The timestamp when this avro data is written. Useful for identifying the newest row of data sharing keys.
	AvroWriteTime int64

	// This is set to true when the Avro data is recording a delete in the source data.
	AvroDeleted bool
	Height      *UnionNullLong
	SomeDateObj *UnionNullSomeDateObj_record
	Type        string
	Visible     bool
	Width       *UnionNullDouble
}

func NewSimpleWriter(writer io.Writer, codec container.Codec, recordsPerBlock int64) (*container.Writer, error) {
	str := &Simple{}
	return container.NewWriter(writer, codec, recordsPerBlock, str.Schema())
}

func DeserializeSimple(r io.Reader) (*Simple, error) {
	t := NewSimple()

	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	return t, err
}

func NewSimple() *Simple {
	return &Simple{}
}

func (r *Simple) Schema() string {
	return "{\"fields\":[{\"doc\":\"The timestamp when this avro data is written. Useful for identifying the newest row of data sharing keys.\",\"logicalType\":\"timestamp-millis\",\"name\":\"AvroWriteTime\",\"type\":\"long\"},{\"default\":false,\"doc\":\"This is set to true when the Avro data is recording a delete in the source data.\",\"name\":\"AvroDeleted\",\"type\":\"boolean\"},{\"name\":\"height\",\"type\":[\"null\",\"long\"]},{\"name\":\"someDateObj\",\"type\":[\"null\",{\"fields\":[{\"name\":\"dates\",\"namespace\":\"someDateObj\",\"type\":{\"items\":{\"logicalType\":\"timestamp-millis\",\"type\":\"long\"},\"type\":\"array\"}}],\"name\":\"someDateObj_record\",\"namespace\":\"someDateObj\",\"type\":\"record\"}]},{\"name\":\"type\",\"type\":\"string\"},{\"default\":false,\"name\":\"visible\",\"type\":\"boolean\"},{\"name\":\"width\",\"type\":[\"null\",\"double\"]}],\"name\":\"Simple\",\"type\":\"record\"}"
}

func (r *Simple) SchemaName() string {
	return "Simple"
}

func (r *Simple) Serialize(w io.Writer) error {
	return writeSimple(r, w)
}

func (_ *Simple) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *Simple) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *Simple) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *Simple) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *Simple) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *Simple) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *Simple) SetString(v string)   { panic("Unsupported operation") }
func (_ *Simple) SetUnionElem(v int64) { panic("Unsupported operation") }
func (r *Simple) Get(i int) types.Field {
	switch i {
	case 0:
		return (*types.Long)(&r.AvroWriteTime)
	case 1:
		return (*types.Boolean)(&r.AvroDeleted)
	case 2:
		r.Height = NewUnionNullLong()
		return r.Height
	case 3:
		r.SomeDateObj = NewUnionNullSomeDateObj_record()
		return r.SomeDateObj
	case 4:
		return (*types.String)(&r.Type)
	case 5:
		return (*types.Boolean)(&r.Visible)
	case 6:
		r.Width = NewUnionNullDouble()
		return r.Width

	}
	panic("Unknown field index")
}
func (r *Simple) SetDefault(i int) {
	switch i {
	case 1:
		r.AvroDeleted = false
		return
	case 5:
		r.Visible = false
		return

	}
	panic("Unknown field index")
}
func (_ *Simple) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *Simple) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *Simple) Finalize()                        {}

type SimpleReader struct {
	r io.Reader
	p *vm.Program
}

func NewSimpleReader(r io.Reader) (*SimpleReader, error) {
	containerReader, err := container.NewReader(r)
	if err != nil {
		return nil, err
	}

	t := NewSimple()
	deser, err := compiler.CompileSchemaBytes([]byte(containerReader.AvroContainerSchema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	return &SimpleReader{
		r: containerReader,
		p: deser,
	}, nil
}

func (r *SimpleReader) Read() (*Simple, error) {
	t := NewSimple()
	err := vm.Eval(r.r, r.p, t)
	return t, err
}
