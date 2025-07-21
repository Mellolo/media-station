package util

import (
	"fmt"
	"strings"
)

func FormatErrorLog(uuid string, message string, stack ...string) string {
	var builder strings.Builder
	builder.WriteString("-----------------------------------------error----------------------------------------------\n")
	builder.WriteString(fmt.Sprintf("[uuid] %s\n", uuid))
	builder.WriteString(fmt.Sprintf("[message] %s\n", message))
	if len(stack) > 0 && stack[0] != "" {
		builder.WriteString(fmt.Sprintf("[stack]\n%s\n", stack[0]))
	}
	builder.WriteString("-----------------------------------------error----------------------------------------------")
	return builder.String()
}
