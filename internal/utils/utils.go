package utils

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"slices"
)

func OneOf[T comparable](value T, values ...T) bool {
	return slices.Index(values, value) > -1
}

func Ternary[T any](condition bool, first T, second T) T {
	if condition {
		return first
	}

	return second
}

func Second[T any](_ any, second T, _ ...any) T {
	return second
}

func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}

func Map[I any, O any](elements []I, transform func(I) O) []O {
	result := make([]O, len(elements))

	for index, element := range elements {
		result[index] = transform(element)
	}

	return result
}

// calculateHash computes the SHA-256 hash of the input byte slice.
// It returns the hash as a hex-encoded string.
func CalculateHash(data []byte) (hash string, err error) {
	hasher := sha256.New()

	if err = Second(hasher.Write(data)); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// The `ConstantTimeCompare` function is used here to securely compare two hash values.
// It prevents timing-based attacks by ensuring that the comparison takes the same
// amount of time, regardless of whether the values match or not.
// If the hashes are not equal, an error is returned.
func CompareHashes(hash1, hash2 string) error {
	if subtle.ConstantTimeCompare([]byte(hash1), []byte(hash2)) != 1 {
		return fmt.Errorf("the hashes are not equal. needed %s, actual %s", hash1, hash2)
	}

	return nil
}

// uint64ArrayToBytes converts a slice of uint64 integers to a slice of bytes.
// This function is typically used to convert a slice of packed data into bytes,
// which can then be stored in linear memory.
func Uint64ArrayToBytes(data []uint64) []byte {
	// Calculate the total number of bytes required to represent all the uint64
	// integers in the data slice. Since each uint64 integer is 8 bytes long,
	// we multiply the number of uint64 integers by 8 to get the total number of bytes.
	size := len(data) * 8

	result := make([]byte, size)

	for i, d := range data {
		// Convert d to its little-endian byte representation and store it in the
		// result slice. The binary.LittleEndian.PutUint64 function takes a slice
		// of bytes and a uint64 integer, and writes the uint64 integer into the slice
		// of bytes in little-endian order.
		// The result[i<<3:] slice expression ensures that each uint64 integer is
		// written to the correct position in the result slice.
		// i<<3 is equivalent to i*8, but using bit shifting (<<3) is slightly more
		// efficient than multiplication.
		binary.LittleEndian.PutUint64(result[i<<3:], d)
	}

	// Return the result slice of bytes.
	return result
}

// BytesToUint64Array converts a slice of bytes to a slice of uint64 integers.
// This function is typically used to convert a slice of bytes (which might have
// been stored in linear memory) back into packed data in the form of uint64 integers.
func BytesToUint64Array(data []byte) []uint64 {
	// Calculate the number of uint64 integers that can be extracted from the slice of bytes.
	// Since each uint64 integer is represented by 8 bytes, we divide the length of the
	// byte slice by 8 to determine the number of uint64 integers.
	size := len(data) / 8

	result := make([]uint64, size)

	for i := 0; i < size; i++ {
		// Convert a slice of bytes starting at position i*8 into its uint64 representation
		// in little-endian order. The binary.LittleEndian.Uint64 function takes a slice
		// of bytes and returns its representation as a uint64 integer.
		// The data[i<<3:] slice expression ensures that we're reading the correct bytes
		// for each uint64 integer. As before, i<<3 is a more efficient way of computing i*8.
		result[i] = binary.LittleEndian.Uint64(data[i<<3:])
	}

	// Return the result slice of uint64 integers.
	return result
}
