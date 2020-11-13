package bin

import (
	"encoding/binary"
	"encoding/hex"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecoder_Remaining(t *testing.T) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b, 1)
	binary.LittleEndian.PutUint16(b[2:], 2)

	d := NewDecoder(b)

	n, err := d.ReadUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(1), n)
	assert.Equal(t, 2, d.remaining())

	n, err = d.ReadUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(2), n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_int8(t *testing.T) {
	buf := []byte{
		0x9d, // -99
		0x64, // 100
	}

	d := NewDecoder(buf)

	n, err := d.ReadInt8()
	assert.NoError(t, err)
	assert.Equal(t, int8(-99), n)
	assert.Equal(t, 1, d.remaining())

	n, err = d.ReadInt8()
	assert.NoError(t, err)
	assert.Equal(t, int8(100), n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_int16(t *testing.T) {
	buf := []byte{
		0xae, 0xff, // -82
		0x49, 0x00, // 73
	}

	d := NewDecoder(buf)

	n, err := d.ReadInt16()
	assert.NoError(t, err)
	assert.Equal(t, int16(-82), n)
	assert.Equal(t, 2, d.remaining())

	n, err = d.ReadInt16()
	assert.NoError(t, err)
	assert.Equal(t, int16(73), n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_int32(t *testing.T) {
	buf := []byte{
		0xd8, 0x8d, 0x8a, 0xef, // -276132392
		0x4f, 0x9f, 0x3, 0x00, // 237391
	}

	d := NewDecoder(buf)

	n, err := d.ReadInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(-276132392), n)
	assert.Equal(t, 4, d.remaining())

	n, err = d.ReadInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(237391), n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_int64(t *testing.T) {
	buf := []byte{
		0x91, 0x7d, 0xf3, 0xff, 0xff, 0xff, 0xff, 0xff, //-819823
		0xe3, 0x1c, 0x1, 0x00, 0x00, 0x00, 0x00, 0x00, //72931
	}

	d := NewDecoder(buf)

	n, err := d.ReadInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(-819823), n)
	assert.Equal(t, 8, d.remaining())

	n, err = d.ReadInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(72931), n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_uint8(t *testing.T) {
	buf := []byte{
		0x63, // 99
		0x64, // 100
	}

	d := NewDecoder(buf)

	n, err := d.ReadUint8()
	assert.NoError(t, err)
	assert.Equal(t, uint8(99), n)
	assert.Equal(t, 1, d.remaining())

	n, err = d.ReadUint8()
	assert.NoError(t, err)
	assert.Equal(t, uint8(100), n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_uint16(t *testing.T) {
	buf := []byte{
		0x52, 0x00, // 82
		0x49, 0x00, // 73
	}

	d := NewDecoder(buf)

	n, err := d.ReadUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(82), n)
	assert.Equal(t, 2, d.remaining())

	n, err = d.ReadUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(73), n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_uint32(t *testing.T) {
	buf := []byte{
		0x28, 0x72, 0x75, 0x10, // 276132392 as LE
		0x4f, 0x9f, 0x03, 0x00, // 237391 as LE
	}

	d := NewDecoder(buf)

	n, err := d.ReadUint32()
	assert.NoError(t, err)
	assert.Equal(t, uint32(276132392), n)
	assert.Equal(t, 4, d.remaining())

	n, err = d.ReadUint32()
	assert.NoError(t, err)
	assert.Equal(t, uint32(237391), n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_uint64(t *testing.T) {
	buf := []byte{
		0x6f, 0x82, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, //819823
		0xe3, 0x1c, 0x1, 0x00, 0x00, 0x00, 0x00, 0x00, //72931
	}

	d := NewDecoder(buf)

	n, err := d.ReadUint64()
	assert.NoError(t, err)
	assert.Equal(t, uint64(819823), n)
	assert.Equal(t, 8, d.remaining())

	n, err = d.ReadUint64()
	assert.NoError(t, err)
	assert.Equal(t, uint64(72931), n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_float32(t *testing.T) {
	buf := []byte{
		0xc3, 0xf5, 0xa8, 0x3f,
		0xa4, 0x70, 0x4d, 0xc0,
	}

	d := NewDecoder(buf)

	n, err := d.ReadFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(1.32), n)
	assert.Equal(t, 4, d.remaining())

	n, err = d.ReadFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(-3.21), n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_float64(t *testing.T) {
	buf := []byte{
		0x3d, 0x0a, 0xd7, 0xa3, 0x70, 0x1d, 0x4f, 0xc0,
		0x77, 0xbe, 0x9f, 0x1a, 0x2f, 0x3d, 0x37, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf0, 0x7f,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf0, 0xff,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf8, 0x7f,
	}

	d := NewDecoder(buf)

	n, err := d.ReadFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(-62.23), n)
	assert.Equal(t, 32, d.remaining())

	n, err = d.ReadFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(23.239), n)
	assert.Equal(t, 24, d.remaining())

	n, err = d.ReadFloat64()
	assert.NoError(t, err)
	assert.Equal(t, math.Inf(1), n)
	assert.Equal(t, 16, d.remaining())

	n, err = d.ReadFloat64()
	assert.NoError(t, err)
	assert.Equal(t, math.Inf(-1), n)
	assert.Equal(t, 8, d.remaining())

	n, err = d.ReadFloat64()
	assert.NoError(t, err)
	assert.True(t, math.IsNaN(n))
}

func TestDecoder_string(t *testing.T) {
	buf := []byte{
		0x03, 0x31, 0x32, 0x33, // "123"
		0x00,                   // ""
		0x03, 0x61, 0x62, 0x63, // "abc
	}

	d := NewDecoder(buf)

	s, err := d.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, "123", s)
	assert.Equal(t, 5, d.remaining())

	s, err = d.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, "", s)
	assert.Equal(t, 4, d.remaining())

	s, err = d.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, "abc", s)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Decode_String_Err(t *testing.T) {
	buf := []byte{
		0x0a,
	}

	decoder := NewDecoder(buf)

	var s string
	err := decoder.Decode(&s)
	assert.EqualError(t, err, "byte array: varlen=10, missing 10 bytes")
}

func TestDecoder_Byte(t *testing.T) {
	buf := []byte{
		0x00, 0x01,
	}

	d := NewDecoder(buf)

	n, err := d.ReadByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(0), n)
	assert.Equal(t, 1, d.remaining())

	n, err = d.ReadByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(1), n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_Bool(t *testing.T) {
	buf := []byte{
		0x01, 0x00,
	}

	d := NewDecoder(buf)

	n, err := d.ReadBool()
	assert.NoError(t, err)
	assert.Equal(t, true, n)
	assert.Equal(t, 1, d.remaining())

	n, err = d.ReadBool()
	assert.NoError(t, err)
	assert.Equal(t, false, n)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_ByteArray(t *testing.T) {
	buf := []byte{
		0x03, 0x01, 0x02, 0x03,
		0x03, 0x04, 0x05, 0x06,
	}

	d := NewDecoder(buf)

	data, err := d.ReadByteArray()
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, data)
	assert.Equal(t, 4, d.remaining())

	data, err = d.ReadByteArray()
	assert.Equal(t, []byte{4, 5, 6}, data)
	assert.Equal(t, 0, d.remaining())
}

func TestDecoder_ByteArray_MissingData(t *testing.T) {
	buf := []byte{
		0x0a,
	}

	d := NewDecoder(buf)

	_, err := d.ReadByteArray()
	assert.EqualError(t, err, "byte array: varlen=10, missing 10 bytes")
}

func TestDecoder_Array(t *testing.T) {
	buf := []byte{1, 2, 4}

	decoder := NewDecoder(buf)

	var decoded [3]byte
	decoder.Decode(&decoded)
	assert.Equal(t, [3]byte{1, 2, 4}, decoded)
}

func TestDecoder_Array_Err(t *testing.T) {

	decoder := NewDecoder([]byte{1})

	toDecode := [1]time.Duration{}
	err := decoder.Decode(&toDecode)

	assert.EqualError(t, err, "decode: unsupported type time.Duration")
}

func TestDecoder_Slice_Err(t *testing.T) {
	buf := []byte{}

	decoder := NewDecoder(buf)
	var s []string
	err := decoder.Decode(&s)
	assert.Equal(t, err, ErrVarIntBufferSize)

	buf = []byte{0x01}

	decoder = NewDecoder(buf)
	err = decoder.Decode(&s)
	assert.Equal(t, err, ErrVarIntBufferSize)
}

func TestDecoder_Int64(t *testing.T) {
	buf := []byte{
		0x91, 0x7d, 0xf3, 0xff, 0xff, 0xff, 0xff, 0xff, //-819823
		0xe3, 0x1c, 0x1, 0x00, 0x00, 0x00, 0x00, 0x00, //72931
	}

	d := NewDecoder(buf)

	n, err := d.ReadInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(-819823), n)
	assert.Equal(t, 8, d.remaining())

	n, err = d.ReadInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(72931), n)
	assert.Equal(t, 0, d.remaining())
}

type binaryTestStruct struct {
	F1  string
	F2  int16
	F3  uint16
	F4  int32
	F5  uint32
	F6  int64
	F7  uint64
	F8  float32
	F9  float64
	F10 []string
	F11 [2]string
	F12 byte
	F13 []byte
	F14 bool
	F15 Int64
	F16 Uint64
	F17 JSONFloat64
	F18 Uint128
	F19 Int128
	F20 Float128
	F21 Varuint32
	F22 Varint32
	F23 Bool
	F24 HexBytes
}

func TestDecoder_BinaryStruct(t *testing.T) {
	cnt, err := hex.DecodeString("03616263b5ff630019ffffffe703000051ccffffffffffff9f860100000000003d0ab9c15c8fc2f5285c0f4002036465660337383903666f6f03626172ff05010203040501e9ffffffffffffff17000000000000001f85eb51b81e09400a000000000000005200000000000000070000000000000003000000000000000a000000000000005200000000000000e707cd0f01050102030405")
	require.NoError(t, err)

	s := &binaryTestStruct{}
	decoder := NewDecoder(cnt)
	assert.NoError(t, decoder.Decode(s))

	assert.Equal(t, "abc", s.F1)
	assert.Equal(t, int16(-75), s.F2)
	assert.Equal(t, uint16(99), s.F3)
	assert.Equal(t, int32(-231), s.F4)
	assert.Equal(t, uint32(999), s.F5)
	assert.Equal(t, int64(-13231), s.F6)
	assert.Equal(t, uint64(99999), s.F7)
	assert.Equal(t, float32(-23.13), s.F8)
	assert.Equal(t, float64(3.92), s.F9)
	assert.Equal(t, []string{"def", "789"}, s.F10)
	assert.Equal(t, [2]string{"foo", "bar"}, s.F11)
	assert.Equal(t, uint8(0xff), s.F12)
	assert.Equal(t, []byte{1, 2, 3, 4, 5}, s.F13)
	assert.Equal(t, true, s.F14)
	assert.Equal(t, Int64(-23), s.F15)
	assert.Equal(t, Uint64(23), s.F16)
	assert.Equal(t, JSONFloat64(3.14), s.F17)
	assert.Equal(t, Uint128{
		Lo: 10,
		Hi: 82,
	}, s.F18)
	assert.Equal(t, Int128{
		Lo: 7,
		Hi: 3,
	}, s.F19)
	assert.Equal(t, Float128{
		Lo: 10,
		Hi: 82,
	}, s.F20)
	assert.Equal(t, Varuint32(999), s.F21)
	assert.Equal(t, Varint32(-999), s.F22)
	assert.Equal(t, Bool(true), s.F23)
	assert.Equal(t, HexBytes([]byte{1, 2, 3, 4, 5}), s.F24)
}

type binaryInvalidTestStruct struct {
	F1 time.Duration
}

func TestDecoder_BinaryStruct_Err(t *testing.T) {
	s := binaryInvalidTestStruct{}
	decoder := NewDecoder([]byte{})
	err := decoder.Decode(&s)
	assert.EqualError(t, err, "decode: unsupported type time.Duration")
}

func TestDecoder_Decode_No_Ptr(t *testing.T) {
	decoder := NewDecoder([]byte{})
	err := decoder.Decode(1)
	assert.EqualError(t, err, "can only decode to pointer type, got int")
}