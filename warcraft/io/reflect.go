package io

import "reflect"

func ReflectRead(reader Reader, v interface{}) error {
    return reflectRead(reader, reflect.ValueOf(v))
}

func reflectRead(reader Reader, r reflect.Value) (err error) {
    switch r.Kind() {
    case reflect.Ptr:
        err = reflectRead(reader, r.Elem())
        return
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
        reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        err = ReadInteger(reader, r.Addr().Interface())
        return
    case reflect.String:
        var str string
        str, err = ReadNullTerminatedString(reader)
        if err != nil {
            return
        }
        r.SetString(str)
    case reflect.Array:
        if r.Type().Elem().Kind() == reflect.Uint8 {
            _, err = reader.Read(r.Slice(0, r.Len()).Bytes())
            if err != nil {
                return
            }
        } else {
            for i := 0; i < r.Len(); i++ {
                err = reflectRead(reader, r.Index(i))
                if err != nil {
                    return err
                }
            }
        }
    default:
        for i := 0; i < r.NumField(); i++ {
            field := r.Field(i)
            switch r.Type().Field(i).Tag.Get("encoding") {
            case "nested":
                nestedReader := NewNestedReader(reader)
                err = reflectRead(nestedReader, field)
                if err != nil {
                    return
                }
                _, err = nestedReader.SkipAll()
            default:
                err = reflectRead(reader, field)
            }
            if err != nil {
                return
            }
        }
    }

    return
}

func ReflectWrite(writer Writer, v interface{}) error {
    return reflectWrite(writer, reflect.ValueOf(v))
}

func reflectWrite(writer Writer, r reflect.Value) (err error) {
    switch r.Kind() {
    case reflect.Ptr:
        return reflectWrite(writer, r.Elem())
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
        reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return WriteInteger(writer, r.Interface())
    case reflect.String:
        err = WriteNullTerminatedString(writer, r.String())
        if err != nil {
            return
        }
    case reflect.Array:
        // not addressable, Slice() should not be used
        for i := 0; i < r.Len(); i++ {
            err = reflectWrite(writer, r.Index(i))
            if err != nil {
                return err
            }
        }
    default:
        for i := 0; i < r.NumField(); i++ {
            field := r.Field(i)
            switch r.Type().Field(i).Tag.Get("encoding") {
            case "nested":
                nestedWriter := NewNestedWriter(writer)
                err = reflectWrite(nestedWriter, field)
                if err != nil {
                    return
                }
                err = nestedWriter.Close()
            default:
                err = reflectWrite(writer, field)
            }
            if err != nil {
                return
            }
        }
    }

    return
}
