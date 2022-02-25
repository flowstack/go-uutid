package uutid

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

type UUTID [16]byte

var (
	// version is the UUID version to use
	version = 4

	// math/rand is faster than crypto/rand, but not cryptographically secure
	// rander = io.Reader(rand.New(rand.NewSource(int64(time.Now().UnixNano()))))
	rander = rand.Reader

	// NilUUTID is an empty UUTID, all zeros
	NilUUTID UUTID
)

// SetRand sets the random number generator.
// Calling with nil will set the random number generator to the default (math/rand).
// For slower but cryptographically secure randomness, use rand.Reader from crypto/rand.
func SetRand(r io.Reader) {
	if r == nil {
		rander = rand.Reader
		return
	}
	rander = r
}

// SetVersion set the UUID version to use
func SetVersion(v int) error {
	if v < 0 || v > 9 {
		return errors.New("version must be a positive integer smaller than 10")
	}
	version = v
	return nil
}

// New return a UUTID that looks like a UUID but is not directly compatible with UUID.
// UUTID can be converted to any UUID type.
func New() UUTID {
	return NewWithTime(time.Now())
}

// NewWithTime is used by New which uses time.Now() as t
func NewWithTime(t time.Time) UUTID {
	var uutid UUTID

	sec := t.Unix()
	nsec := t.Nanosecond()

	// Shift left to get the most of the high part of the nanoseconds in the first 16bit
	// This is to make room for the version
	// The first 2 bits will never be used anyway as its > 999999999
	nsec = nsec << 2

	// Extract the first and highest part of the nanoseconds
	ns1 := nsec >> 16 & 0xffff

	// Extract the lowest part of the nanoseconds
	ns2 := nsec & 0xffff

	// Utilize the four zeros in the lowest bits
	ns2 = (ns2 >> 4) & 0x0fff

	// Set the version in the last part of the timestamp
	// ns2 |= 0x4000 // Version 4
	ns2 |= version << 12 // e.g. 0x4000

	// Write the timestamp and version to the uutid
	binary.BigEndian.PutUint32(uutid[0:4], uint32(sec))
	binary.BigEndian.PutUint16(uutid[4:6], uint16(ns1))
	binary.BigEndian.PutUint16(uutid[6:8], uint16(ns2))

	// Fill the rest of the uutid with randomness
	_, err := io.ReadFull(rander, uutid[8:])
	if err != nil {
		return NilUUTID
	}

	// Finally set the variant to 1 (big endianness)
	uutid[8] = (uutid[8] & 0x3f) | 0x80

	return uutid
}

// FromBytes converts a byte slice to a UUTID
func FromBytes(uutidSlice []byte) UUTID {
	if len(uutidSlice) != 16 {
		return NilUUTID
	}

	return UUTID{
		0:  uutidSlice[0],
		1:  uutidSlice[1],
		2:  uutidSlice[2],
		3:  uutidSlice[3],
		4:  uutidSlice[4],
		5:  uutidSlice[5],
		6:  uutidSlice[6],
		7:  uutidSlice[7],
		8:  uutidSlice[8],
		9:  uutidSlice[9],
		10: uutidSlice[10],
		11: uutidSlice[11],
		12: uutidSlice[12],
		13: uutidSlice[13],
		14: uutidSlice[14],
		15: uutidSlice[15],
	}
}

// FromBase64 returns uutid from a base 64 encoded uutid
func FromBase64(str string) (UUTID, error) {
	if len(str) != 22 {
		return UUTID{}, errors.New("unable to extract uutid from base64 string")
	}

	uutid := UUTID{}
	base64.RawURLEncoding.Decode(uutid[:], []byte(str[:]))

	return uutid, nil
}

// FromBase32 returns uutid from a base 32 encoded uutid
const crockfordAlphabet = "0123456789abcdefghjkmnpqrstvwxyz"

var base32Encoder = base32.NewEncoding(crockfordAlphabet).WithPadding(base32.NoPadding)

func FromBase32(str string) (UUTID, error) {
	if len(str) != 24 {
		return UUTID{}, errors.New("unable to extract uutid from base32 string")
	}

	uutid := UUTID{}
	base32Encoder.Decode(uutid[:], []byte(str[:]))

	return uutid, nil
}

// FromBase16 returns uutid from a base 16 encoded uutid
func FromBase16(base16 string) (UUTID, error) {
	if len(base16) != 32 {
		return UUTID{}, errors.New("unable to extract uutid from base16 string")
	}

	uutid := UUTID{}
	hex.Decode(uutid[:], []byte(base16[:]))

	return uutid, nil
}

// FromUUID returns uutid from a UUID
func FromUUID(uuid string) (UUTID, error) {
	if len(uuid) == 32 {
		return FromBase16(uuid)

	} else if len(uuid) == 36 {
		uutid := UUTID{}
		hex.Decode(uutid[:4], []byte(uuid[0:8]))
		hex.Decode(uutid[4:6], []byte(uuid[9:13]))
		hex.Decode(uutid[6:8], []byte(uuid[14:18]))
		hex.Decode(uutid[8:10], []byte(uuid[19:23]))
		hex.Decode(uutid[10:], []byte(uuid[24:]))
		return uutid, nil
	}

	return UUTID{}, errors.New("unable to extract uutid")
}

// String returns uutid as a hex encoded string
func (uutid UUTID) String() string {
	return uutid.Base16()
}

// Base32 returns uutid as a regular base 32 encoded string
func (uutid UUTID) Base64() string {
	var buf [22]byte
	base64.RawURLEncoding.Encode(buf[:], uutid[:])
	return string(buf[:])
}

// Base32 returns uutid as a regular base 32 encoded string
func (uutid UUTID) Base32() string {
	var buf [24]byte
	base32Encoder.Encode(buf[:], uutid[:])
	return string(buf[:])
}

// Base16 returns uutid as a regular base 16 encoded string
func (uutid UUTID) Base16() string {
	var buf [32]byte
	hex.Encode(buf[:], uutid[:])
	return string(buf[:])
}

// UUID returns uutid as a UUID with version set - default is 4
func (uutid UUTID) UUID() string {
	var buf [36]byte

	hex.Encode(buf[0:8], uutid[:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], uutid[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], uutid[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], uutid[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], uutid[10:])
	return string(buf[:])
}

// Bytes returns uutid as a byte slice
func (uutid UUTID) Bytes() []byte {
	return []byte{
		0:  uutid[0],
		1:  uutid[1],
		2:  uutid[2],
		3:  uutid[3],
		4:  uutid[4],
		5:  uutid[5],
		6:  uutid[6],
		7:  uutid[7],
		8:  uutid[8],
		9:  uutid[9],
		10: uutid[10],
		11: uutid[11],
		12: uutid[12],
		13: uutid[13],
		14: uutid[14],
		15: uutid[15],
	}
}

// Time returns the timestamp of the UUTID
func (uutid UUTID) Time() time.Time {
	if len(uutid) < 10 {
		return time.Time{}
	}

	sec := int64(binary.BigEndian.Uint32(uutid[0:4]))
	ns1 := int64(binary.BigEndian.Uint16(uutid[4:6]))
	ns2 := int64(binary.BigEndian.Uint16(uutid[6:8]))

	// Remove the version
	ns2 = ns2 & 0x0fff

	// Move the lower part of the nanoseconds back to it's original position
	ns2 = ns2 << 4

	// Merge the ns1 and ns2 part of the nanoseconds
	nsec := ns1<<16 | ns2

	// Move the entire nanoseconds part back into it's original position
	nsec = nsec >> 2

	return time.Unix(sec, nsec)
}
