package rendezvoushash

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"math"
)

type (
	Hashable interface {
		Hash(h hash.Hash)
	}

	NodeHasher[ID ordered] interface {
		Hash(nodeID ID, item any) uint64
	}

	DefaultNodeHasher[ID ordered] struct{}
)

func (DefaultNodeHasher[ID]) Hash(nodeID ID, item any) uint64 {
	h := sha256.New()
	writeHash(h, item)
	writeHash(h, nodeID)
	sum := h.Sum(nil)
	return binary.BigEndian.Uint64(sum[:8])
}

const (
	tagHashable byte = iota + 1
	tagString
	tagBytes
	tagBool
	tagInt64
	tagUint64
	tagFloat64
	tagStringer
	tagFallback
)

func writeHash(h hash.Hash, v any) {
	if hv, ok := v.(Hashable); ok {
		_, _ = h.Write([]byte{tagHashable})
		hv.Hash(h)
		return
	}

	switch x := v.(type) {
	case string:
		writeString(h, tagString, x)
	case []byte:
		writeBytes(h, tagBytes, x)
	case bool:
		_, _ = h.Write([]byte{tagBool})
		if x {
			_, _ = h.Write([]byte{1})
		} else {
			_, _ = h.Write([]byte{0})
		}
	case int:
		writeUint64(h, tagInt64, uint64(int64(x)))
	case int8:
		writeUint64(h, tagInt64, uint64(int64(x)))
	case int16:
		writeUint64(h, tagInt64, uint64(int64(x)))
	case int32:
		writeUint64(h, tagInt64, uint64(int64(x)))
	case int64:
		writeUint64(h, tagInt64, uint64(x))
	case uint:
		writeUint64(h, tagUint64, uint64(x))
	case uint8:
		writeUint64(h, tagUint64, uint64(x))
	case uint16:
		writeUint64(h, tagUint64, uint64(x))
	case uint32:
		writeUint64(h, tagUint64, uint64(x))
	case uint64:
		writeUint64(h, tagUint64, x)
	case uintptr:
		writeUint64(h, tagUint64, uint64(x))
	case float32:
		writeUint64(h, tagFloat64, uint64(math.Float32bits(x)))
	case float64:
		writeUint64(h, tagFloat64, math.Float64bits(x))
	case fmt.Stringer:
		writeString(h, tagStringer, x.String())
	default:
		writeString(h, tagFallback, fmt.Sprintf("%T:%v", v, v))
	}
}

func writeString(h hash.Hash, tag byte, s string) {
	_, _ = h.Write([]byte{tag})
	writeUint64(h, 0, uint64(len(s)))
	_, _ = h.Write([]byte(s))
}

func writeBytes(h hash.Hash, tag byte, b []byte) {
	_, _ = h.Write([]byte{tag})
	writeUint64(h, 0, uint64(len(b)))
	_, _ = h.Write(b)
}

func writeUint64(h hash.Hash, tag byte, v uint64) {
	if tag != 0 {
		_, _ = h.Write([]byte{tag})
	}
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], v)
	_, _ = h.Write(buf[:])
}
