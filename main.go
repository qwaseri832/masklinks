package main

import "fmt"

func MaskLinks(input string) string {
    data := []byte(input)
    result := make([]byte, 0, len(data))
    prefix := []byte("http://")
    n := len(prefix)

    for i := 0; i < len(data); {
        if i+n <= len(data) && bytesEqual(data[i:i+n], prefix) {
            start := i
            i += n
            for i < len(data) && data[i] != ' ' {
                i++
            }
            for j := start; j < i; j++ {
                result = append(result, '*')
            }
        } else {
            result = append(result, data[i])
            i++
        }
    }
    return string(result)
}

func bytesEqual(a, b []byte) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}

func main() {
    input := "Visit http://example.com and also http://foo.bar/baz now!"
    data := []byte(input)
    result := make([]byte, 0, len(data))
    prefix := []byte("http://")
    n := len(prefix)

    for i := 0; i < len(data); {
        if i+n <= len(data) && bytesEqual(data[i:i+n], prefix) {
            start := i
            i += n
            for i < len(data) && data[i] != ' ' {
                i++
            }
            for j := start; j < i; j++ {
                result = append(result, '*')
            }
        } else {
            result = append(result, data[i])
            i++
        }
    }
    fmt.Printf("%s\n", result)
}