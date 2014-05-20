package io

import (
    "bytes"
    "io"
)

type NestedReader struct {
    reader Reader
    mask   byte
    pos    uint
}

func NewNestedReader(reader Reader) *NestedReader {
    return &NestedReader{reader: reader}
}

func (nestedReader *NestedReader) ReadByte() (c byte, err error) {
    c, err = nestedReader.reader.ReadByte()
    if err != nil {
        return
    }
    if c == 0 {
        err = io.EOF
        return
    }

    pos := nestedReader.pos
    nestedReader.pos++
    switch pos {
    case 0:
        nestedReader.mask = c
        c, err = nestedReader.ReadByte()
        return
    case 7:
        nestedReader.pos = 0
        fallthrough
    default:
        if nestedReader.mask&byte(0x01<<pos) == 0 {
            c--
        }
    }
    return
}

func (nestedReader *NestedReader) Read(data []byte) (n int, err error) {
    n = 0
    for i := 0; i < len(data); i++ {
        data[i], err = nestedReader.ReadByte()
        if err != nil {
            return
        }
        n++
    }
    return
}

func (nestedReader *NestedReader) SkipAll() (n int, err error) {
    n = 0
    for _, err = nestedReader.ReadByte(); err == nil; _, err = nestedReader.ReadByte() {
        n++
    }
    if err == io.EOF {
        err = nil
    }
    return
}

type NestedWriter struct {
    writer Writer
    pos    uint
    data   [8]byte
}

func NewNestedWriter(writer Writer) *NestedWriter {
    return &NestedWriter{writer: writer}
}

func (nestedWriter *NestedWriter) flush() (err error) {
    _, err = nestedWriter.writer.Write(nestedWriter.data[:nestedWriter.pos])
    nestedWriter.pos = 0
    return
}

func (nestedWriter *NestedWriter) Close() (err error) {
    err = nestedWriter.flush()
    if err != nil {
        return
    }

    err = nestedWriter.writer.WriteByte(0)
    if err != nil {
        return
    }

    return
}

func (nestedWriter *NestedWriter) WriteByte(c byte) (err error) {
    if nestedWriter.pos == 0 {
        nestedWriter.data[0] = 1
        nestedWriter.pos++ // allocate mask when necessary
    }
    pos := nestedWriter.pos
    nestedWriter.pos++
    if pos == 7 {
        defer func() {
            if err == nil {
                err = nestedWriter.flush()
            }
        }()
    }

    if c%2 == 0 {
        nestedWriter.data[pos] = c + 1
    } else {
        nestedWriter.data[pos] = c
        nestedWriter.data[0] |= byte(0x01 << pos)
    }

    return
}

func (nestedWriter *NestedWriter) Write(data []byte) (n int, err error) {
    n = 0
    for i := 0; i < len(data); i++ {
        err = nestedWriter.WriteByte(data[i])
        if err != nil {
            return
        }
        n++
    }
    return
}

func EncodeNestedBytes(data []byte) []byte {
    var buffer bytes.Buffer
    writer := NewNestedWriter(&buffer)

    _, err := writer.Write(data)
    if err != nil {
        panic(err)
    }

    err = writer.Close()
    if err != nil {
        panic(err)
    }

    return buffer.Bytes()
}

func DecodeNestedBytes(data []byte) []byte {
    reader := NewNestedReader(bytes.NewReader(data))
    var buffer bytes.Buffer

    _, err := io.Copy(&buffer, reader)
    if err != nil {
        panic(err)
    }

    return buffer.Bytes()
}
