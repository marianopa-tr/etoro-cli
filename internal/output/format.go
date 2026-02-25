package output

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Format int

const (
	Table Format = iota
	JSON
)

var currentFormat Format

func SetFormat(f Format) { currentFormat = f }
func GetFormat() Format  { return currentFormat }

func PrintJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

func PrintError(err error) {
	if currentFormat == JSON {
		fmt.Fprintf(os.Stdout, `{"error":%q}`+"\n", err.Error())
	} else {
		Errorf("%s", err)
	}
}

func Errorf(format string, args ...any) {
	c := color.New(color.FgRed, color.Bold)
	c.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
}

func Warnf(format string, args ...any) {
	c := color.New(color.FgYellow, color.Bold)
	c.Fprintf(os.Stderr, "Warning: "+format+"\n", args...)
}

func Successf(format string, args ...any) {
	c := color.New(color.FgGreen)
	c.Fprintf(os.Stderr, format+"\n", args...)
}

func Infof(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func Bold(s string) string {
	return color.New(color.Bold).Sprint(s)
}

func Green(s string) string {
	return color.GreenString(s)
}

func Red(s string) string {
	return color.RedString(s)
}

func Yellow(s string) string {
	return color.YellowString(s)
}

func Cyan(s string) string {
	return color.CyanString(s)
}

func FormatPnL(val float64) string {
	if val > 0 {
		return color.GreenString("+%.2f", val)
	} else if val < 0 {
		return color.RedString("%.2f", val)
	}
	return fmt.Sprintf("%.2f", val)
}

func FormatPercent(val float64) string {
	if val > 0 {
		return color.GreenString("+%.2f%%", val)
	} else if val < 0 {
		return color.RedString("%.2f%%", val)
	}
	return fmt.Sprintf("%.2f%%", val)
}

func FormatMoney(val float64) string {
	if val >= 1_000_000 {
		return fmt.Sprintf("$%.1fM", val/1_000_000)
	} else if val >= 1_000 {
		return fmt.Sprintf("$%.1fK", val/1_000)
	}
	return fmt.Sprintf("$%.2f", val)
}

func FormatBool(v bool, trueStr, falseStr string) string {
	if v {
		return color.GreenString(trueStr)
	}
	return color.RedString(falseStr)
}

func NewTable(headers ...string) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.Style().Format.Header = text.FormatDefault
	row := make(table.Row, len(headers))
	for i, h := range headers {
		row[i] = Bold(h)
	}
	t.AppendHeader(row)
	return t
}

func RenderTable(t table.Writer) {
	t.Render()
}

func NewDetailTable() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.Style().Options.SeparateRows = true
	t.Style().Format.Header = text.FormatDefault
	return t
}

func DetailRow(t table.Writer, label string, value any) {
	t.AppendRow(table.Row{Bold(label), fmt.Sprintf("%v", value)})
}

func Confirm(prompt string) bool {
	fmt.Fprintf(os.Stderr, "%s [y/N] ", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false
	}
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return answer == "y" || answer == "yes"
}

func ConfirmDanger(message string) bool {
	c := color.New(color.FgRed, color.Bold)
	c.Fprintln(os.Stderr, message)
	fmt.Fprintf(os.Stderr, "Type YES to confirm: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false
	}
	return strings.TrimSpace(scanner.Text()) == "YES"
}
