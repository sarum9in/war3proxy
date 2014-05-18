package io

import (
    "bytes"
    "io"
)

type CodeReader struct {
    reader Reader
    mask byte
    pos uint
}

func NewCodeReader(reader Reader) *CodeReader {
    return &CodeReader{ reader: reader }
}

func (codeReader *CodeReader) ReadByte() (c byte, err error) {
    c, err = codeReader.reader.ReadByte()
    if err != nil {
        return
    }
    if c == 0 {
        err = io.EOF
        return
    }

    pos := codeReader.pos
    codeReader.pos++
    switch pos {
    case 0:
        codeReader.mask = c
        c, err = codeReader.ReadByte()
        return
    case 7:
        codeReader.pos = 0
        fallthrough
    default:
        if codeReader.mask & byte(0x01 << pos) == 0 {
            c--
        }
    }
    return
}

func (codeReader *CodeReader) Read(data []byte) (n int, err error) {
    n = 0
    for i := 0; i < len(data); i++ {
        data[i], err = codeReader.ReadByte()
        if err != nil {
            return
        }
        n++
    }
    return
}

func (codeReader *CodeReader) SkipAll() (n int, err error) {
    n = 0
    for _, err = codeReader.ReadByte(); err == nil; _, err = codeReader.ReadByte() {
        n++
    }
    if err == io.EOF {
        err = nil
    }
    return
}

type CodeWriter struct {
    writer Writer
    pos uint
    data [8]byte
}

func NewCodeWriter(writer Writer) *CodeWriter {
    return &CodeWriter{ writer: writer }
}

func (codeWriter *CodeWriter) flush() (err error) {
    _, err = codeWriter.writer.Write(codeWriter.data[:codeWriter.pos])
    codeWriter.pos = 0
    return
}

func (codeWriter *CodeWriter) Close() (err error) {
    err = codeWriter.flush()
    if err != nil {
        return
    }

    err = codeWriter.writer.WriteByte(0)
    if err != nil {
        return
    }

    return
}

func (codeWriter *CodeWriter) WriteByte(c byte) (err error) {
    if codeWriter.pos == 0 {
        codeWriter.data[0] = 1
        codeWriter.pos++ // allocate mask when necessary
    }
    pos := codeWriter.pos
    codeWriter.pos++
    if pos == 7 {
        defer func() {
            if err == nil {
                err = codeWriter.flush()
            }
        }()
    }

    if c % 2 == 0 {
        codeWriter.data[pos] = c + 1
    } else {
        codeWriter.data[pos] = c
        codeWriter.data[0] |= byte(0x01 << pos)
    }

    return
}

func (codeWriter *CodeWriter) Write(data []byte) (n int, err error) {
    n = 0
    for i := 0; i < len(data); i++ {
        err = codeWriter.WriteByte(data[i])
        if err != nil {
            return
        }
        n++
    }
    return
}

func EncodeBytes(data []byte) []byte {
    var buffer bytes.Buffer
    writer := NewCodeWriter(&buffer)

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

func DecodeBytes(data []byte) []byte {
    reader := NewCodeReader(bytes.NewReader(data))
    var buffer bytes.Buffer

    _, err := io.Copy(&buffer, reader)
    if err != nil {
        panic(err)
    }

    return buffer.Bytes()
}
