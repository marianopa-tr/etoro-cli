package shell

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
	"github.com/etoro/etoro-cli/internal/config"
	"github.com/etoro/etoro-cli/internal/output"
	"github.com/fatih/color"
)

type ExecuteFunc func(args []string) error

type Shell struct {
	execute  ExecuteFunc
	commands []string
	symbols  []string
}

func New(execute ExecuteFunc, commands []string, symbols []string) *Shell {
	return &Shell{
		execute:  execute,
		commands: commands,
		symbols:  symbols,
	}
}

func (s *Shell) Run() error {
	historyFile := filepath.Join(config.ConfigDir(), "history")
	_ = os.MkdirAll(config.ConfigDir(), 0o700)

	completer := readline.NewPrefixCompleter(s.buildCompletionTree()...)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          promptString(),
		HistoryFile:     historyFile,
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return fmt.Errorf("initializing shell: %w", err)
	}
	defer rl.Close()

	output.PrintBanner("Interactive Shell")

	dim := color.New(color.Faint)
	dim.Fprintln(os.Stderr, "  Type 'help' for commands, 'exit' to quit.")
	fmt.Fprintln(os.Stderr)

	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		}
		if err == io.EOF {
			fmt.Fprintln(os.Stderr, "\nGoodbye!")
			return nil
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if line == "exit" || line == "quit" {
			fmt.Fprintln(os.Stderr, "Goodbye!")
			return nil
		}

		if line == "help" {
			line = "--help"
		}

		args := strings.Fields(line)
		if err := s.execute(args); err != nil {
			color.New(color.FgRed).Fprintf(os.Stderr, "Error: %s\n", err)
		}
	}
}

func promptString() string {
	return color.New(color.FgGreen, color.Bold).Sprint("etoro") +
		color.New(color.FgWhite).Sprint("> ")
}

func (s *Shell) buildCompletionTree() []readline.PrefixCompleterInterface {
	items := make([]readline.PrefixCompleterInterface, 0, len(s.commands))
	for _, cmd := range s.commands {
		items = append(items, readline.PcItem(cmd))
	}
	items = append(items, readline.PcItem("exit"))
	items = append(items, readline.PcItem("quit"))
	items = append(items, readline.PcItem("help"))
	return items
}
