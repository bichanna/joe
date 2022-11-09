package compiler

import (
	"fmt"
	"os"
)

const helpMsg = `Usage: bootstrap (OPTIONS) SOURCE FILE(S)
Source file must have a .joe extention to be compiled.

[-options]
			-v			--version	print compiler version and exit
			-o<file>				set the output object file
			-c						compile only and do not generate executable
			-a						enable aggressive error reporting
			-s						string debugging info
			-O						optimize executable
			-w						disable all warnings
			-v<version>				set the application version
			-u			--unsafe	allow unsafe code
						--objdmp	create dump file for generated assembly
						--target	target the specified platform of Joe to run on
						--werror	enable warnings as errors
			-r			--release	generate a release build executable
			-h			--help		display this message and exit`

// help displays helpMsg.
func help() {
	fmt.Println(helpMsg)
}

func bootstrap() {
	if len(os.Args) < 2 {
		help()
		return
	}
}
