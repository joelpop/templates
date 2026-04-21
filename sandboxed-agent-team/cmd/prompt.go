package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var stdinReader = bufio.NewReader(os.Stdin)

// Prompt asks a question and returns the trimmed response.
func Prompt(question string) (string, error) {
	fmt.Print(question)
	line, err := stdinReader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}
	return strings.TrimSpace(line), nil
}

// PromptWithDefault prompts with a visible default. Empty input returns
// def.
func PromptWithDefault(question, def string) (string, error) {
	var suffix string
	if def != "" {
		suffix = fmt.Sprintf(" [%s]", def)
	}
	resp, err := Prompt(fmt.Sprintf("%s%s: ", question, suffix))
	if err != nil {
		return "", err
	}
	if resp == "" {
		return def, nil
	}
	return resp, nil
}

// PromptChoice presents a numbered list of options and returns the
// chosen value. Accepts either the option number (1-indexed) or an
// exact (case-insensitive) match of the option text. Empty input
// returns def.
func PromptChoice(question string, options []string, def string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("PromptChoice called with no options")
	}

	fmt.Println(question)
	for i, o := range options {
		suffix := ""
		if o == def {
			suffix = "  (default)"
		}
		fmt.Printf("  %d. %s%s\n", i+1, o, suffix)
	}

	prompt := "Enter number or name: "
	if def != "" {
		prompt = fmt.Sprintf("Enter number or name [%s]: ", def)
	}

	for {
		raw, err := Prompt(prompt)
		if err != nil {
			return "", err
		}
		if raw == "" {
			if def == "" {
				fmt.Println("  (no default — please pick an option)")
				continue
			}
			return def, nil
		}
		// Try as an index.
		if idx, err := strconv.Atoi(raw); err == nil && idx >= 1 && idx <= len(options) {
			return options[idx-1], nil
		}
		// Try as an exact match.
		for _, o := range options {
			if strings.EqualFold(raw, o) {
				return o, nil
			}
		}
		fmt.Printf("  (didn't match %q; try again)\n", raw)
	}
}

// PromptYesNo returns true for "yes" or "y", false for "no" or "n",
// and the default on empty input.
func PromptYesNo(question string, def bool) (bool, error) {
	var suffix string
	if def {
		suffix = " [Y/n]: "
	} else {
		suffix = " [y/N]: "
	}
	for {
		raw, err := Prompt(question + suffix)
		if err != nil {
			return false, err
		}
		if raw == "" {
			return def, nil
		}
		switch strings.ToLower(raw) {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		}
		fmt.Printf("  (answer y or n; got %q)\n", raw)
	}
}
