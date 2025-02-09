package main

import (
	"encoding/csv"
	"fmt"
    "bufio"
	"io"
	"os"
	"strings"
    "flag"
)

var (
    print_headers bool = true
    expect_headers bool = true
    args []string
)

func init() {
    flag.BoolVar(&print_headers, "print_headers", false, "Print the CSV headers with the values")
    flag.BoolVar(&expect_headers, "expect_headers", true, "Expect the first line of input to be headers, not data")
    flag.Parse()

    if !expect_headers && print_headers {
        fmt.Fprintf(os.Stderr, "To print headers the data must have the first line be headers\n")
        os.Exit(1)
    }
    args = flag.Args()
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
	if len(args) > 1 {
		infile, err := os.Open(args[1])
        if err != nil {
            panic(err)
        }
        scanner = bufio.NewScanner(infile)
	}

    headers := make([]string, 0)
    linecount := 0
    for scanner.Scan() {
        line := scanner.Text()

        // Surely we can re-use the object...
        reader := csv.NewReader(strings.NewReader(line))
        record, err := reader.Read()
        if err != nil {
            panic(err)
        }
        for i, field := range record {
            if linecount == 0 && expect_headers {
                headers = append(headers, field)
            } else {
                if print_headers {
                    if len(headers) == 0 {
                        panic("assertion failure, headers should not be 0 here")
                    }
                    fmt.Println(headers[i] + ": " + field)
                } else {
                    fmt.Println(field)
                }
            }
        }
        fmt.Println("------------------------")
        linecount += 1
    }
    if err := scanner.Err(); err != nil {
        if err != io.EOF {
            panic(err)
        }
    }
}
