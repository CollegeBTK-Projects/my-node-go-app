package main

import (
	"fmt"
	"os"
	"slices"

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
			if output == "" {
				output = input + "." + format
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
	root.AddCommand(newGenerateCmd(), newConvertCmd(), newGreetCmd(), newInfoCmd())
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
