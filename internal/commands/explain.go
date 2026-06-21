package commands

import (
    "fmt"
    "path/filepath"
    "sort"
    "strings"

    "github.com/ctx/ctx/internal/core"
    "github.com/ctx/ctx/internal/i18n"
    "github.com/ctx/ctx/pkg/models"
    "github.com/ctx/ctx/pkg/utils"
    "github.com/fatih/color"
    "github.com/spf13/cobra"
)

type ExplainOptions struct {
    Detailed     bool
    IncludeBugs  bool
    IncludeDecisions bool
    IncludeHistory bool
    IncludeDeps  bool
    IncludeTests bool
    Format       string
    Output       string
    FromGit      bool
    Since        string
}

func NewExplainCmd() *cobra.Command {
    opts := &ExplainOptions{}
    
    cmd := &cobra.Command{
        Use:   "explain [file or function]",
        Short: i18n.T("explain_short"),
        Long:  i18n.T("explain_long"),
        Example: `  ctx explain auth.go
  ctx explain internal/service/auth.go
  ctx explain Login --detailed
  ctx explain auth.go --include-history
  ctx explain . --format json`,
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            target := args[0]
            
            // Se for diretório, explica o diretório
            if target == "." || target == "./" {
                return explainDirectory(ctx, opts)
            }

            // Tenta encontrar como arquivo
            if fileInfo, ok := ctx.Files[target]; ok {
                return explainFile(ctx, target, fileInfo, opts)
            }

            // Tenta encontrar como função
            for path, info := range ctx.Files {
                for _, fn := range info.Functions {
                    if fn == target {
                        return explainFunction(ctx, path, fn, opts)
                    }
                }
            }

            // Tenta encontrar por padrão (glob)
            matches, err := filepath.Glob(target)
            if err == nil && len(matches) > 0 {
                for _, match := range matches {
                    relPath, _ := filepath.Rel(ctx.Path, match)
                    if info, ok := ctx.Files[relPath]; ok {
                        explainFile(ctx, relPath, info, opts)
                        fmt.Println()
                    }
                }
                return nil
            }

            return fmt.Errorf(i18n.T("explain_not_found"), target)
        },
    }

    cmd.Flags().BoolVarP(&opts.Detailed, "detailed", "d", false, i18n.T("explain_flag_detailed"))
    cmd.Flags().BoolVarP(&opts.IncludeBugs, "bugs", "b", true, i18n.T("explain_flag_bugs"))
    cmd.Flags().BoolVarP(&opts.IncludeDecisions, "decisions", "", true, i18n.T("explain_flag_decisions"))
    cmd.Flags().BoolVarP(&opts.IncludeHistory, "history", "H", false, i18n.T("explain_flag_history"))
    cmd.Flags().BoolVarP(&opts.IncludeDeps, "deps", "", false, i18n.T("explain_flag_deps"))
    cmd.Flags().BoolVarP(&opts.IncludeTests, "tests", "t", false, i18n.T("explain_flag_tests"))
    cmd.Flags().StringVarP(&opts.Format, "format", "f", "text", i18n.T("explain_flag_format"))
    cmd.Flags().StringVarP(&opts.Output, "output", "o", "", i18n.T("explain_flag_output"))
    cmd.Flags().BoolVarP(&opts.FromGit, "from-git", "g", false, i18n.T("explain_flag_git"))
    cmd.Flags().StringVarP(&opts.Since, "since", "", "30d", i18n.T("explain_flag_since"))

    return cmd
}

func explainDirectory(ctx *core.Context, opts *ExplainOptions) error {
    color.Cyan(i18n.T("explain_directory"), ctx.Name)
    fmt.Println(strings.Repeat("─", 60))

    // Estatísticas gerais
    fmt.Printf("📊 %s\n", color.WhiteString(i18n.T("explain_stats")))
    fmt.Printf("  • %s: %d\n", i18n.T("stats_files"), len(ctx.Files))
    fmt.Printf("  • %s: %d\n", i18n.T("stats_memories"), len(ctx.Memories))
    fmt.Printf("  • %s: %d\n", i18n.T("stats_decisions"), len(ctx.Decisions))
    fmt.Printf("  • %s: %s\n", i18n.T("stats_language"), ctx.Language)
    fmt.Printf("  • %s: %s\n", i18n.T("stats_created"), ctx.CreatedAt.Format("2006-01-02"))

    // Top arquivos mais problemáticos
    if opts.IncludeBugs {
        buggyFiles := findBuggyFiles(ctx)
        if len(buggyFiles) > 0 {
            fmt.Println()
            fmt.Printf("🐛 %s\n", color.RedString(i18n.T("explain_buggy_files")))
            for i, f := range buggyFiles {
                if i >= 5 {
                    break
                }
                fmt.Printf("  • %s (%d bugs)\n", f.File, f.Count)
            }
        }
    }

    // Decisões recentes
    if opts.IncludeDecisions && len(ctx.Decisions) > 0 {
        fmt.Println()
        fmt.Printf("⚖️ %s\n", color.BlueString(i18n.T("explain_recent_decisions")))
        recent := getRecentDecisions(ctx, 5)
        for _, d := range recent {
            fmt.Printf("  • %s (%s)\n", d.Title, d.Date.Format("2006-01-02"))
        }
    }

    // Recomendações
    fmt.Println()
    fmt.Printf("💡 %s\n", color.GreenString(i18n.T("explain_recommendations")))
    recommendations := generateDirectoryRecommendations(ctx)
    for _, r := range recommendations {
        fmt.Printf("  • %s\n", r)
    }

    return nil
}

func explainFile(ctx *core.Context, path string, info *core.FileInfo, opts *ExplainOptions) error {
    color.Cyan(i18n.T("explain_file"), path)
    fmt.Println(strings.Repeat("─", 60))

    // Metadados básicos
    fmt.Printf("📄 %s\n", color.WhiteString(path))
    fmt.Printf("  • %s: %s\n", i18n.T("explain_package"), info.Package)
    fmt.Printf("  • %s: %d\n", i18n.T("explain_functions"), len(info.Functions))
    fmt.Printf("  • %s: %s\n", i18n.T("explain_modified"), info.ModifiedAt.Format("2006-01-02 15:04"))
    fmt.Printf("  • %s: %d bytes\n", i18n.T("explain_size"), info.Size)

    // Funções
    if len(info.Functions) > 0 && opts.Detailed {
        fmt.Println()
        fmt.Printf("🔧 %s\n", color.YellowString(i18n.T("explain_functions_list")))
        for _, fn := range info.Functions {
            fmt.Printf("  • %s\n", fn)
        }
    }

    // Dependências
    if opts.IncludeDeps {
        deps := ctx.Dependencies[path]
        if len(deps) > 0 {
            fmt.Println()
            fmt.Printf("📦 %s\n", color.CyanString(i18n.T("explain_dependencies")))
            for _, dep := range deps {
                fmt.Printf("  • %s\n", dep)
            }
        }

        // Quem depende deste arquivo
        dependents := findDependents(ctx, path)
        if len(dependents) > 0 {
            fmt.Println()
            fmt.Printf("🎯 %s\n", color.MagentaString(i18n.T("explain_dependents")))
            for _, dep := range dependents {
                fmt.Printf("  • %s\n", dep)
            }
        }
    }

    // Memórias relacionadas
    if opts.IncludeBugs {
        memories := findMemoriesForFile(ctx, path)
        if len(memories) > 0 {
            fmt.Println()
            fmt.Printf("📝 %s\n", color.RedString(i18n.T("explain_memories")))
            for _, m := range memories {
                fmt.Printf("  • [%s] %s (%s)\n", 
                    colorForSeverity(m.Severity),
                    m.Title,
                    m.CreatedAt.Format("2006-01-02"))
            }
        }
    }

    // Decisões relacionadas
    if opts.IncludeDecisions {
        decisions := findDecisionsForFile(ctx, path)
        if len(decisions) > 0 {
            fmt.Println()
            fmt.Printf("⚖️ %s\n", color.BlueString(i18n.T("explain_decisions")))
            for _, d := range decisions {
                fmt.Printf("  • %s (%s)\n", d.Title, d.Date.Format("2006-01-02"))
            }
        }
    }

    // Histórico Git
    if opts.IncludeHistory {
        history, err := utils.GetFileHistory(ctx.Path, path, opts.Since)
        if err == nil && len(history.Commits) > 0 {
            fmt.Println()
            fmt.Printf("📜 %s\n", color.CyanString(i18n.T("explain_history")))
            fmt.Printf("  • %s: %d\n", i18n.T("history_commits"), len(history.Commits))
            fmt.Printf("  • %s: %s\n", i18n.T("history_authors"), strings.Join(history.Authors, ", "))
            fmt.Printf("  • %s: %d\n", i18n.T("history_bugs"), history.BugCount)
            fmt.Printf("  • %s: %.2f/day\n", i18n.T("history_frequency"), history.ChangeFreq)

            if opts.Detailed {
                fmt.Println()
                fmt.Printf("  %s:\n", i18n.T("history_recent"))
                for i, commit := range history.Commits {
                    if i >= 5 {
                        break
                    }
                    fmt.Printf("    • %s - %s\n", 
                        commit.Date.Format("2006-01-02"),
                        commit.Message)
                }
            }
        }
    }

    // Testes relacionados
    if opts.IncludeTests {
        tests := findTestsForFile(ctx, path)
        if len(tests) > 0 {
            fmt.Println()
            fmt.Printf("🧪 %s\n", color.GreenString(i18n.T("explain_tests")))
            for _, test := range tests {
                fmt.Printf("  • %s\n", test)
            }
        }
    }

    // Riscos e recomendações
    risks := analyzeFileRisks(ctx, path, info)
    if len(risks) > 0 {
        fmt.Println()
        fmt.Printf("⚠️ %s\n", color.YellowString(i18n.T("explain_risks")))
        for _, risk := range risks {
            fmt.Printf("  • %s\n", risk)
        }
    }

    return nil
}

func explainFunction(ctx *core.Context, file, function string, opts *ExplainOptions) error {
    _ = ctx.Files[file]
    
    color.Cyan(i18n.T("explain_function"), function, file)
    fmt.Println(strings.Repeat("─", 60))

    // Onde é usada
    usages := findFunctionUsages(ctx, file, function)
    if len(usages) > 0 {
        fmt.Printf("📞 %s\n", color.YellowString(i18n.T("explain_usages")))
        for _, usage := range usages {
            fmt.Printf("  • %s\n", usage)
        }
    } else {
        fmt.Printf("📞 %s\n", color.YellowString(i18n.T("explain_no_usages")))
    }

    // Memórias relacionadas a esta função
    memories := findMemoriesForFunction(ctx, file, function)
    if len(memories) > 0 {
        fmt.Println()
        fmt.Printf("📝 %s\n", color.RedString(i18n.T("explain_function_memories")))
        for _, m := range memories {
            fmt.Printf("  • [%s] %s\n", 
                colorForSeverity(m.Severity),
                m.Title)
            if m.Solution != "" {
                fmt.Printf("    Solution: %s\n", m.Solution)
            }
        }
    }

    // Complexidade estimada
    complexity := estimateFunctionComplexity(ctx, file, function)
    fmt.Println()
    fmt.Printf("📊 %s: %s\n", i18n.T("explain_complexity"), complexity)

    return nil
}

// Tipos auxiliares
type BuggyFile struct {
    File  string
    Count int
}

func findBuggyFiles(ctx *core.Context) []BuggyFile {
    counts := make(map[string]int)
    
    for _, m := range ctx.Memories {
        if m.Type == "bug" {
            for _, f := range m.Files {
                counts[f]++
            }
        }
    }

    result := []BuggyFile{}
    for f, c := range counts {
        result = append(result, BuggyFile{File: f, Count: c})
    }

    sort.Slice(result, func(i, j int) bool {
        return result[i].Count > result[j].Count
    })

    return result
}

func getRecentDecisions(ctx *core.Context, limit int) []*models.Decision {
    decisions := make([]*models.Decision, len(ctx.Decisions))
    copy(decisions, ctx.Decisions)
    
    sort.Slice(decisions, func(i, j int) bool {
        return decisions[i].Date.After(decisions[j].Date)
    })

    if len(decisions) > limit {
        decisions = decisions[:limit]
    }

    return decisions
}

func generateDirectoryRecommendations(ctx *core.Context) []string {
    recommendations := []string{}

    // Recomendações baseadas em estatísticas
    if len(ctx.Memories) == 0 {
        recommendations = append(recommendations, i18n.T("rec_add_memories"))
    }

    buggyFiles := findBuggyFiles(ctx)
    if len(buggyFiles) > 0 {
        recommendations = append(recommendations, 
            fmt.Sprintf(i18n.T("rec_review_buggy"), buggyFiles[0].File))
    }

    // Verificar arquivos sem documentação
    undocumented := findUndocumentedFiles(ctx)
    if len(undocumented) > 0 {
        recommendations = append(recommendations,
            fmt.Sprintf(i18n.T("rec_document_files"), undocumented[0]))
    }

    return recommendations
}

func findDependents(ctx *core.Context, file string) []string {
    dependents := []string{}
    for path, deps := range ctx.Dependencies {
        for _, dep := range deps {
            if dep == file {
                dependents = append(dependents, path)
                break
            }
        }
    }
    return dependents
}

func findMemoriesForFile(ctx *core.Context, file string) []*models.Memory {
    memories := []*models.Memory{}
    for _, m := range ctx.Memories {
        for _, f := range m.Files {
            if f == file {
                memories = append(memories, m)
                break
            }
        }
    }
    return memories
}

func findDecisionsForFile(ctx *core.Context, file string) []*models.Decision {
    decisions := []*models.Decision{}
    for _, d := range ctx.Decisions {
        for _, f := range d.Files {
            if f == file {
                decisions = append(decisions, d)
                break
            }
        }
    }
    return decisions
}

func findTestsForFile(ctx *core.Context, file string) []string {
    tests := []string{}
    base := strings.TrimSuffix(file, ".go")
    testFile := base + "_test.go"
    
    if _, ok := ctx.Files[testFile]; ok {
        tests = append(tests, testFile)
    }
    
    return tests
}

func findFunctionUsages(ctx *core.Context, targetFile, function string) []string {
    usages := []string{}
    
    for file, info := range ctx.Files {
        for _, call := range info.Calls {
            if strings.Contains(call, function) {
                usages = append(usages, fmt.Sprintf("%s calls %s", file, function))
                break
            }
        }
    }
    
    return usages
}

func findMemoriesForFunction(ctx *core.Context, file, function string) []*models.Memory {
    memories := []*models.Memory{}
    
    for _, m := range ctx.Memories {
        for _, f := range m.Files {
            if f == file {
                // Simplificado: assumimos que a memória é sobre a função
                // Em produção, poderia ter linha específica
                memories = append(memories, m)
                break
            }
        }
    }
    
    return memories
}

func estimateFunctionComplexity(ctx *core.Context, file, function string) string {
    // Simplificado: em produção, analisaria AST
    return "medium"
}

func analyzeFileRisks(ctx *core.Context, file string, info *core.FileInfo) []string {
    risks := []string{}
    
    // Muitas funções pode indicar complexidade
    if len(info.Functions) > 10 {
        risks = append(risks, i18n.T("risk_many_functions"))
    }
    
    // Muitos bugs no passado
    bugCount := 0
    for _, m := range ctx.Memories {
        if m.Type == "bug" {
            for _, f := range m.Files {
                if f == file {
                    bugCount++
                    break
                }
            }
        }
    }
    if bugCount > 3 {
        risks = append(risks, fmt.Sprintf(i18n.T("risk_many_bugs"), bugCount))
    }
    
    // Muitas dependências
    deps := ctx.Dependencies[file]
    if len(deps) > 5 {
        risks = append(risks, fmt.Sprintf(i18n.T("risk_many_deps"), len(deps)))
    }
    
    return risks
}

func findUndocumentedFiles(ctx *core.Context) []string {
    undocumented := []string{}
    
    for file := range ctx.Files {
        // Verifica se tem memórias ou decisões associadas
        hasDoc := false
        for _, m := range ctx.Memories {
            for _, f := range m.Files {
                if f == file {
                    hasDoc = true
                    break
                }
            }
        }
        for _, d := range ctx.Decisions {
            for _, f := range d.Files {
                if f == file {
                    hasDoc = true
                    break
                }
            }
        }
        
        if !hasDoc && !strings.HasSuffix(file, "_test.go") {
            undocumented = append(undocumented, file)
        }
    }
    
    return undocumented
}