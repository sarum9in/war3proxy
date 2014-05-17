package warcraft

import (
    "bytes"
    "encoding/binary"
)

func ParseInteger(data []byte, integer interface{}) {
    buffer := bytes.NewBuffer(data)
    err := binary.Read(buffer, binary.LittleEndian, integer)
    if err != nil {
        panic(err)
    }
}

func ParseString(data []byte) (str string, rawSize int) {
    raw := ParseRawString(data)
    str = string(raw) // UTF8
    rawSize = len(raw) + 1
    return
}

func ParseRawString(data []byte) []byte {
    var output bytes.Buffer
    for _, c := range data {
        if c == 0 {
            break
        } else {
            output.WriteByte(byte(c))
        }
    }
    return output.Bytes()
}

func EncodeBytes(data []byte) []byte {
    var mask byte = 1
    groups := (len(data) + 6) / 7
    result := make([]byte, len(data) + groups)

    for pos, c := range data {
        // for each 7 bytes from data save [mask, c0, c1, ..., c6]
        rgroup := uint(pos / 7)
        rpos := uint(pos % 7) + 1
        dst := 8 * rgroup + rpos
        if c % 2 == 0 {
            result[dst] = c + 1
        } else {
            result[dst] = c
            mask |= byte(0x1 << rpos)
        }
        if rpos == 7 || pos + 1 == len(data) {
            result[8 * rgroup] = mask
            mask = 1
        }
    }

    return result
}

func DecodeBytes(data []byte) []byte {
    var buffer bytes.Buffer
    var mask byte = 0
    for pos, c := range data {
        rpos := uint(pos % 8)
        if rpos == 0 {
            mask = c
        } else {
            if mask & byte(0x1 << rpos) == 0 {
                buffer.WriteByte(c - 1)
            } else {
                buffer.WriteByte(c)
            }
        }
    }
    return buffer.Bytes()
}
