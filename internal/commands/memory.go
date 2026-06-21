package commands

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
    "time"

    "github.com/ctx/ctx/internal/core"
    "github.com/ctx/ctx/internal/i18n"
    "github.com/ctx/ctx/pkg/models"
    "github.com/fatih/color"
    "github.com/spf13/cobra"
    "github.com/google/uuid"
)

type MemoryOptions struct {
    ID          string
    Title       string
    Description string
    Type        string
    Severity    string
    Files       []string
    Lines       []int
    Commit      string
    Branch      string
    Solution    string
    Workaround  string
    Tags        []string
    Author      string
    Reviewers   []string
    Links       []string
    From        string
    Output      string
    Limit       int
    Sort        string
    ResolvedAt  string
}

func NewMemoryCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "memory",
        Short: i18n.T("memory_short"),
        Long:  i18n.T("memory_long"),
    }

    cmd.AddCommand(NewMemoryAddCmd())
    cmd.AddCommand(NewMemoryListCmd())
    cmd.AddCommand(NewMemoryShowCmd())
    cmd.AddCommand(NewMemoryUpdateCmd())
    cmd.AddCommand(NewMemoryDeleteCmd())
    cmd.AddCommand(NewMemorySearchCmd())
    cmd.AddCommand(NewMemoryExportCmd())
    cmd.AddCommand(NewMemoryImportCmd())

    return cmd
}

func NewMemoryAddCmd() *cobra.Command {
    opts := &MemoryOptions{}
    
    cmd := &cobra.Command{
        Use:   "add",
        Short: i18n.T("memory_add_short"),
        Example: `  ctx memory add --title "Bug fix" --type bug --files auth.go
  ctx memory add --title "Decision" --type decision --tags architecture
  ctx memory add --from issue-123.json`,
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            // Se --from for fornecido, importa de arquivo
            if opts.From != "" {
                return addMemoryFromFile(ctx, opts.From)
            }

            return addMemoryInteractive(ctx, opts)
        },
    }

    // Flags básicas
    cmd.Flags().StringVarP(&opts.Title, "title", "t", "", i18n.T("memory_flag_title"))
    cmd.Flags().StringVarP(&opts.Description, "description", "d", "", i18n.T("memory_flag_desc"))
    cmd.Flags().StringVarP(&opts.Type, "type", "", "note", i18n.T("memory_flag_type"))
    cmd.Flags().StringVarP(&opts.Severity, "severity", "s", "medium", i18n.T("memory_flag_severity"))
    
    // Arquivos e linhas
    cmd.Flags().StringSliceVarP(&opts.Files, "files", "f", []string{}, i18n.T("memory_flag_files"))
    cmd.Flags().IntSliceVarP(&opts.Lines, "lines", "l", []int{}, i18n.T("memory_flag_lines"))
    
    // Git
    cmd.Flags().StringVarP(&opts.Commit, "commit", "c", "", i18n.T("memory_flag_commit"))
    cmd.Flags().StringVarP(&opts.Branch, "branch", "b", "", i18n.T("memory_flag_branch"))
    
    // Solução
    cmd.Flags().StringVarP(&opts.Solution, "solution", "", "", i18n.T("memory_flag_solution"))
    cmd.Flags().StringVarP(&opts.Workaround, "workaround", "", "", i18n.T("memory_flag_workaround"))
    
    // Metadados
    cmd.Flags().StringSliceVarP(&opts.Tags, "tags", "", []string{}, i18n.T("memory_flag_tags"))
    cmd.Flags().StringVarP(&opts.Author, "author", "a", "", i18n.T("memory_flag_author"))
    cmd.Flags().StringSliceVarP(&opts.Reviewers, "reviewers", "r", []string{}, i18n.T("memory_flag_reviewers"))
    cmd.Flags().StringSliceVarP(&opts.Links, "links", "", []string{}, i18n.T("memory_flag_links"))
    
    // Import/Export
    cmd.Flags().StringVarP(&opts.From, "from", "", "", i18n.T("memory_flag_from"))

    // Marca flags obrigatórias
    cmd.MarkFlagRequired("title")

    return cmd
}

func NewMemoryListCmd() *cobra.Command {
    opts := &MemoryOptions{}
    
    cmd := &cobra.Command{
        Use:   "list",
        Short: i18n.T("memory_list_short"),
        Example: `  ctx memory list
  ctx memory list --type bug --severity high
  ctx memory list --tags performance,security
  ctx memory list --limit 10 --sort date`,
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            memories := ctx.Memories

            // Filtros
            if opts.Type != "" {
                memories = filterMemoriesByType(memories, opts.Type)
            }
            if opts.Severity != "" {
                memories = filterMemoriesBySeverity(memories, opts.Severity)
            }
            if len(opts.Tags) > 0 {
                memories = filterMemoriesByTags(memories, opts.Tags)
            }
            if len(opts.Files) > 0 {
                memories = filterMemoriesByFiles(memories, opts.Files)
            }

            // Ordenação
            memories = sortMemories(memories, opts.Sort)

            // Limite
            if opts.Limit > 0 && len(memories) > opts.Limit {
                memories = memories[:opts.Limit]
            }

            // Output
            switch opts.Output {
            case "json":
                return printMemoriesJSON(memories)
            case "table":
                printMemoriesTable(memories)
            default:
                printMemoriesList(memories)
            }

            return nil
        },
    }

    cmd.Flags().StringVarP(&opts.Type, "type", "t", "", i18n.T("memory_filter_type"))
    cmd.Flags().StringVarP(&opts.Severity, "severity", "s", "", i18n.T("memory_filter_severity"))
    cmd.Flags().StringSliceVarP(&opts.Tags, "tags", "", []string{}, i18n.T("memory_filter_tags"))
    cmd.Flags().StringSliceVarP(&opts.Files, "files", "f", []string{}, i18n.T("memory_filter_files"))
    cmd.Flags().StringVarP(&opts.Output, "output", "o", "text", i18n.T("memory_flag_output"))
    cmd.Flags().IntVarP(&opts.Limit, "limit", "l", 0, i18n.T("memory_flag_limit"))
    cmd.Flags().StringVarP(&opts.Sort, "sort", "", "date", i18n.T("memory_flag_sort"))

    return cmd
}

func NewMemoryShowCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "show [id]",
        Short: i18n.T("memory_show_short"),
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            memory := findMemoryByID(ctx, args[0])
            if memory == nil {
                return fmt.Errorf(i18n.T("memory_not_found"), args[0])
            }

            printMemoryDetails(memory)
            return nil
        },
    }
}

func NewMemoryUpdateCmd() *cobra.Command {
    opts := &MemoryOptions{}
    
    cmd := &cobra.Command{
        Use:   "update [id]",
        Short: i18n.T("memory_update_short"),
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            memory := findMemoryByID(ctx, args[0])
            if memory == nil {
                return fmt.Errorf(i18n.T("memory_not_found"), args[0])
            }

            // Atualiza campos
            updateMemory(memory, opts)
            memory.UpdatedAt = time.Now()

            if err := ctx.Save(); err != nil {
                return err
            }

            color.Green(i18n.T("memory_updated"), args[0])
            return nil
        },
    }

    // Mesmas flags do add, mas opcionais
    cmd.Flags().StringVarP(&opts.Title, "title", "t", "", i18n.T("memory_flag_title"))
    cmd.Flags().StringVarP(&opts.Description, "description", "d", "", i18n.T("memory_flag_desc"))
    cmd.Flags().StringVarP(&opts.Type, "type", "", "", i18n.T("memory_flag_type"))
    cmd.Flags().StringVarP(&opts.Severity, "severity", "s", "", i18n.T("memory_flag_severity"))
    cmd.Flags().StringSliceVarP(&opts.Files, "files", "f", []string{}, i18n.T("memory_flag_files"))
    cmd.Flags().StringVarP(&opts.Solution, "solution", "", "", i18n.T("memory_flag_solution"))
    cmd.Flags().StringSliceVarP(&opts.Tags, "tags", "", []string{}, i18n.T("memory_flag_tags"))
    cmd.Flags().StringVarP(&opts.ResolvedAt, "resolved", "", "", i18n.T("memory_flag_resolved"))

    return cmd
}

func NewMemoryDeleteCmd() *cobra.Command {
    var force bool
    
    cmd := &cobra.Command{
        Use:   "delete [id]",
        Short: i18n.T("memory_delete_short"),
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            if !force {
                color.Yellow(i18n.T("memory_delete_confirm"), args[0])
                fmt.Print(i18n.T("confirm") + " (y/N): ")
                var response string
                fmt.Scanln(&response)
                if strings.ToLower(response) != "y" {
                    return nil
                }
            }

            if !deleteMemory(ctx, args[0]) {
                return fmt.Errorf(i18n.T("memory_not_found"), args[0])
            }

            if err := ctx.Save(); err != nil {
                return err
            }

            color.Green(i18n.T("memory_deleted"), args[0])
            return nil
        },
    }

    cmd.Flags().BoolVarP(&force, "force", "f", false, i18n.T("memory_flag_force"))
    return cmd
}

func NewMemorySearchCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "search [query]",
        Short: i18n.T("memory_search_short"),
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            query := strings.ToLower(args[0])
            results := []*models.Memory{}

            for _, memory := range ctx.Memories {
                if strings.Contains(strings.ToLower(memory.Title), query) ||
                   strings.Contains(strings.ToLower(memory.Description), query) ||
                   strings.Contains(strings.ToLower(memory.Solution), query) ||
                   containsAny(strings.ToLower(strings.Join(memory.Tags, " ")), query) {
                    results = append(results, memory)
                }
            }

            if len(results) == 0 {
                color.Yellow(i18n.T("memory_search_none"), query)
                return nil
            }

            color.Cyan(i18n.T("memory_search_results"), len(results), query)
            printMemoriesList(results)
            
            return nil
        },
    }

    return cmd
}

func NewMemoryExportCmd() *cobra.Command {
    var format string
    var output string
    
    cmd := &cobra.Command{
        Use:   "export",
        Short: i18n.T("memory_export_short"),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            var data []byte
            switch format {
            case "json":
                data, err = json.MarshalIndent(ctx.Memories, "", "  ")
            case "yaml":
                // Implementar YAML export
            case "csv":
                // Implementar CSV export
            default:
                return fmt.Errorf(i18n.T("error_format"), format)
            }

            if err != nil {
                return err
            }

            if output == "-" {
                fmt.Println(string(data))
            } else {
                if err := os.WriteFile(output, data, 0644); err != nil {
                    return err
                }
                color.Green(i18n.T("memory_exported"), output)
            }

            return nil
        },
    }

    cmd.Flags().StringVarP(&format, "format", "f", "json", i18n.T("memory_flag_format"))
    cmd.Flags().StringVarP(&output, "output", "o", "-", i18n.T("memory_flag_output"))

    return cmd
}

func NewMemoryImportCmd() *cobra.Command {
    var format string
    
    cmd := &cobra.Command{
        Use:   "import [file]",
        Short: i18n.T("memory_import_short"),
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            data, err := os.ReadFile(args[0])
            if err != nil {
                return err
            }

            var memories []*models.Memory
            switch format {
            case "json":
                err = json.Unmarshal(data, &memories)
            default:
                return fmt.Errorf(i18n.T("error_format"), format)
            }

            if err != nil {
                return err
            }

            for _, m := range memories {
                m.ID = uuid.New().String()
                m.CreatedAt = time.Now()
                m.UpdatedAt = time.Now()
                ctx.Memories = append(ctx.Memories, m)
            }

            if err := ctx.Save(); err != nil {
                return err
            }

            color.Green(i18n.T("memory_imported"), len(memories), args[0])
            return nil
        },
    }

    cmd.Flags().StringVarP(&format, "format", "f", "json", i18n.T("memory_flag_format"))
    return cmd
}

// Funções auxiliares
func addMemoryInteractive(ctx *core.Context, opts *MemoryOptions) error {
    memory := &models.Memory{
        ID:          uuid.New().String(),
        Title:       opts.Title,
        Description: opts.Description,
        Type:        opts.Type,
        Severity:    opts.Severity,
        Files:       opts.Files,
        Lines:       opts.Lines,
        Commit:      opts.Commit,
        Branch:      opts.Branch,
        Solution:    opts.Solution,
        Workaround:  opts.Workaround,
        Tags:        opts.Tags,
        Author:      opts.Author,
        Reviewers:   opts.Reviewers,
        Links:       opts.Links,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    // Se não forneceu autor, tenta pegar do git
    if memory.Author == "" {
        if author, err := getGitUser(); err == nil {
            memory.Author = author
        }
    }

    // Se não forneceu branch, tenta pegar do git
    if memory.Branch == "" {
        if branch, err := getGitBranch(); err == nil {
            memory.Branch = branch
        }
    }

    ctx.AddMemory(memory)
    
    if err := ctx.Save(); err != nil {
        return err
    }

    color.Green(i18n.T("memory_added"), memory.ID)
    printMemorySummary(memory)
    
    return nil
}

func addMemoryFromFile(ctx *core.Context, path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    var memory models.Memory
    if err := json.Unmarshal(data, &memory); err != nil {
        return err
    }

    memory.ID = uuid.New().String()
    memory.CreatedAt = time.Now()
    memory.UpdatedAt = time.Now()

    ctx.AddMemory(&memory)
    
    if err := ctx.Save(); err != nil {
        return err
    }

    color.Green(i18n.T("memory_added_from"), memory.ID, path)
    return nil
}

func findMemoryByID(ctx *core.Context, id string) *models.Memory {
    for _, m := range ctx.Memories {
        if m.ID == id {
            return m
        }
    }
    return nil
}

func updateMemory(memory *models.Memory, opts *MemoryOptions) {
    if opts.Title != "" {
        memory.Title = opts.Title
    }
    if opts.Description != "" {
        memory.Description = opts.Description
    }
    if opts.Type != "" {
        memory.Type = opts.Type
    }
    if opts.Severity != "" {
        memory.Severity = opts.Severity
    }
    if len(opts.Files) > 0 {
        memory.Files = opts.Files
    }
    if opts.Solution != "" {
        memory.Solution = opts.Solution
    }
    if len(opts.Tags) > 0 {
        memory.Tags = opts.Tags
    }
    if opts.ResolvedAt != "" {
        if t, err := time.Parse(time.RFC3339, opts.ResolvedAt); err == nil {
            memory.ResolvedAt = &t
        }
    }
}

func deleteMemory(ctx *core.Context, id string) bool {
    for i, m := range ctx.Memories {
        if m.ID == id {
            ctx.Memories = append(ctx.Memories[:i], ctx.Memories[i+1:]...)
            return true
        }
    }
    return false
}

func filterMemoriesByType(memories []*models.Memory, mtype string) []*models.Memory {
    result := []*models.Memory{}
    for _, m := range memories {
        if m.Type == mtype {
            result = append(result, m)
        }
    }
    return result
}

func filterMemoriesBySeverity(memories []*models.Memory, severity string) []*models.Memory {
    result := []*models.Memory{}
    for _, m := range memories {
        if m.Severity == severity {
            result = append(result, m)
        }
    }
    return result
}

func filterMemoriesByTags(memories []*models.Memory, tags []string) []*models.Memory {
    result := []*models.Memory{}
    for _, m := range memories {
        for _, tag := range tags {
            if contains(m.Tags, tag) {
                result = append(result, m)
                break
            }
        }
    }
    return result
}

func filterMemoriesByFiles(memories []*models.Memory, files []string) []*models.Memory {
    result := []*models.Memory{}
    for _, m := range memories {
        for _, file := range files {
            if contains(m.Files, file) {
                result = append(result, m)
                break
            }
        }
    }
    return result
}

func sortMemories(memories []*models.Memory, sortBy string) []*models.Memory {
    // Implementar ordenação
    return memories
}

func printMemoriesList(memories []*models.Memory) {
    if len(memories) == 0 {
        color.Yellow(i18n.T("memory_none"))
        return
    }

    for _, m := range memories {
        printMemorySummary(m)
    }
}

func printMemorySummary(m *models.Memory) {
    // Cores por tipo
    var typeColor func(format string, a ...interface{}) string
    switch m.Type {
    case "bug":
        typeColor = color.RedString
    case "decision":
        typeColor = color.BlueString
    case "warning":
        typeColor = color.YellowString
    default:
        typeColor = color.WhiteString
    }

    // Cores por severidade
    var severityColor func(format string, a ...interface{}) string
    switch m.Severity {
    case "critical":
        severityColor = color.RedString
    case "high":
        severityColor = color.YellowString
    case "medium":
        severityColor = color.CyanString
    default:
        severityColor = color.GreenString
    }

    fmt.Printf("%s [%s] %s\n",
        severityColor("●"),
        typeColor(strings.ToUpper(m.Type)),
        color.WhiteString(m.Title))
    
    fmt.Printf("  ID: %s\n", color.CyanString(m.ID))
    fmt.Printf("  Tags: %s\n", color.GreenString(strings.Join(m.Tags, ", ")))
    fmt.Printf("  Files: %s\n", color.YellowString(strings.Join(m.Files, ", ")))
    fmt.Printf("  Author: %s, %s\n", m.Author, m.CreatedAt.Format("2006-01-02"))
    fmt.Println()
}

func printMemoryDetails(m *models.Memory) {
    fmt.Println(strings.Repeat("═", 60))
    fmt.Printf("📝 %s\n", color.WhiteString(m.Title))
    fmt.Println(strings.Repeat("─", 60))
    
    fmt.Printf("ID:        %s\n", color.CyanString(m.ID))
    fmt.Printf("Type:      %s\n", colorForType(m.Type))
    fmt.Printf("Severity:  %s\n", colorForSeverity(m.Severity))
    fmt.Printf("Created:   %s by %s\n", m.CreatedAt.Format("2006-01-02 15:04"), m.Author)
    if m.ResolvedAt != nil {
        fmt.Printf("Resolved:  %s\n", m.ResolvedAt.Format("2006-01-02 15:04"))
    }
    
    fmt.Println(strings.Repeat("─", 60))
    fmt.Println("DESCRIPTION:")
    fmt.Println(m.Description)
    
    if m.Solution != "" {
        fmt.Println(strings.Repeat("─", 60))
        fmt.Println("SOLUTION:")
        fmt.Println(m.Solution)
    }
    
    if len(m.Files) > 0 {
        fmt.Println(strings.Repeat("─", 60))
        fmt.Println("FILES:")
        for _, f := range m.Files {
            fmt.Printf("  📄 %s\n", f)
        }
    }
    
    if len(m.Tags) > 0 {
        fmt.Println(strings.Repeat("─", 60))
        fmt.Println("TAGS:")
        fmt.Printf("  %s\n", strings.Join(m.Tags, ", "))
    }
    
    if len(m.Links) > 0 {
        fmt.Println(strings.Repeat("─", 60))
        fmt.Println("LINKS:")
        for _, l := range m.Links {
            fmt.Printf("  🔗 %s\n", l)
        }
    }
    
    fmt.Println(strings.Repeat("═", 60))
}

func printMemoriesTable(memories []*models.Memory) {
    // Implementar formato tabela
}

func printMemoriesJSON(memories []*models.Memory) error {
    data, err := json.MarshalIndent(memories, "", "  ")
    if err != nil {
        return err
    }
    fmt.Println(string(data))
    return nil
}

func colorForType(t string) string {
    switch t {
    case "bug":
        return color.RedString(t)
    case "decision":
        return color.BlueString(t)
    case "warning":
        return color.YellowString(t)
    default:
        return color.WhiteString(t)
    }
}

func colorForSeverity(s string) string {
    switch s {
    case "critical":
        return color.RedString(s)
    case "high":
        return color.YellowString(s)
    case "medium":
        return color.CyanString(s)
    default:
        return color.GreenString(s)
    }
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

func containsAny(s, substr string) bool {
    return strings.Contains(s, substr)
}

func getGitUser() (string, error) {
    // Implementar git user
    return "unknown", nil
}

func getGitBranch() (string, error) {
    // Implementar git branch
    return "main", nil
}