package main

import (
	"shell-quote"
	"diff"
)

type Token = shell_quote.Token
type ControlOperator = shell_quote.ControlOperator
type ParseEntry = shell_quote.ParseEntry

func ParseShellCommand(cmd string, env map[string]string) []Token {
	return shell_quote.parse(cmd, env)
}

func QuoteShellArgs(args []string) string {
	return shell_quote.quote(args)
}

func CreateTwoFilesPatch(oldFileName, newFileName, oldStr, newStr, oldHeader, newHeader string, options map[string]interface{}) string {
	return diff.createTwoFilesPatch(oldFileName, newFileName, oldStr, newStr, oldHeader, newHeader, options)
}
