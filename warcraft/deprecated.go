package warcraft

import "bytes"

func ParseInteger(data []byte, integer interface{}) {
    ReadIntegerRequired(bytes.NewBuffer(data), integer)
}

func ParseString(data []byte) (str string, rawSize int) {
    raw := ParseRawString(data)
    str = string(raw) // UTF8
    rawSize = len(raw) + 1
    return
}

func ParseRawString(data []byte) []byte {
    return ReadNullTerminatedBytesRequired(bytes.NewBuffer(data))
}
