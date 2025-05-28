package main

import "os"

func main() {
	os.Exit(1) // want `прямой вызов os.Exit\(\) в функции main\(\) запрещен`
}
