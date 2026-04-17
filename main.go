package main

import "fmt"

func MaskLinks(input string) string {
    data := []byte(input)
    out := make([]byte, 0, len(data))
    
    p := []byte("http://")
    pl := len(p)

    for i := 0; i < len(data); {
        if i+pl <= len(data) {
            match := true
            for k := 0; k < pl; k++ {
                if data[i+k] != p[k] {
                    match = false
                    break
                }
            }
            if match {
                out = append(out, p...)
                i += pl
                for i < len(data) && data[i] != ' ' {
                    out = append(out, '*')
                    i++
                }
                continue
            }
        }
        out = append(out, data[i])
        i++
    }
    return string(out)
}

func main() {
    input := "Hello, its my page: http://localhost123.com See you"
    result := MaskLinks(input)
    fmt.Println(result)
}
