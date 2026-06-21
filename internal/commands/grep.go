package commands

import (
    "fmt"
    "strings"

    "github.com/ctx/ctx/internal/core"
    "github.com/ctx/ctx/internal/i18n"
    "github.com/ctx/ctx/pkg/models"
    "github.com/fatih/color"
    "github.com/spf13/cobra"
)

type GrepOptions struct {
    FilesOnly   bool
    Memories    bool
    Decisions   bool
    Code        bool
    Tags        []string
    Type        string
    Output      string
}

func NewGrepCmd() *cobra.Command {
    opts := &GrepOptions{}
    
    cmd := &cobra.Command{
        Use:   "grep [pattern]",
        Short: i18n.T("grep_short"),
        Long:  i18n.T("grep_long"),
        Example: `  ctx grep login
  ctx grep performance --tags critical
  ctx grep database --files-only`,
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            pattern := strings.ToLower(args[0])
            results := []string{}

            // Busca em memórias
            if opts.Memories {
                for _, m := range ctx.Memories {
                    if matchesPattern(m, pattern, opts) {
                        results = append(results, formatMemoryResult(m))
                    }
                }
            }

            // Busca em decisões
            if opts.Decisions {
                for _, d := range ctx.Decisions {
                    if matchesDecision(d, pattern, opts) {
                        results = append(results, formatDecisionResult(d))
                    }
                }
            }

            // Busca em código
            if opts.Code {
                for path, info := range ctx.Files {
                    if strings.Contains(strings.ToLower(path), pattern) {
                        results = append(results, fmt.Sprintf("📄 %s", path))
                    }
                    for _, fn := range info.Functions {
                        if strings.Contains(strings.ToLower(fn), pattern) {
                            results = append(results, fmt.Sprintf("🔧 %s -> %s", path, fn))
                        }
                    }
                }
            }

            if len(results) == 0 {
                color.Yellow(i18n.T("grep_no_results"), pattern)
                return nil
            }

            color.Cyan(i18n.T("grep_results"), len(results), pattern)
            for _, r := range results {
                fmt.Println(r)
            }

            return nil
        },
    }

    cmd.Flags().BoolVarP(&opts.FilesOnly, "files-only", "l", false, i18n.T("grep_flag_files"))
    cmd.Flags().BoolVarP(&opts.Memories, "memories", "m", true, i18n.T("grep_flag_memories"))
    cmd.Flags().BoolVarP(&opts.Decisions, "decisions", "d", true, i18n.T("grep_flag_decisions"))
    cmd.Flags().BoolVarP(&opts.Code, "code", "c", false, i18n.T("grep_flag_code"))
    cmd.Flags().StringSliceVarP(&opts.Tags, "tags", "t", []string{}, i18n.T("grep_flag_tags"))
    cmd.Flags().StringVarP(&opts.Type, "type", "", "", i18n.T("grep_flag_type"))

    return cmd
}

func matchesPattern(m *models.Memory, pattern string, opts *GrepOptions) bool {
    // Implementar lógica de busca
    return strings.Contains(strings.ToLower(m.Title), pattern) ||
           strings.Contains(strings.ToLower(m.Description), pattern)
}

func matchesDecision(d *models.Decision, pattern string, opts *GrepOptions) bool {
    return strings.Contains(strings.ToLower(d.Title), pattern) ||
           strings.Contains(strings.ToLower(d.Decision), pattern)
}

func formatMemoryResult(m *models.Memory) string {
    return fmt.Sprintf("📝 %s (%s)", m.Title, m.ID[:8])
}

func formatDecisionResult(d *models.Decision) string {
    return fmt.Sprintf("⚖️ %s (%s)", d.Title, d.ID[:8])
}