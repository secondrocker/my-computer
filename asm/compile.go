package asm

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func Compile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	reg, _ := regexp.Compile("//.*$")
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " ")
		line = reg.ReplaceAllString(line, "")
		fmt.Println(line)
	}
}
