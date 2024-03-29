package uutid

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"log"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	now := time.Now().Truncate(0)

	uutid := New()
	if uutid == NilUUTID {
		t.Fatal("expected utid to not be nil")
	}

	uutidTime := uutid.Time().Truncate(0)

	diff := uutidTime.Sub(now)
	if diff < 0 {
		t.Fatalf("expected diff to be more or equal to 0, got: %s", diff)
	}
	if diff > 5*time.Microsecond {
		t.Fatalf("expected UUTID time to be closer to now than 5 microsecond, got: %s", diff)
	}
}

func TestNewWithTimeNow(t *testing.T) {
	now := time.Now().Truncate(0)

	uutid := NewWithTime(now)
	if uutid == NilUUTID {
		t.Fatal("expected utid to not be nil")
	}

	uutidTime := uutid.Time().Truncate(0)

	diff := now.Sub(uutidTime)
	// UUTID are accurate down to the 100th of nanoseconds
	if diff > 100*time.Nanosecond {
		t.Fatalf("expected UUTID time to be closer to now than 100 nanoseconds, got: %s", diff)
	}
}

func TestNewWithCustomTime(t *testing.T) {
	testTime := time.Date(2021, 1, 17, 1, 5, 10, 123456900, time.UTC).Truncate(0)

	uutid := NewWithTime(testTime)
	if uutid == NilUUTID {
		t.Fatal("expected utid to not be nil")
	}

	uutidTime := uutid.Time().Truncate(0)

	diff := testTime.Sub(uutidTime)
	if diff != 0 {
		t.Fatalf("expected UUTID time to be the same, diff was: %s", diff)
	}

	if uutid.UUID()[0:19] != "60038d46-1d6f-4361-" {
		t.Fatalf(`expected UUTID to start with: "60038d46-1d6f-4361-", got: %s`, uutid.UUID()[0:19])
	}
}

func TestUUIDUnixTime(t *testing.T) {
	testTime := time.Date(2021, 1, 17, 1, 5, 10, 123456900, time.UTC).Truncate(0)
	expectedTime := time.Date(2021, 1, 17, 1, 5, 10, 0, time.UTC).Truncate(0)

	uutid := NewWithTime(testTime)
	if uutid == NilUUTID {
		t.Fatal("expected utid to not be nil")
	}

	uuidTime := uutid.UUID()[0:8]
	rawTime, _ := hex.DecodeString(string(uuidTime))
	sec := int64(binary.BigEndian.Uint32(rawTime))

	unixTime := time.Unix(sec, 0).Truncate(0)

	diff := expectedTime.Sub(unixTime)
	if diff != 0 {
		t.Fatalf("expected unix time to be the same, diff was: %s", diff)
	}
}

func TestUUID(t *testing.T) {
	sec := int64(7952935226)
	nsec := int64(782162000)

	testTime := time.Unix(sec, nsec).Truncate(0)
	uutid := NewWithTime(testTime)
	if uutid == NilUUTID {
		t.Fatal("expected utid to not be nil")
	}

	if uutid.UUID()[0:19] != "da08293a-ba7b-4614-" {
		t.Fatalf(`expected UUID to be: "da08293a-ba7b-4614-", got: %s`, uutid.UUID()[0:19])
	}
}

func TestNewWithCryptoRand(t *testing.T) {
	SetRand(rand.Reader)

	now := time.Now().Truncate(0)

	uutid := New()
	if uutid == NilUUTID {
		t.Fatal("expected utid to not be nil")
	}

	uutidTime := uutid.Time().Truncate(0)

	diff := uutidTime.Sub(now)
	if diff < 0 {
		t.Fatalf("expected diff to be more or equal to 0, got: %s", diff)
	}
	if diff > 5*time.Microsecond {
		t.Fatalf("expected UUTID time to be closer to now than 5 microsecond, got: %s", diff)
	}

	// Reset the rander to default
	SetRand(nil)
}

func TestNewWithVersion(t *testing.T) {
	SetVersion(5)

	now := time.Now().Truncate(0)

	uutid := New()
	if uutid == NilUUTID {
		t.Fatal("expected utid to not be nil")
	}

	uutidTime := uutid.Time().Truncate(0)

	diff := uutidTime.Sub(now)
	if diff < 0 {
		t.Fatalf("expected diff to be more or equal to 0, got: %s", diff)
	}
	if diff > 5*time.Microsecond {
		t.Fatalf("expected UUTID time to be closer to now than 5 microsecond, got: %s", diff)
	}

	// Reset the version to 4
	SetVersion(4)
}

func TestFromBase64(t *testing.T) {
	// uutid := New()
	// expected := uutid.Base64()
	expected := "YhbQUTKtTQ6mYtDKTsbxcg"

	toUUTID, err := FromBase64(expected)
	if err != nil {
		log.Fatal(err)
	}

	actual := toUUTID.Base64()

	if actual != expected {
		t.Fatalf("actual and expected base 64 doesn't match.\nexpected: %s, got: %s", expected, actual)
	}
}

func TestFromBase16(t *testing.T) {
	// uutid := New()
	// expected := uutid.Base16()
	expected := "6216b0a7290a42a28f945b18df4e0537"

	toUUTID, err := FromBase16(expected)
	if err != nil {
		log.Fatal(err)
	}

	actual := toUUTID.Base16()

	if actual != expected {
		t.Fatalf("actual and expected base 16 doesn't match.\nexpected: %s, got: %s", expected, actual)
	}
}

func TestFromUUID(t *testing.T) {
	// uutid := New()
	// expected := uutid.UUID()
	expected := "6216b0a7-290a-42a2-8f94-5b18df4e0537"

	toUUTID, err := FromUUID(expected)
	if err != nil {
		log.Fatal(err)
	}

	actual := toUUTID.UUID()

	if actual != expected {
		t.Fatalf("actual and expected UUID doesn't match.\nexpected: %s, got: %s", expected, actual)
	}
}

func TestFromString(t *testing.T) {
	expected := New()

	str64 := expected.Base64()
	actual, err := FromString(str64)
	if err != nil {
		log.Fatal(err)
	}
	if !bytes.Equal(expected[:], actual[:]) {
		t.Fatalf("Base64: actual and expected UUTID doesn't match.\nexpected:\t%b\ngot:\t\t%b", expected[:], actual[:])
	}

	str16 := expected.Base16()
	actual, err = FromString(str16)
	if err != nil {
		log.Fatal(err)
	}
	if !bytes.Equal(expected[:], actual[:]) {
		t.Fatalf("Base16: actual and expected UUTID doesn't match.\nexpected:\t%b\ngot:\t\t%b", expected[:], actual[:])
	}

	strUUID := expected.UUID()
	actual, err = FromString(strUUID)
	if err != nil {
		log.Fatal(err)
	}
	if !bytes.Equal(expected[:], actual[:]) {
		t.Fatalf("UUID: actual and expected UUTID doesn't match.\nexpected:\t%b\ngot:\t\t%b", expected[:], actual[:])
	}

	strBinary := string(expected[:])
	actual, err = FromString(strBinary)
	if err != nil {
		log.Fatal(err)
	}
	if !bytes.Equal(expected[:], actual[:]) {
		t.Fatalf("binary: actual and expected UUTID doesn't match.\nexpected:\t%b\ngot:\t\t%b", expected[:], actual[:])
	}
}

func TestAllCombosOnSameUUTID(t *testing.T) {
	testTime := time.Date(2021, 1, 17, 1, 5, 10, 123456900, time.UTC).Truncate(0)
	testUUTID := NewWithTime(testTime)

	fromFuncs := map[string]func(string) (UUTID, error){
		"Base64": FromBase64,
		"Base16": FromBase16,
		"UUID":   FromUUID,
	}

	toFuncs := map[string]func() string{
		"Base64": testUUTID.Base64,
		"Base16": testUUTID.Base16,
		"UUID":   testUUTID.UUID,
	}

	for typ, toFunc := range toFuncs {
		toStr := toFunc()
		uutid, err := fromFuncs[typ](toStr)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(testUUTID[:], uutid[:]) {
			t.Fatalf("binary: actual and expected UUTID doesn't match.\nexpected:\t%b\ngot:\t\t%b", testUUTID[:], uutid[:])
		}

		uutid2, err := FromString(toStr)
		if !bytes.Equal(uutid[:], uutid2[:]) {
			t.Fatalf("binary: actual and expected UUTID doesn't match.\nexpected:\t%b\ngot:\t\t%b", uutid[:], uutid2[:])
		}
	}
}

func BenchmarkNew(b *testing.B) {
	var uutid UUTID
	for i := 0; i < b.N; i++ {
		uutid = New()
	}
	_ = uutid
}

func BenchmarkNewWithTime(b *testing.B) {
	var uutid UUTID
	now := time.Now()
	for i := 0; i < b.N; i++ {
		uutid = NewWithTime(now)
	}
	_ = uutid
}

func BenchmarkNewCryptoRand(b *testing.B) {
	SetRand(rand.Reader)
	var uutid UUTID
	for i := 0; i < b.N; i++ {
		uutid = New()
	}
	_ = uutid
	// Reset the rander to default
	SetRand(nil)
}

func BenchmarkNewWithTimeCryptoRand(b *testing.B) {
	SetRand(rand.Reader)
	var uutid UUTID
	now := time.Now()
	for i := 0; i < b.N; i++ {
		uutid = NewWithTime(now)
	}
	_ = uutid
	// Reset the rander to default
	SetRand(nil)
}

func BenchmarkFromBytes(b *testing.B) {
	uutid := New()
	uutidBytes := uutid.Bytes()

	for i := 0; i < b.N; i++ {
		uutid, _ = FromBytes(uutidBytes)
	}
	_ = uutid
}

func BenchmarkFromBase64(b *testing.B) {
	uutid := New()
	base64 := uutid.Base64()
	var err error

	for i := 0; i < b.N; i++ {
		uutid, err = FromBase64(base64)
	}
	_, _ = uutid, err
}

func BenchmarkFromBase16(b *testing.B) {
	uutid := New()
	base16 := uutid.Base16()
	var err error

	for i := 0; i < b.N; i++ {
		uutid, err = FromBase16(base16)
	}
	_, _ = uutid, err
}

func BenchmarkFromUUID(b *testing.B) {
	uutid := New()
	uuid := uutid.UUID()
	var err error

	for i := 0; i < b.N; i++ {
		uutid, err = FromUUID(uuid)
	}
	_, _ = uutid, err
}

func BenchmarkBase64(b *testing.B) {
	uutid := New()
	var base64 string
	for i := 0; i < b.N; i++ {
		base64 = uutid.Base64()
	}
	_, _ = uutid, base64
}

func BenchmarkBase16(b *testing.B) {
	uutid := New()
	var base16 string
	for i := 0; i < b.N; i++ {
		base16 = uutid.Base16()
	}
	_, _ = uutid, base16
}

func BenchmarkUUID(b *testing.B) {
	uutid := New()
	var uuid string
	for i := 0; i < b.N; i++ {
		uuid = uutid.UUID()
	}
	_, _ = uutid, uuid
}
