package main

import "fmt"

func MaskLinks(input string) string {
    data := []byte(input)
    result := make([]byte, 0, len(data))

    for i := 0; i < len(data); {
        if i+7 <= len(data) &&
            data[i] == 'h' && data[i+1] == 't' && data[i+2] == 't' &&
            data[i+3] == 'p' && data[i+4] == ':' &&
            data[i+5] == '/' && data[i+6] == '/' {

            start := i
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

func main() {
    input := "Visit http://example.com and http://test.com"
    masked := MaskLinks(input)
    fmt.Println(masked)
}
