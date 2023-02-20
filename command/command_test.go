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
}

func TestGetCommandWithArgs(t *testing.T) {
	command, arg, err := ParseCommandWithArg("/ls")
	testify.Assert(t, err == nil)
	testify.Assert(t, command == "ls")
	testify.Assert(t, len(arg) == 0)

	command, arg, err = ParseCommandWithArg("/ls -al")
	testify.Assert(t, err == nil)
	testify.Assert(t, command == "ls")
	testify.Assert(t, len(arg) == 1)
	testify.Assert(t, arg == "-al")

	command, arg, err = ParseCommandWithArg("/ls -a -l -x")
	testify.Assert(t, err == nil)
	testify.Assert(t, command == "ls")
	testify.Assert(t, len(arg) == 1)
	testify.Assert(t, arg == "-a -l -x")
}
