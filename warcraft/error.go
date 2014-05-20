package warcraft

import "fmt"

type ParseError struct {
    Name string
    Err  error
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("Unable to parse %s: %v", e.Name, e.Err)
}
