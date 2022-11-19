package compiler

import (
	"fmt"
	"os"
)

const helpMsg = `Usage: bootstrap (OPTIONS) SOURCE FILE(S)
Source file must have a .joe extention to be compiled.

[-options]
			-o<file>				set the output object file
			-c						compile only and do not generate executable
			-a						enable aggressive error reporting
			-O						optimize executable
			-V<version>				set the application version
			-w			--warnings	disable all warnings
			-d			--debug		string debugging info
			-r			--release	generate a release build executable
			-u			--unsafe	allow unsafe code
						--strip		strip away metadata
						--objdmp	create dump file for generated assembly
						--werror	enable warnings as errors
			
			-v			--version	print compiler version and exit
			-h			--help		display this message and exit`

// help displays helpMsg.
func help() {
	fmt.Println(helpMsg)
}

// printErr prints error.
func printErr(msg string) {
	fmt.Println("joec: err: " + msg)
	os.Exit(1)
}

func printVersion()

func bootstrap() int {
	if len(os.Args) < 2 {
		help()
		return 1
	}

	// Initializes the error manager's errors
	InitializeErrors()

	var files []string

	for i := 0; i < len(os.Args); i++ {
	args_:
		if os.Args[i] == "-a" {
			COptionsAggressiveErrorReporting = true
		} else if os.Args[i] == "-c" {
			COptionsCompileOnly = true
		} else if os.Args[i] == "-o" {
			if i+1 >= len(os.Args) {
				printErr("output file required after option '-o'")
			} else {
				i++
				COptionsOutputFile = os.Args[i]
			}
		} else if os.Args[i] == "-v" || os.Args[i] == "--version" {
			printVersion()
			os.Exit(0)
		} else if os.Args[i] == "-O" {
			COptionsOptimize = true
		} else if os.Args[i] == "-h" || os.Args[i] == "--help" {
			help()
			os.Exit(0)
		} else if os.Args[i] == "-r" || os.Args[i] == "--release" {
			COptionsOptimize = true
			COptionsDebug = false
			COptionsStrip = true
		} else if os.Args[i] == "-d" || os.Args[i] == "--debug" {
			COptionsDebug = true
		} else if os.Args[i] == "--strip" {
			COptionsStrip = true
		} else if os.Args[i] == "--magic" { // Easter egg
			COptionsMagic = true
		} else if os.Args[i] == "-w" || os.Args[i] == "--warnings" {
			COptionsWarnings = false
		} else if os.Args[i] == "-V" {
			if i+1 >= len(os.Args) {
				printErr("file version required after option '-V'")
			} else {
				i++
				COptionsVersion = os.Args[i]
			}
		} else if os.Args[i] == "-u" || os.Args[i] == "--unsafe" {
			COptionsUnsafe = true
		} else if os.Args[i] == "--werror" {
			COptionsWarnings = true
			COptionsWErrors = true
		} else if os.Args[i] == "--objdmp" {
			COptionsObjDmp = true
		} else {
			for ok := true; ok; ok = i < len(os.Args) {
				if os.Args[i][0] == '-' {
					goto args_
				}
				i++
				files = append(files, os.Args[i])
			}
			break
		}
	}

	if len(files) == 0 {
		help()
		return 1
	}

	return 0
}
