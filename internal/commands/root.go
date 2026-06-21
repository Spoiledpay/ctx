package commands

import (
    "github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "ctx",
        Short: "CTX - Code Context & Memory",
        Long: `CTX is a source code diary that stores decisions, 
bugs, and context about your project.`,
    }

    // Adiciona todos os comandos
    cmd.AddCommand(NewInitCmd())
    cmd.AddCommand(NewScanCmd())
    cmd.AddCommand(NewMemoryCmd())
    cmd.AddCommand(NewDecisionCmd())
    cmd.AddCommand(NewExplainCmd())
    cmd.AddCommand(NewImpactCmd())
    cmd.AddCommand(NewGrepCmd())
    cmd.AddCommand(NewVersionCmd())

    return cmd
}