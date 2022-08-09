<head>
<h1 align=center>Franky - Shellcode Injection</h1>
</head>

![Franky Logo](images/franky.gif "Franky")

Franky is a shellcode generation and injection tool, designed to generate an executable payload to inject custom shellcode into another process.

---

## Features

* Inject shellcode into a remote process
* Generate shellcode from an executable to be embedded into payload
* Supports 32-bit or 64-bit architecture
* Configurable process to use for injection
* Builder program works on Windows, Linux (MacOS not tested)
---
## Installation

Clone the repository:<br>
```git clone https://github.com/misthi0s/Franky```

Change the working directory:<br>
```cd Franky```

Build the project with Go:<br>
```go build -o franky .```

Help information can be accessed by running the `-h` switch:<br>
```./franky -h```

---
## Usage

Franky supports the following commandline switches:
<table>
<tr>
<th>Switch</th>
<th>Description</th>
<th>Mandatory?</th>
</tr>
<tr>
<td>-path</td>
<td>Path to shellcode or executable file</td>
<td>Yes</td>
</tr>
<tr>
<td>-arch</td>
<td>Architecture of the output payload; 32 (default) or 64</td>
<td>No</td>
</tr>
<tr>
<td>-process</td>
<td>Full path of process to inject into (default: notepad.exe)</td>
<td>No</td>
</tr>
<tr>
<td>-exe</td>
<td>Generate shellcode from executable file to use as payload</td>
<td>No</td>
</tr>
</table>

---
## Examples

Generate a 32-bit payload to inject shellcode "shell.bin" into notepad.exe:<br>
``` ./franky -path shell.bin```

Generate a 64-bit payload to inject shellcode "shell.bin" into calc.exe:<br>
``` ./franky -arch 64 -path shell.bin -process C:\Windows\System32\calc.exe```

Convert executable "shell.exe" into shellcode and generate a 64-bit payload to inject into notepad.exe:<br>
``` ./franky -arch 64 -path shell.exe -exe```

---
## Additional Notes

* To avoid issues, architectures should all match! 64-bit shellcode should be generated as a 64-bit payload to inject into a 64-bit process.

* The "pkgs" folder (and contents) must exist in the same directory as the builder program. If the builder is moved, the "pkgs" folder must be moved with it.

* All payloads will be saved in the "output" folder in the same directory as the builder program.

---
## Credits

Special thanks go to:

* [TheWover](https://github.com/TheWover/donut) and [Binject](https://github.com/Binject/go-donut) for their excellent Donut and Go-Donut, respectively, projects. Without these, the executable to shellcode function would not be included in Franky.

---
## Issues

If you run into any issues with Franky, feel free to open an issue. It is likely that only issues related to the builder will be worked on, due to the near infinite number of issues that could stem from shellcode-specific problems.