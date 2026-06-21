package commands

import (
    "fmt"

    "github.com/ctx/ctx/internal/analyzer"
    "github.com/ctx/ctx/internal/core"
    "github.com/ctx/ctx/internal/i18n"
    "github.com/spf13/cobra"
)

func NewImpactCmd() *cobra.Command {
    opts := &analyzer.ImpactOptions{}
    
    cmd := &cobra.Command{
        Use:   "impact [files...]",
        Short: i18n.T("impact_short"),
        Long:  i18n.T("impact_long"),
        Example: `  ctx impact auth.go
  ctx impact internal/service/auth.go internal/handler/user.go
  ctx impact --recursive --depth 3`,
        Args: cobra.MinimumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            analyzer := analyzer.NewImpactAnalyzer(ctx, opts)
            results, err := analyzer.Analyze(args)
            if err != nil {
                return err
            }

            analyzer.PrintResults(results)
            return nil
        },
    }

    cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", true, i18n.T("impact_flag_recursive"))
    cmd.Flags().IntVarP(&opts.Depth, "depth", "d", 3, i18n.T("impact_flag_depth"))
    cmd.Flags().BoolVarP(&opts.IncludeTests, "include-tests", "t", false, i18n.T("impact_flag_tests"))
    cmd.Flags().BoolVarP(&opts.ShowMemories, "show-memories", "m", true, i18n.T("impact_flag_memories"))
    cmd.Flags().BoolVarP(&opts.ShowDecisions, "show-decisions", "", true, i18n.T("impact_flag_decisions"))
    cmd.Flags().BoolVarP(&opts.ShowHistory, "show-history", "H", false, i18n.T("impact_flag_history"))

    return cmd
}