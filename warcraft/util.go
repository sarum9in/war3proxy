package warcraft

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "io"
)

type Reader interface {
    io.Reader
    io.ByteReader
}

type Writer interface {
    io.Writer
    io.ByteWriter
}

func ReadInteger(reader Reader, integer interface{}) {
    err := binary.Read(reader, binary.LittleEndian, integer)
    if err != nil {
        panic(err)
    }
}

func WriteInteger(writer Writer, integer interface{}) {
    err := binary.Write(writer, binary.LittleEndian, integer)
    if err != nil {
        panic(err)
    }
}

func ReadNullTerminatedString(reader Reader) string {
    raw := ReadNullTerminatedBytes(reader)
    return string(raw) // UTF8
}

func WriteNullTerminatedString(writer Writer, data string) {
    WriteNullTerminatedBytes(writer, []byte(data))
}

func ReadNullTerminatedBytes(reader Reader) []byte {
    var buffer bytes.Buffer

    var err error

    for c, err := reader.ReadByte(); err == nil; c, err = reader.ReadByte() {
        if c == 0 {
            break
        } else {
            buffer.WriteByte(byte(c))
        }
    }
    if err != nil && err != io.EOF {
        panic(err)
    }

    return buffer.Bytes()
}

func WriteNullTerminatedBytes(writer Writer, data []byte) {
    for _, c := range data {
        if c == 0 {
            panic(fmt.Errorf("Invalid data with null character"))
        }
    }

    _, err := writer.Write(data)
    if err != nil {
        panic(err)
    }

    err = writer.WriteByte(0)
    if err != nil {
        panic(err)
    }
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
