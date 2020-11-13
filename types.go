package bin

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type SafeString string

func (ss *SafeString) UnmarshalBinary(d *Decoder) error {
	s, e := d.SafeReadUTF8String()
	if e != nil {
		return e
	}

	*ss = SafeString(s)
	return nil
}

type Bool bool

func (b *Bool) UnmarshalJSON(data []byte) error {
	var num int
	err := json.Unmarshal(data, &num)
	if err == nil {
		*b = Bool(num != 0)
		return nil
	}

	var boolVal bool
	if err := json.Unmarshal(data, &boolVal); err != nil {
		return fmt.Errorf("couldn't unmarshal bool as int or true/false: %s", err)
	}

	*b = Bool(boolVal)
	return nil
}

func (b *Bool) UnmarshalBinary(decoder *Decoder) error {
	value, err := decoder.ReadBool()
	if err != nil {
		return err
	}

	*b = Bool(value)
	return nil
}

func (b Bool) MarshalBinary(encoder *Encoder) error {
	return encoder.writeBool(bool(b))
}

type HexBytes []byte

func (t HexBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(t))
}

func (t *HexBytes) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	*t, err = hex.DecodeString(s)
	return
}

func (t HexBytes) String() string {
	return hex.EncodeToString(t)
}

func (o *HexBytes) UnmarshalBinary(decoder *Decoder) error {
	value, err := decoder.ReadByteArray()
	if err != nil {
		return fmt.Errorf("hex bytes: %s", err)
	}

	*o = HexBytes(value)
	return nil
}

func (o HexBytes) MarshalBinary(encoder *Encoder) error {
	return encoder.writeByteArray([]byte(o), true)
}

type Varint16 int16

func (o *Varint16) UnmarshalBinary(decoder *Decoder) error {
	value, err := decoder.ReadVarint16()
	if err != nil {
		return fmt.Errorf("varint16: %s", err)
	}

	*o = Varint16(value)
	return nil
}

func (o Varint16) MarshalBinary(encoder *Encoder) error {
	return encoder.writeVarInt(int(o))
}

type Varuint16 uint16

func (o *Varuint16) UnmarshalBinary(decoder *Decoder) error {
	value, err := decoder.ReadUvarint16()
	if err != nil {
		return fmt.Errorf("varuint16: %s", err)
	}

	*o = Varuint16(value)
	return nil
}

func (o Varuint16) MarshalBinary(encoder *Encoder) error {
	return encoder.writeUVarInt(int(o))
}

type Varuint32 uint32

func (o *Varuint32) UnmarshalBinary(decoder *Decoder) error {
	value, err := decoder.ReadUvarint64()
	if err != nil {
		return fmt.Errorf("varuint32: %s", err)
	}

	*o = Varuint32(value)
	return nil
}

func (o Varuint32) MarshalBinary(encoder *Encoder) error {
	return encoder.writeUVarInt(int(o))
}

type Varint32 int32

func (o *Varint32) UnmarshalBinary(decoder *Decoder) error {
	value, err := decoder.ReadVarint32()
	if err != nil {
		return err
	}

	*o = Varint32(value)
	return nil
}

func (o Varint32) MarshalBinary(encoder *Encoder) error {
	return encoder.writeVarInt(int(o))
}

type JSONFloat64 float64

func (f *JSONFloat64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty value")
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}

		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}

		*f = JSONFloat64(val)

		return nil
	}

	var fl float64
	if err := json.Unmarshal(data, &fl); err != nil {
		return err
	}

	*f = JSONFloat64(fl)

	return nil
}

func (f *JSONFloat64) UnmarshalBinary(decoder *Decoder) error {
	value, err := decoder.ReadFloat64()
	if err != nil {
		return err
	}

	*f = JSONFloat64(value)
	return nil
}

func (f JSONFloat64) MarshalBinary(encoder *Encoder) error {
	return encoder.writeFloat64(float64(f))
}

type Int64 int64

func (i Int64) MarshalJSON() (data []byte, err error) {
	if i > 0xffffffff || i < -0xffffffff {
		encodedInt, err := json.Marshal(int64(i))
		if err != nil {
			return nil, err
		}
		data = append([]byte{'"'}, encodedInt...)
		data = append(data, '"')
		return data, nil
	}
	return json.Marshal(int64(i))
}

func (i *Int64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty value")
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}

		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}

		*i = Int64(val)

		return nil
	}

	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = Int64(v)

	return nil
}

func (i *Int64) UnmarshalBinary(decoder *Decoder) error {
	value, err := decoder.ReadInt64()
	if err != nil {
		return err
	}

	*i = Int64(value)
	return nil
}

func (i Int64) MarshalBinary(encoder *Encoder) error {
	return encoder.writeInt64(int64(i))
}

type Uint64 uint64

func (i Uint64) MarshalJSON() (data []byte, err error) {
	if i > 0xffffffff {
		encodedInt, err := json.Marshal(uint64(i))
		if err != nil {
			return nil, err
		}
		data = append([]byte{'"'}, encodedInt...)
		data = append(data, '"')
		return data, nil
	}
	return json.Marshal(uint64(i))
}

func (i *Uint64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty value")
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}

		val, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}

		*i = Uint64(val)

		return nil
	}

	var v uint64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = Uint64(v)

	return nil
}

func (i *Uint64) UnmarshalBinary(decoder *Decoder) error {
	value, err := decoder.ReadUint64()
	if err != nil {
		return err
	}

	*i = Uint64(value)
	return nil
}

func (i Uint64) MarshalBinary(encoder *Encoder) error {
	return encoder.writeUint64(uint64(i))
}

// uint128
type Uint128 struct {
	Lo uint64
	Hi uint64
}

func (i Uint128) BigInt() *big.Int {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[:], i.Hi)
	binary.BigEndian.PutUint64(buf[8:], i.Lo)
	value := (&big.Int{}).SetBytes(buf)
	return value
}

func (i Uint128) String() string {
	//Same for Int128, Float128
	number := make([]byte, 16)
	binary.LittleEndian.PutUint64(number[:], i.Lo)
	binary.LittleEndian.PutUint64(number[8:], i.Hi)
	return fmt.Sprintf("0x%s%s", hex.EncodeToString(number[:8]), hex.EncodeToString(number[8:]))
}

func (i Uint128) DecimalString() string {
	return i.BigInt().String()
}

func (i Uint128) MarshalJSON() (data []byte, err error) {
	return []byte(`"` + i.String() + `"`), nil
}

func (i *Uint128) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if !strings.HasPrefix(s, "0x") && !strings.HasPrefix(s, "0X") {
		return fmt.Errorf("int128 expects 0x prefix")
	}

	truncatedVal := s[2:]
	if len(truncatedVal) != 32 {
		return fmt.Errorf("int128 expects 32 characters after 0x, had %d", len(truncatedVal))
	}

	loHex := truncatedVal[:16]
	hiHex := truncatedVal[16:]

	lo, err := hex.DecodeString(loHex)
	if err != nil {
		return err
	}

	hi, err := hex.DecodeString(hiHex)
	if err != nil {
		return err
	}

	loUint := binary.LittleEndian.Uint64(lo)
	hiUint := binary.LittleEndian.Uint64(hi)

	i.Lo = loUint
	i.Hi = hiUint

	return nil
}

func (i *Uint128) UnmarshalBinary(decoder *Decoder) error {
	value, err := decoder.ReadUint128("uint128")
	if err != nil {
		return err
	}

	*i = value
	return nil
}

func (i Uint128) MarshalBinary(encoder *Encoder) error {
	return encoder.writeUint128(i)
}

// Int128
type Int128 Uint128

func (i Int128) BigInt() *big.Int {
	comp := byte(0x80)
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[:], i.Hi)
	binary.BigEndian.PutUint64(buf[8:], i.Lo)

	var value *big.Int
	if (buf[0] & comp) == comp {
		buf = twosComplement(buf)
		value = (&big.Int{}).SetBytes(buf)
		value = value.Neg(value)
	} else {
		value = (&big.Int{}).SetBytes(buf)
	}
	return value
}

func (i Int128) String() string {
	return Uint128(i).String()
}

func (i Int128) DecimalString() string {
	return i.BigInt().String()
}

func (i Int128) MarshalJSON() (data []byte, err error) {
	return []byte(`"` + Uint128(i).String() + `"`), nil
}

func (i *Int128) UnmarshalJSON(data []byte) error {
	var el Uint128
	if err := json.Unmarshal(data, &el); err != nil {
		return err
	}

	out := Int128(el)
	*i = out

	return nil
}

func (i *Int128) UnmarshalBinary(decoder *Decoder) error {
	value, err := decoder.ReadInt128()
	if err != nil {
		return err
	}

	*i = value
	return nil
}

func (i Int128) MarshalBinary(encoder *Encoder) error {
	return encoder.writeInt128(i)
}

type Float128 Uint128

func (i Float128) MarshalJSON() (data []byte, err error) {
	return []byte(`"` + Uint128(i).String() + `"`), nil
}

func (i *Float128) UnmarshalJSON(data []byte) error {
	var el Uint128
	if err := json.Unmarshal(data, &el); err != nil {
		return err
	}

	out := Float128(el)
	*i = out

	return nil
}

func (i *Float128) UnmarshalBinary(decoder *Decoder) error {
	value, err := decoder.ReadFloat128()
	if err != nil {
		return err
	}

	*i = Float128(value)
	return nil
}

func (i Float128) MarshalBinary(encoder *Encoder) error {
	return encoder.writeUint128(Uint128(i))
}

// Blob

// Blob is base64 encoded data
// https://github.com/EOSIO/fc/blob/0e74738e938c2fe0f36c5238dbc549665ddaef82/include/fc/variant.hpp#L47
type Blob string

// Data returns decoded base64 data
func (b Blob) Data() ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(b))
}

// String returns the blob as a string
func (b Blob) String() string {
	return string(b)
}

func twosComplement(v []byte) []byte {
	buf := make([]byte, len(v))
	for i, b := range v {
		buf[i] = b ^ byte(0xff)
	}
	one := big.NewInt(1)
	value := (&big.Int{}).SetBytes(buf)
	return value.Add(value, one).Bytes()
}
