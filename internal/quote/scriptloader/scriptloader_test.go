package scriptloader

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"testing"
)

var forwardKeyTable [256]byte

func init() {
	for i := 0; i < 256; i++ {
		forwardKeyTable[inverseKeyTable[i]] = byte(i)
	}
}

func xorEncode(data []byte, xorA, xorB byte) []byte {
	out := make([]byte, len(data))
	for i, b := range data {
		out[i] = forwardKeyTable[b^xorB] ^ xorA
	}
	return out
}

func encodeTestPayload(plaintext []byte) []byte {
	pass1Encoded := xorEncode(plaintext, pass1XorA, pass1XorB)

	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	w.Write(pass1Encoded)
	w.Close()

	pass2Encoded := xorEncode(compressed.Bytes(), pass2XorA, pass2XorB)

	var header [16]byte
	copy(header[:4], "ONS2")
	binary.LittleEndian.PutUint32(header[4:8], uint32(len(pass2Encoded)))
	binary.LittleEndian.PutUint32(header[8:12], uint32(len(plaintext)))
	binary.LittleEndian.PutUint32(header[12:16], 110)

	var out bytes.Buffer
	out.Write(header[:])
	out.Write(pass2Encoded)
	return out.Bytes()
}

func TestDecode_ValidPayload(t *testing.T) {
	plaintext := []byte("preset_define 0,6,36,#FFFFFF\nnew_episode 1\nd `Hello world`[\\]")
	encoded := encodeTestPayload(plaintext)

	decoded, err := decode(encoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(decoded, plaintext) {
		t.Errorf("decoded mismatch:\ngot:  %q\nwant: %q", decoded, plaintext)
	}
}

func TestDecode_TooShort(t *testing.T) {
	_, err := decode([]byte("ONS2"))
	if err == nil {
		t.Fatal("expected error for data too short")
	}
}

func TestDecode_InvalidMagic(t *testing.T) {
	data := make([]byte, 16)
	copy(data[:4], "NOPE")

	_, err := decode(data)
	if err == nil {
		t.Fatal("expected error for invalid magic")
	}
}

func TestDecode_CorruptZlib(t *testing.T) {
	var header [16]byte
	copy(header[:4], "ONS2")
	binary.LittleEndian.PutUint32(header[4:8], 10)
	binary.LittleEndian.PutUint32(header[8:12], 100)
	binary.LittleEndian.PutUint32(header[12:16], 110)

	data := append(header[:], []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}...)

	_, err := decode(data)
	if err == nil {
		t.Fatal("expected error for corrupt zlib data")
	}
}

func TestDecode_EmptyPayload(t *testing.T) {
	var plaintext []byte
	encoded := encodeTestPayload(plaintext)

	decoded, err := decode(encoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(decoded) != 0 {
		t.Errorf("expected empty decoded output, got %d bytes", len(decoded))
	}
}
