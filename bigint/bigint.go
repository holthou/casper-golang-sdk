package bigint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
)

type BigInt struct {
	*big.Int
}

// Big 返回底层 *big.Int 数据
func (i BigInt) Big() *big.Int {
	if i.Int == nil {
		return new(big.Int)
	}
	return i.Int
}

func (i BigInt) MarshalJSON() ([]byte, error) {
	if i.Int == nil {
		return []byte("null"), nil
	}
	return json.Marshal(i.Int.String())
}

var quote = []byte(`"`)
var null = []byte(`null`)
var hprx = []byte(`0x`) // `0x`

func (i *BigInt) UnmarshalJSON(text []byte) error {
	var ok bool
	if bytes.HasPrefix(text, quote) {
		n := text[1 : len(text)-1]
		if bytes.HasPrefix(n, hprx) {
			r := string(n[2:])
			if i.Int, ok = new(big.Int).SetString(r, 16); !ok {
				return fmt.Errorf(`bigint: can't convert "0x%s" to *big.Int`, r)
			}
			return nil
		}

		r := string(n)
		if i.Int, ok = new(big.Int).SetString(r, 10); !ok {
			return fmt.Errorf(`bigint: can't convert "%s" to *big.Int`, r)
		}
		return nil
	}

	if bytes.Equal(text, null) {
		i.Int = new(big.Int)
		return nil
	}

	r := string(text)
	if i.Int, ok = new(big.Int).SetString(r, 10); !ok {
		return fmt.Errorf("bigint: can't convert %s to *big.Int", r)
	}
	return nil
}
