package warcraft

import (
    "fmt"
    "reflect"
)

func ReflectRead(reader Reader, v interface{}) error {
    return reflectRead(reader, reflect.ValueOf(v))
}

func reflectRead(reader Reader, r reflect.Value) error {
    if readFrom := r.MethodByName("ReadFrom"); readFrom.IsValid() {
        ret := readFrom.Call([]reflect.Value { reflect.ValueOf(reader) })
        if len(ret) != 1 {
            panic(fmt.Errorf("Invalid number of return values from ReadFrom(), should be 1"))
        }
        switch x := ret[0].Interface().(type) {
        case error:
            return x
        default:
            panic(fmt.Errorf("Invalid return value from ReadFrom(), should be error"))
        }
    } else {
        return defaultReadFrom(reader, r)
    }
}

func defaultReadFrom(reader Reader, r reflect.Value) error {
    switch r.Kind() {
    case reflect.Ptr:
        return defaultReadFrom(reader, r.Elem())
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
         reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return ReadInteger(reader, r.Addr().Interface())
    case reflect.String:
        str, err := ReadNullTerminatedString(reader)
        if err != nil {
            return err
        }
        r.SetString(str)
    case reflect.Array: // we use only byte arrays
        _, err := reader.Read(r.Slice(0, r.Len()).Bytes())
        if err != nil {
            return err
        }
    default:
        for i := 0; i < r.NumField(); i++ {
            field := r.Field(i)
            err := reflectRead(reader, field)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func ReflectWrite(writer Writer, v interface{}) error {
    return reflectWrite(writer, reflect.ValueOf(v))
}

func reflectWrite(writer Writer, r reflect.Value) error {
    if writeTo := r.MethodByName("WriteTo"); writeTo.IsValid() {
        ret := writeTo.Call([]reflect.Value { reflect.ValueOf(writer) })
        if len(ret) != 1 {
            panic(fmt.Errorf("Invalid number of return values from ReadFrom(), should be 1"))
        }
        switch x := ret[0].Interface().(type) {
        case error:
            return x
        default:
            panic(fmt.Errorf("Invalid return value from ReadFrom(), should be error"))
        }
    } else {
        return defaultWriteTo(writer, r)
    }
}

func defaultWriteTo(writer Writer, r reflect.Value) error {
    switch r.Kind() {
    case reflect.Ptr:
        return defaultWriteTo(writer, r.Elem())
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
         reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return WriteInteger(writer, r.Interface())
    case reflect.String:
        err := WriteNullTerminatedString(writer, r.String())
        if err != nil {
            return err
        }
    case reflect.Array: // we use only byte arrays
        _, err := writer.Write(r.Slice(0, r.Len()).Bytes())
        if err != nil {
            return err
        }
    default:
        for i := 0; i < r.NumField(); i++ {
            field := r.Field(i)
            err := reflectWrite(writer, field)
            if err != nil {
                return err
            }
        }
    }

    return nil
}
