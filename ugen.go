package ugen // import go.katupy.io/ugen

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type Generator struct {
	AnyCharacter bool
	Digit        bool
	Interval     string
	Base64       bool
	Hex          bool
	Ulid         bool
	UlidAsUuid   bool
	Uuid         bool
	Uuid7        bool
	Lower        bool
	Upper        bool
	Prefix       string
	Suffix       string
	Separator    string
	WithLineFeed bool
}

func (g *Generator) Gen(writer io.Writer, count, length int) error {
	anyCharacter := g.AnyCharacter || g.Base64 || g.Hex
	builder := new(strings.Builder)

	if g.Interval != "" {
		// Need to calculate first because of length.

		parts := strings.SplitN(g.Interval, ",", 2)

		first, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse number %q: %w", parts[0], err)
		}

		var second int64

		if len(parts) == 1 {
			second = first
			first = 0
		} else {
			second, err = strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse number %q: %w", parts[1], err)
			}
		}

		if first >= second {
			return fmt.Errorf("interval begin is greater than or equals to end")
		}

		max := second - first - 1

		nBig, err := rand.Int(rand.Reader, big.NewInt(max))
		if err != nil {
			return fmt.Errorf("failed to generate random number: %w", err)
		}

		builder.WriteString(strconv.FormatInt(first+nBig.Int64(), 10))
		length = builder.Len()
	}

	var genBuf []byte

	switch {
	case g.Ulid, g.Uuid, g.Uuid7:
		length = 16
	default:
		genBuf = make([]byte, length)
	}

	var encBuf []byte

	switch {
	case g.Base64:
		encBuf = make([]byte, base64.StdEncoding.EncodedLen(length))
	case g.Hex:
		encBuf = make([]byte, hex.EncodedLen(length))
	}

	for i := 0; i < count; i++ {
		switch {
		case g.Ulid:
			v, err := ulid.New(ulid.Timestamp(time.Now().UTC()), ulid.DefaultEntropy())
			if err != nil {
				return fmt.Errorf("failed to generate ULID: %w", err)
			}

			switch {
			case g.Base64, g.Hex:
				genBuf = v[:]
			default:
				if g.UlidAsUuid {
					genBuf = []byte(uuid.UUID(v).String())
				} else {
					genBuf = []byte(v.String())
				}
			}
		case g.Uuid:
			v, err := uuid.NewRandom()
			if err != nil {
				return fmt.Errorf("failed to generate UUIDv4: %w", err)
			}

			switch {
			case g.Base64, g.Hex:
				genBuf = v[:]
			default:
				genBuf = []byte(v.String())
			}
		case g.Uuid7:
			v, err := uuid.NewV7()
			if err != nil {
				return fmt.Errorf("failed to generate UUIDv7: %w", err)
			}

			switch {
			case g.Base64, g.Hex:
				genBuf = v[:]
			default:
				genBuf = []byte(v.String())
			}
		case anyCharacter:
			gen(length, genBuf, nil)
		case g.Digit:
			if length == 1 {
				gen(length, genBuf, digit)
			} else {
				gen(1, genBuf[:1], digit[1:])
				gen(length, genBuf[1:], digit)
			}
		default:
			gen(length, genBuf, digitLowerUpper)
		}

		builder.Reset()

		if i > 0 {
			builder.WriteString(g.Separator)
		}

		if g.Prefix != "" {
			builder.WriteString(g.Prefix)
		}

		switch {
		case g.Base64:
			base64.StdEncoding.Encode(encBuf, genBuf)
			builder.Write(encBuf)
		case g.Hex:
			hex.Encode(encBuf, genBuf)
			builder.Write(encBuf)
		default:
			builder.Write(genBuf)
		}

		if g.Suffix != "" {
			builder.WriteString(g.Suffix)
		}

		switch {
		case g.Lower:
			fmt.Fprint(writer, strings.ToLower(builder.String()))
		case g.Upper:
			fmt.Fprint(writer, strings.ToUpper(builder.String()))
		default:
			fmt.Fprint(writer, builder.String())
		}
	}

	if g.WithLineFeed {
		fmt.Fprint(writer, "\n")
	}

	return nil
}

// GenInterval works in a fundamentally different way than Gen since the string
// length can be variable, so we can't reuse the same (fixed-length) buffer
// for the results and possible encoding (base64, hex).
func (g *Generator) GenInterval(writer io.Writer, count int) error {
	parts := strings.SplitN(g.Interval, ",", 2)

	first, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse number %q: %w", parts[0], err)
	}

	var second int64

	if len(parts) == 1 {
		second = first
		first = 0
	} else {
		second, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse number %q: %w", parts[1], err)
		}
	}

	if first >= second {
		return fmt.Errorf("interval begin is greater than or equals to end")
	}

	max := second - first
	builder := new(strings.Builder)

	for i := 0; i < count; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(max))
		if err != nil {
			return fmt.Errorf("failed to generate random number: %w", err)
		}

		genBuf := []byte(strconv.FormatInt(first+n.Int64(), 10))
		length := len(genBuf)

		builder.Reset()

		if i > 0 {
			builder.WriteString(g.Separator)
		}

		if g.Prefix != "" {
			builder.WriteString(g.Prefix)
		}

		switch {
		case g.Base64:
			encBuf := make([]byte, base64.StdEncoding.EncodedLen(length))
			base64.StdEncoding.Encode(encBuf, genBuf)
			builder.Write(encBuf)
		case g.Hex:
			encBuf := make([]byte, hex.EncodedLen(length))
			hex.Encode(encBuf, genBuf)
			builder.Write(encBuf)
		default:
			builder.Write(genBuf)
		}

		fmt.Fprint(writer, builder.String())
	}

	if g.WithLineFeed {
		fmt.Fprint(writer, "\n")
	}

	return nil
}
