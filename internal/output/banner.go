package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
)

var bannerLines = []string{
	"╔██████╗   ╔████████╗   ╔██████╗   ╔██████╗   ╔██████╗",
	"██╔═══██╗  ╚═══██╔══╝  ██╔═══██╗  ██╔═══██╗  ██╔═══██╗",
	"████████║      ██║     ██║   ██║  ██║   ╚═╝  ██║   ██║",
	"██╔═════╝      ██║     ██║   ██║  ██║        ██║   ██║",
	"╚███████╗      ██║     ╚██████╔╝  ██║        ╚██████╔╝",
	" ╚══════╝      ╚═╝      ╚═════╝   ╚═╝         ╚═════╝",
}

func PrintBanner(subtitle string) {
	green := color.New(color.FgGreen, color.Bold)
	indent := "   "

	fmt.Fprintln(os.Stderr)

	bannerWidth := 0
	for _, line := range bannerLines {
		if w := runewidth.StringWidth(line); w > bannerWidth {
			bannerWidth = w
		}
		green.Fprintln(os.Stderr, indent+line)
	}

	if subtitle != "" {
		fmt.Fprintln(os.Stderr)
		dim := color.New(color.FgWhite)
		pad := (runewidth.StringWidth(indent) + bannerWidth - runewidth.StringWidth(subtitle)) / 2
		if pad < 3 {
			pad = 3
		}
		dim.Fprintf(os.Stderr, "%s%s\n", strings.Repeat(" ", pad), subtitle)
	}

	fmt.Fprintln(os.Stderr)
}

func PrintNoticeBox(message string) {
	displayWidth := runewidth.StringWidth(message)
	innerWidth := displayWidth + 6
	if innerWidth < 52 {
		innerWidth = 52
	}

	padTotal := innerWidth - displayWidth
	padLeft := padTotal / 2
	padRight := padTotal - padLeft

	top := "╭" + strings.Repeat("─", innerWidth) + "╮"
	mid := "│" + strings.Repeat(" ", padLeft) + message + strings.Repeat(" ", padRight) + "│"
	bot := "╰" + strings.Repeat("─", innerWidth) + "╯"

	c := color.New(color.FgCyan)
	c.Fprintln(os.Stderr, "   "+top)
	c.Fprintln(os.Stderr, "   "+mid)
	c.Fprintln(os.Stderr, "   "+bot)
	fmt.Fprintln(os.Stderr)
}

func PrintStepHeader(step, total int, name string) {
	fmt.Fprintln(os.Stderr)
	dim := color.New(color.Faint)
	bold := color.New(color.Bold, color.FgWhite)

	dim.Fprintf(os.Stderr, "  [%d/%d] ", step, total)
	bold.Fprintln(os.Stderr, name)

	dim.Fprintln(os.Stderr, "  "+strings.Repeat("─", 50))
	fmt.Fprintln(os.Stderr)
}

func PrintSavedBox(path string) {
	green := color.New(color.FgGreen, color.Bold)
	dim := color.New(color.Faint)

	fmt.Fprintln(os.Stderr)
	green.Fprintln(os.Stderr, "  ✔  Configuration saved!")
	fmt.Fprintln(os.Stderr)
	dim.Fprintf(os.Stderr, "     %s\n", path)
	fmt.Fprintln(os.Stderr)
}
