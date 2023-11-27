package command

import (
	"testing"

	"github.com/go-zoox/testify"
)

func TestIsCommand(t *testing.T) {
	testify.Assert(t, !IsCommand("/"))
	testify.Assert(t, IsCommand("/ls"))
	testify.Assert(t, IsCommand("/ask arg1"))
	testify.Assert(t, IsCommand("/ask arg1 arg2 arg3"))
	testify.Assert(t, IsCommand("/画图 清明上河图"), "\"/画图 清明上河图\" should be a command")
}

func TestGetCommandWithArgs(t *testing.T) {
	command, arg, err := ParseCommandWithArg("/ls")
	testify.Assert(t, err == nil)
	testify.Assert(t, command == "ls")
	testify.Assert(t, arg == "")

	command, arg, err = ParseCommandWithArg("/ls -al")
	testify.Assert(t, err == nil)
	testify.Assert(t, command == "ls")
	testify.Assert(t, arg == "-al")

	command, arg, err = ParseCommandWithArg("/ls -a -l -x")
	testify.Assert(t, err == nil)
	testify.Assert(t, command == "ls")
	testify.Assert(t, arg == "-a -l -x")

	command, arg, err = ParseCommandWithArg("/画图 清明上河图")
	testify.Assert(t, err == nil)
	testify.Assert(t, command == "画图")
	testify.Assert(t, arg == "清明上河图")
}
