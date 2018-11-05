package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func decodeKey(text string) (string, error) {
	var buf []byte
	r := bytes.NewBuffer([]byte(text))
	for {
		c, err := r.ReadByte()
		if err != nil {
			if err != io.EOF {
				return "", err
			}
			break
		}
		if c != '\\' {
			buf = append(buf, c)
			continue
		}
		n := r.Next(1)
		if len(n) == 0 {
			return "", io.EOF
		}
		// See: https://golang.org/ref/spec#Rune_literals
		if idx := strings.IndexByte(`abfnrtv\'"`, n[0]); idx != -1 {
			buf = append(buf, []byte("\a\b\f\n\r\t\v\\'\"")[idx])
			continue
		}

		switch n[0] {
		case 'x':
			fmt.Sscanf(string(r.Next(2)), "%02x", &c)
			buf = append(buf, c)
		default:
			n = append(n, r.Next(2)...)
			_, err := fmt.Sscanf(string(n), "%03o", &c)
			if err != nil {
				return "", err
			}
			buf = append(buf, c)
		}
	}
	return string(buf), nil
}

var indexTypeToString = map[byte]string{
	0:  "Null",
	1:  "Int64",
	2:  "Uint64",
	3:  "Float32",
	4:  "Float64",
	5:  "String",
	6:  "Bytes",
	7:  "BinaryLiteral",
	8:  "MysqlDecimal",
	9:  "MysqlDuration",
	10: "MysqlEnum",
	11: "MysqlBit",
	12: "MysqlSet",
	13: "MysqlTime",
	14: "Interface",
	15: "MinNotNull",
	16: "MaxValue",
	17: "Raw",
	18: "MysqlJSON",
}
