package logger

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"time"
)

func Fatal(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(color.Output,"%s %s\n", header("fatal"), color.RedString(format, a...))
	os.Exit(1)
}

func Error(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(color.Output,"%s %s\n", header("error"), color.RedString(format, a...))
}

func Warning(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(color.Output,"%s %s\n", header("warning"), color.YellowString(format, a...))
}

func Info(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(color.Output, "%v %v\n", header("info"), color.CyanString(format, a...))
}

func Success(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(color.Output,"%s %s\n", header("success"), color.GreenString(format, a...))
}

func Log(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(color.Output,"%s %s\n", header("log"), fmt.Sprintf(format, a...))
}

func header(t string) string {
	return color.MagentaString("[%s][%s]", time.Now().Format("15:04:05"), t)
}