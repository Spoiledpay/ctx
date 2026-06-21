package main

import (
    _ "embed"
    "fmt"
    "os"
    "strings"

    "github.com/ctx/ctx/internal/commands"
    "github.com/ctx/ctx/internal/i18n"
    "github.com/fatih/color"
    "github.com/spf13/cobra"
)

//go:embed version.bin
var versionData string

var (
    version = strings.TrimSpace(versionData)
    lang    string
    verbose bool
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "ctx",
        Short: i18n.T("ctx_short"),
        Long:  i18n.T("ctx_long"),
        PersistentPreRun: func(cmd *cobra.Command, args []string) {
            showBanner()
            if lang == "auto" {
                lang = detectSystemLanguage()
            }
            i18n.SetLanguage(lang)
            if verbose {
                color.Cyan(i18n.T("verbose_mode"), version)
            }
        },
        Run: func(cmd *cobra.Command, args []string) {
            cmd.Help()
        },
        Version: version,
    }

    // Flags globais
    rootCmd.PersistentFlags().StringVarP(&lang, "lang", "l", "auto", i18n.T("flag_lang"))
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, i18n.T("flag_verbose"))

    // Comandos
    rootCmd.AddCommand(commands.NewInitCmd())
    rootCmd.AddCommand(commands.NewScanCmd())
    rootCmd.AddCommand(commands.NewMemoryCmd())
    rootCmd.AddCommand(commands.NewDecisionCmd())
    rootCmd.AddCommand(commands.NewGrepCmd())
    rootCmd.AddCommand(commands.NewImpactCmd())
    rootCmd.AddCommand(commands.NewExplainCmd())

    if err := rootCmd.Execute(); err != nil {
        color.Red(i18n.T("error"), err)
        os.Exit(1)
    }
}

func showBanner() {
    fmt.Println("LabsObjects (R) CTX tool code diary version " + version)
    fmt.Println("Copyright (C) LabsOjects Project. All rights reserved.")
}

func detectSystemLanguage() string {
    // Detecta variáveis de ambiente
    for _, env := range []string{"LANG", "LANGUAGE", "LC_ALL"} {
        if val := os.Getenv(env); val != "" {
            switch {
            case len(val) >= 2 && val[:2] == "pt":
                return "pt"
            case len(val) >= 2 && val[:2] == "es":
                return "es"
            case len(val) >= 2 && val[:2] == "ru":
                return "ru"
            case len(val) >= 2 && val[:2] == "zh":
                return "zh"
            case len(val) >= 2 && val[:2] == "ja":
                return "jp"
            }
        }
    }
    return "en" // padrão
}