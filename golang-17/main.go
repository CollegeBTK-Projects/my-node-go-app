package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

const version = "0.1"

var (
	verbose      bool
	reportsTypes = []string{"html", "pdf", "json", "csv"}
)

func logf(format string, args ...any) {
	if verbose {
		fmt.Printf("[verbose] "+format+"\n", args...)
	}
}

func newConvertCmd() *cobra.Command {
	var input, output, format string
	var force, dryRun bool

	cmd := &cobra.Command{
		Use:   "convert",
		Short: "конвертировать файл",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := os.ReadFile(input)
			if err != nil {
				return fmt.Errorf("Ошибка: входной файл %q не найден", input)
			}
			dot := strings.LastIndex(input, ".")

			if output == "" {
				output = input[:dot] + "." + format
			}

			if dryRun {
				fmt.Printf("Файл %s будет конвертирован в %s\n", input, format)
				fmt.Printf("Результат: %s\n", output)
				return nil
			}

			logf("Чтение файла: %s", input)
			if _, err := os.Stat(output); err == nil && !force {
				return fmt.Errorf("Ошибка: файл %q уже существует. Используйте --force для перезаписи", output)
			}
			if err := os.WriteFile(output, data, 0o644); err != nil {
				return fmt.Errorf("Ошибка: %w", err)
			}
			if err = os.Remove(input); err != nil {
				return fmt.Errorf("Ошибка: %w", err)
			}
			fmt.Println("Файл конвертирован:", output)
			return nil
		},
	}

	cmd.Flags().StringVarP(&input, "input", "i", "", "входной файл")
	cmd.Flags().StringVarP(&output, "output", "o", "", "путь для сохранения")
	cmd.Flags().StringVarP(&format, "format", "F", "json", "целевой формат")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "перезаписать существующий файл")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "показать что будет сделано без выполнения")
	_ = cmd.MarkFlagRequired("input")
	return cmd
}

func newInteractiveCmd() *cobra.Command {
	var n, t, l string
	var git bool
	var noInteractive bool

	answers := struct {
		Project  string
		Type     string
		Language string
		Git      bool
	}{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "вызвать меню",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			nameSet := cmd.Flags().Changed("name")
			typeSet := cmd.Flags().Changed("type")

			info, err := os.Stdin.Stat()
			if err != nil {
				return err
			}

			isPipe := info.Mode()&os.ModeCharDevice == 0

			if isPipe {
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					answers.Project = scanner.Text()
				}
				if err := scanner.Err(); err != nil {
					return err
				}
			}

			qs := []*survey.Question{
				{
					Name: "Project",
					Prompt: &survey.Input{
						Message: "Название проекта: ",
					},
				},
				{
					Name: "Type",
					Prompt: &survey.Select{
						Message: "Выбери тип проекта: ",
						Options: []string{"Микросервис", "CLI-утилита", "Монолит"},
						Default: "Микросервис",
					},
				},
				{
					Name: "Language",
					Prompt: &survey.Select{
						Message: "Выберите язык проекта: ",
						Options: []string{"Go", "Rust", "Python"},
						Default: "Go",
					},
				},
				{
					Name: "Git",
					Prompt: &survey.Confirm{
						Message: "Использовать Git?",
						Default: false,
					},
				},
			}

			if !noInteractive && !isPipe && !nameSet && !typeSet {
				if err := survey.Ask(qs, &answers); err != nil {
					return fmt.Errorf("ошибка в вопросах: %w", err)
				}
			} else {
				if answers.Project == "" {
					answers.Project = n
				}
				answers.Type = t
				answers.Language = l
				answers.Git = git
			}

			if answers.Project == "" {
				return fmt.Errorf("недостаточно данных. Укажите --name")
			}

			if answers.Type == "" {
				return fmt.Errorf("недостаточно данных. Укажите --type")
			}

			if err := os.Mkdir(answers.Project, 0o755); err != nil {
				return fmt.Errorf("ошибка в создании проекта: %w", err)
			}

			if err := os.Chdir(answers.Project); err != nil {
				return fmt.Errorf("ошибка при переходе в директорию: %w", err)
			}

			if answers.Git {
				gt := exec.Command("git", "init")
				gt.Stdout = os.Stdout
				gt.Stderr = os.Stderr

				if err := gt.Run(); err != nil {
					return fmt.Errorf("ошибка в init git: %w", err)
				}
			}

			fmt.Printf("создан проект: %s\nтип проекта: %s\nязык проекта: %s\ngit: %t\n", answers.Project, answers.Type, answers.Language, answers.Git)

			return nil
		},
	}

	cmd.Flags().StringVarP(&t, "type", "t", "", "тип проекта")
	cmd.Flags().StringVarP(&l, "lang", "l", "Go", "язык проекта")
	cmd.Flags().StringVarP(&n, "name", "n", "", "название проекта")
	cmd.Flags().BoolVarP(&git, "git", "g", false, "git в проект")
	cmd.Flags().BoolVar(&noInteractive, "no-interactive", false, "отключить интерактивный режим")

	return cmd
}

func newGenerateCmd() *cobra.Command {
	var t, o string
	var f, dry bool

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "сгенерировать отчет",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !slices.Contains(reportsTypes, t) {
				return fmt.Errorf("Недопустимый тип: %s", t)
			}

			if o == "" {
				o = "./out." + t
			}

			if dry {
				fmt.Printf("Будет сгенерирован отчёт типа: %s\n", t)
				fmt.Printf("Файл будет сохранён в: %s\n", o)
				return nil
			}

			logf("Запуск генерации отчёта...")
			logf("Тип отчёта: %s", t)
			logf("Путь сохранения: %s", o)

			if _, err := os.Stat(o); err == nil && !f {
				return fmt.Errorf("Ошибка: файл %q уже существует. Используйте --force для перезаписи", o)
			}

			if err := os.WriteFile(o, []byte("Отчёт ("+t+")\n"), 0o644); err != nil {
				return fmt.Errorf("Ошибка: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&t, "type", "t", "html", "тип отчета")
	cmd.Flags().StringVarP(&o, "output", "o", "", "путя для сохранения")
	cmd.Flags().BoolVarP(&f, "force", "f", false, "перезаписать существующий файл")
	cmd.Flags().BoolVar(&dry, "dry-run", false, "показать что будет сделано без выполнения")
	return cmd
}

func main() {
	root := &cobra.Command{
		Use:           "main-cli",
		Short:         "CLI-приложение для лабы 17",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.CompletionOptions.DisableDefaultCmd = true

	root.Flags().BoolP("version", "V", false, "версия программы")
	root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "подробный вывод")
	root.AddCommand(newGenerateCmd(), newConvertCmd(), newGreetCmd(), newInfoCmd(), newInteractiveCmd())
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
