package uutid

import (
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	now := time.Now().Truncate(0)

	uutid := New()
	if uutid == NilUUTID {
		t.Fatal("expected utid to not be nil")
	}

	uutidTime := uutid.Time()

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

	uutidTime := uutid.Time()

	diff := now.Sub(uutidTime)
	// UUTID are accurate down to the 100th of nanoseconds
	if diff > 100*time.Nanosecond || diff < 100*time.Nanosecond {
		t.Fatalf("expected UUTID time to be closer to now than 100 nanoseconds, got: %s", diff)
	}
}

func TestNewWithCustomTime(t *testing.T) {
	testTime := time.Date(2021, 1, 17, 1, 5, 10, 123456900, time.UTC).Truncate(0)

	uutid := NewWithTime(testTime)
	if uutid == NilUUTID {
		t.Fatal("expected utid to not be nil")
	}

	uutidTime := uutid.Time()

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

	unixTime := time.Unix(sec, 0)

	diff := expectedTime.Sub(unixTime)
	if diff != 0 {
		t.Fatalf("expected unix time to be the same, diff was: %s", diff)
	}
}

func TestUUID(t *testing.T) {
	sec := int64(7952935226)
	nsec := int64(782162000)

	testTime := time.Unix(sec, nsec)
	uutid := NewWithTime(testTime)
	if uutid == NilUUTID {
		t.Fatal("expected utid to not be nil")
	}

	if uutid.UUID()[0:19] != "da08293a-ba7b-4614-" {
		t.Fatalf(`expected UUID to be: "da08293a-ba7b-4614-", got: %s`, uutid.UUID()[0:19])
	}
}

func TestNewWithMathRand(t *testing.T) {
	rander = io.Reader(rand.New(rand.NewSource(int64(time.Now().UnixNano()))))
	SetRand(rander)

	now := time.Now().Truncate(0)

	uutid := New()
	if uutid == NilUUTID {
		t.Fatal("expected utid to not be nil")
	}

	uutidTime := uutid.Time()

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

	uutidTime := uutid.Time()

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

func TestFromBase32(t *testing.T) {
	// uutid := New()
	// expected := uutid.Base32()
	expected := "c8bdfc02794eheda39v63kw9" // Crockford

	toUUTID, err := FromBase32(expected)
	if err != nil {
		log.Fatal(err)
	}

	actual := toUUTID.Base32()

	if actual != expected {
		t.Fatalf("actual and expected base 32 doesn't match.\nexpected: %s, got: %s", expected, actual)
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

func BenchmarkNewMathRand(b *testing.B) {
	rander = io.Reader(rand.New(rand.NewSource(int64(time.Now().UnixNano()))))
	SetRand(rander)
	var uutid UUTID
	for i := 0; i < b.N; i++ {
		uutid = New()
	}
	_ = uutid
	// Reset the rander to default
	SetRand(nil)
}

func BenchmarkNewWithTimeMathRand(b *testing.B) {
	rander = io.Reader(rand.New(rand.NewSource(int64(time.Now().UnixNano()))))
	SetRand(rander)
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
		uutid = FromBytes(uutidBytes)
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

func BenchmarkFromBase32(b *testing.B) {
	uutid := New()
	base32 := uutid.Base32()
	var err error

	for i := 0; i < b.N; i++ {
		uutid, err = FromBase32(base32)
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

func BenchmarkBase32(b *testing.B) {
	uutid := New()
	var base32 string
	for i := 0; i < b.N; i++ {
		base32 = uutid.Base32()
	}
	_, _ = uutid, base32
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
