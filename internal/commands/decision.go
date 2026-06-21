package commands

import (
    "encoding/json"
    "fmt"
    "os"
    "sort"
    "strings"
    "time"

    "github.com/ctx/ctx/internal/core"
    "github.com/ctx/ctx/internal/i18n"
    "github.com/ctx/ctx/pkg/models"
    "github.com/fatih/color"
    "github.com/spf13/cobra"
    "github.com/google/uuid"
)

type DecisionOptions struct {
    ID           string
    Title        string
    Description  string
    Decision     string
    Rationale    string
    Alternatives []string
    Files        []string
    Participants []string
    Date         string
    Status       string
    DeprecatedBy string
    Tags         []string
    Links        []string
    From         string
    Output       string
    Limit        int
    Sort         string
}

func NewDecisionCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "decision",
        Short: i18n.T("decision_short"),
        Long:  i18n.T("decision_long"),
        Example: `  ctx decision add --title "Use PostgreSQL" --decision "PostgreSQL" --rationale "ACID compliance"
  ctx decision list --status accepted --tags database
  ctx decision show abc-123
  ctx decision deprecate abc-123 --reason "New architecture"`,
    }

    cmd.AddCommand(NewDecisionAddCmd())
    cmd.AddCommand(NewDecisionListCmd())
    cmd.AddCommand(NewDecisionShowCmd())
    cmd.AddCommand(NewDecisionUpdateCmd())
    cmd.AddCommand(NewDecisionDeleteCmd())
    cmd.AddCommand(NewDecisionDeprecateCmd())
    cmd.AddCommand(NewDecisionReactivateCmd())
    cmd.AddCommand(NewDecisionExportCmd())
    cmd.AddCommand(NewDecisionImportCmd())

    return cmd
}

func NewDecisionAddCmd() *cobra.Command {
    opts := &DecisionOptions{}
    
    cmd := &cobra.Command{
        Use:   "add",
        Short: i18n.T("decision_add_short"),
        Example: `  ctx decision add --title "Database Choice" --decision "PostgreSQL" --rationale "ACID requirements"
  ctx decision add --title "API Design" --decision "REST" --alternatives "GraphQL,gRPC" --participants "@joao,@maria"
  ctx decision add --from decision-template.json`,
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            if opts.From != "" {
                return addDecisionFromFile(ctx, opts.From)
            }

            return addDecisionInteractive(ctx, opts)
        },
    }

    // Flags básicas
    cmd.Flags().StringVarP(&opts.Title, "title", "t", "", i18n.T("decision_flag_title"))
    cmd.Flags().StringVarP(&opts.Description, "description", "d", "", i18n.T("decision_flag_desc"))
    cmd.Flags().StringVarP(&opts.Decision, "decision", "", "", i18n.T("decision_flag_decision"))
    cmd.Flags().StringVarP(&opts.Rationale, "rationale", "r", "", i18n.T("decision_flag_rationale"))
    
    // Alternativas
    cmd.Flags().StringSliceVarP(&opts.Alternatives, "alternatives", "a", []string{}, i18n.T("decision_flag_alternatives"))
    
    // Contexto
    cmd.Flags().StringSliceVarP(&opts.Files, "files", "f", []string{}, i18n.T("decision_flag_files"))
    cmd.Flags().StringSliceVarP(&opts.Participants, "participants", "p", []string{}, i18n.T("decision_flag_participants"))
    cmd.Flags().StringVarP(&opts.Date, "date", "", time.Now().Format("2006-01-02"), i18n.T("decision_flag_date"))
    
    // Status
    cmd.Flags().StringVarP(&opts.Status, "status", "s", "proposed", i18n.T("decision_flag_status"))
    
    // Metadados
    cmd.Flags().StringSliceVarP(&opts.Tags, "tags", "", []string{}, i18n.T("decision_flag_tags"))
    cmd.Flags().StringSliceVarP(&opts.Links, "links", "l", []string{}, i18n.T("decision_flag_links"))
    
    // Import
    cmd.Flags().StringVarP(&opts.From, "from", "", "", i18n.T("decision_flag_from"))

    // Marca flags obrigatórias
    cmd.MarkFlagRequired("title")
    cmd.MarkFlagRequired("decision")
    cmd.MarkFlagRequired("rationale")

    return cmd
}

func NewDecisionListCmd() *cobra.Command {
    opts := &DecisionOptions{}
    
    cmd := &cobra.Command{
        Use:   "list",
        Short: i18n.T("decision_list_short"),
        Example: `  ctx decision list
  ctx decision list --status accepted
  ctx decision list --tags database,architecture
  ctx decision list --limit 10 --sort date`,
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            decisions := ctx.Decisions

            // Filtros
            if opts.Status != "" {
                decisions = filterDecisionsByStatus(decisions, opts.Status)
            }
            if len(opts.Tags) > 0 {
                decisions = filterDecisionsByTags(decisions, opts.Tags)
            }
            if len(opts.Files) > 0 {
                decisions = filterDecisionsByFiles(decisions, opts.Files)
            }
            if opts.Date != "" {
                decisions = filterDecisionsByDate(decisions, opts.Date)
            }

            // Ordenação
            decisions = sortDecisions(decisions, opts.Sort)

            // Limite
            if opts.Limit > 0 && len(decisions) > opts.Limit {
                decisions = decisions[:opts.Limit]
            }

            // Output
            switch opts.Output {
            case "json":
                return printDecisionsJSON(decisions)
            case "table":
                printDecisionsTable(decisions)
            default:
                printDecisionsList(decisions)
            }

            return nil
        },
    }

    cmd.Flags().StringVarP(&opts.Status, "status", "s", "", i18n.T("decision_filter_status"))
    cmd.Flags().StringSliceVarP(&opts.Tags, "tags", "t", []string{}, i18n.T("decision_filter_tags"))
    cmd.Flags().StringSliceVarP(&opts.Files, "files", "f", []string{}, i18n.T("decision_filter_files"))
    cmd.Flags().StringVarP(&opts.Date, "date", "d", "", i18n.T("decision_filter_date"))
    cmd.Flags().StringVarP(&opts.Output, "output", "o", "text", i18n.T("decision_flag_output"))
    cmd.Flags().IntVarP(&opts.Limit, "limit", "l", 0, i18n.T("decision_flag_limit"))
    cmd.Flags().StringVarP(&opts.Sort, "sort", "", "date", i18n.T("decision_flag_sort"))

    return cmd
}

func NewDecisionShowCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "show [id]",
        Short: i18n.T("decision_show_short"),
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            decision := findDecisionByID(ctx, args[0])
            if decision == nil {
                return fmt.Errorf(i18n.T("decision_not_found"), args[0])
            }

            printDecisionDetails(decision)
            return nil
        },
    }
}

func NewDecisionUpdateCmd() *cobra.Command {
    opts := &DecisionOptions{}
    
    cmd := &cobra.Command{
        Use:   "update [id]",
        Short: i18n.T("decision_update_short"),
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            decision := findDecisionByID(ctx, args[0])
            if decision == nil {
                return fmt.Errorf(i18n.T("decision_not_found"), args[0])
            }

            updateDecision(decision, opts)

            if err := ctx.Save(); err != nil {
                return err
            }

            color.Green(i18n.T("decision_updated"), args[0])
            return nil
        },
    }

    // Flags para update
    cmd.Flags().StringVarP(&opts.Title, "title", "t", "", i18n.T("decision_flag_title"))
    cmd.Flags().StringVarP(&opts.Description, "description", "d", "", i18n.T("decision_flag_desc"))
    cmd.Flags().StringVarP(&opts.Decision, "decision", "", "", i18n.T("decision_flag_decision"))
    cmd.Flags().StringVarP(&opts.Rationale, "rationale", "r", "", i18n.T("decision_flag_rationale"))
    cmd.Flags().StringSliceVarP(&opts.Alternatives, "alternatives", "a", []string{}, i18n.T("decision_flag_alternatives"))
    cmd.Flags().StringSliceVarP(&opts.Files, "files", "f", []string{}, i18n.T("decision_flag_files"))
    cmd.Flags().StringSliceVarP(&opts.Participants, "participants", "p", []string{}, i18n.T("decision_flag_participants"))
    cmd.Flags().StringSliceVarP(&opts.Tags, "tags", "", []string{}, i18n.T("decision_flag_tags"))
    cmd.Flags().StringSliceVarP(&opts.Links, "links", "l", []string{}, i18n.T("decision_flag_links"))

    return cmd
}

func NewDecisionDeleteCmd() *cobra.Command {
    var force bool
    
    cmd := &cobra.Command{
        Use:   "delete [id]",
        Short: i18n.T("decision_delete_short"),
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
                color.Yellow(i18n.T("decision_delete_confirm"), args[0])
                fmt.Print(i18n.T("confirm") + " (y/N): ")
                var response string
                fmt.Scanln(&response)
                if strings.ToLower(response) != "y" {
                    return nil
                }
            }

            if !deleteDecision(ctx, args[0]) {
                return fmt.Errorf(i18n.T("decision_not_found"), args[0])
            }

            if err := ctx.Save(); err != nil {
                return err
            }

            color.Green(i18n.T("decision_deleted"), args[0])
            return nil
        },
    }

    cmd.Flags().BoolVarP(&force, "force", "f", false, i18n.T("decision_flag_force"))
    return cmd
}

func NewDecisionDeprecateCmd() *cobra.Command {
    var reason string
    var replacement string
    
    cmd := &cobra.Command{
        Use:   "deprecate [id]",
        Short: i18n.T("decision_deprecate_short"),
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            decision := findDecisionByID(ctx, args[0])
            if decision == nil {
                return fmt.Errorf(i18n.T("decision_not_found"), args[0])
            }

            decision.Status = "deprecated"
            decision.DeprecatedBy = replacement
            if reason != "" {
                decision.Rationale += "\n\nDeprecated: " + reason
            }

            if err := ctx.Save(); err != nil {
                return err
            }

            color.Yellow(i18n.T("decision_deprecated"), args[0])
            return nil
        },
    }

    cmd.Flags().StringVarP(&reason, "reason", "r", "", i18n.T("decision_flag_reason"))
    cmd.Flags().StringVarP(&replacement, "replacement", "", "", i18n.T("decision_flag_replacement"))
    
    return cmd
}

func NewDecisionReactivateCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "reactivate [id]",
        Short: i18n.T("decision_reactivate_short"),
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            decision := findDecisionByID(ctx, args[0])
            if decision == nil {
                return fmt.Errorf(i18n.T("decision_not_found"), args[0])
            }

            decision.Status = "accepted"
            decision.DeprecatedBy = ""

            if err := ctx.Save(); err != nil {
                return err
            }

            color.Green(i18n.T("decision_reactivated"), args[0])
            return nil
        },
    }

    return cmd
}

func NewDecisionExportCmd() *cobra.Command {
    var format string
    var output string
    var status string
    var tags []string
    
    cmd := &cobra.Command{
        Use:   "export",
        Short: i18n.T("decision_export_short"),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, err := core.LoadContext(".")
            if err != nil {
                return err
            }
            if ctx == nil {
                return fmt.Errorf(i18n.T("error_no_ctx"))
            }

            decisions := ctx.Decisions

            // Aplica filtros
            if status != "" {
                decisions = filterDecisionsByStatus(decisions, status)
            }
            if len(tags) > 0 {
                decisions = filterDecisionsByTags(decisions, tags)
            }

            var data []byte
            switch format {
            case "json":
                data, err = json.MarshalIndent(decisions, "", "  ")
            case "yaml":
                // Implementar YAML
                data = []byte("# YAML export pending")
            case "markdown":
                data = exportDecisionsMarkdown(decisions)
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
                color.Green(i18n.T("decision_exported"), output)
            }

            return nil
        },
    }

    cmd.Flags().StringVarP(&format, "format", "f", "json", i18n.T("decision_flag_format"))
    cmd.Flags().StringVarP(&output, "output", "o", "-", i18n.T("decision_flag_output"))
    cmd.Flags().StringVarP(&status, "status", "s", "", i18n.T("decision_filter_status"))
    cmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, i18n.T("decision_filter_tags"))

    return cmd
}

func NewDecisionImportCmd() *cobra.Command {
    var format string
    var merge bool
    
    cmd := &cobra.Command{
        Use:   "import [file]",
        Short: i18n.T("decision_import_short"),
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

            var decisions []*models.Decision
            switch format {
            case "json":
                err = json.Unmarshal(data, &decisions)
            default:
                return fmt.Errorf(i18n.T("error_format"), format)
            }

            if err != nil {
                return err
            }

            imported := 0
            for _, d := range decisions {
                if !merge {
                    // Gera novo ID
                    d.ID = uuid.New().String()
                    ctx.Decisions = append(ctx.Decisions, d)
                    imported++
                } else {
                    // Tenta fazer merge com existente
                    existing := findDecisionByID(ctx, d.ID)
                    if existing != nil {
                        mergeDecision(existing, d)
                    } else {
                        ctx.Decisions = append(ctx.Decisions, d)
                    }
                    imported++
                }
            }

            if err := ctx.Save(); err != nil {
                return err
            }

            color.Green(i18n.T("decision_imported"), imported, args[0])
            return nil
        },
    }

    cmd.Flags().StringVarP(&format, "format", "f", "json", i18n.T("decision_flag_format"))
    cmd.Flags().BoolVarP(&merge, "merge", "m", false, i18n.T("decision_flag_merge"))

    return cmd
}

// Funções auxiliares
func addDecisionInteractive(ctx *core.Context, opts *DecisionOptions) error {
    date, _ := time.Parse("2006-01-02", opts.Date)
    
    decision := &models.Decision{
        ID:           uuid.New().String(),
        Title:        opts.Title,
        Description:  opts.Description,
        Decision:     opts.Decision,
        Rationale:    opts.Rationale,
        Alternatives: opts.Alternatives,
        Files:        opts.Files,
        Participants: opts.Participants,
        Date:         date,
        Status:       opts.Status,
        Tags:         opts.Tags,
        Links:        opts.Links,
    }

    ctx.Decisions = append(ctx.Decisions, decision)

    if err := ctx.Save(); err != nil {
        return err
    }

    color.Green(i18n.T("decision_added"), decision.ID)
    printDecisionSummary(decision)

    return nil
}

func addDecisionFromFile(ctx *core.Context, path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    var decision models.Decision
    if err := json.Unmarshal(data, &decision); err != nil {
        return err
    }

    decision.ID = uuid.New().String()
    ctx.Decisions = append(ctx.Decisions, &decision)

    if err := ctx.Save(); err != nil {
        return err
    }

    color.Green(i18n.T("decision_added_from"), decision.ID, path)
    return nil
}

func findDecisionByID(ctx *core.Context, id string) *models.Decision {
    for _, d := range ctx.Decisions {
        if d.ID == id {
            return d
        }
    }
    return nil
}

func updateDecision(decision *models.Decision, opts *DecisionOptions) {
    if opts.Title != "" {
        decision.Title = opts.Title
    }
    if opts.Description != "" {
        decision.Description = opts.Description
    }
    if opts.Decision != "" {
        decision.Decision = opts.Decision
    }
    if opts.Rationale != "" {
        decision.Rationale = opts.Rationale
    }
    if len(opts.Alternatives) > 0 {
        decision.Alternatives = opts.Alternatives
    }
    if len(opts.Files) > 0 {
        decision.Files = opts.Files
    }
    if len(opts.Participants) > 0 {
        decision.Participants = opts.Participants
    }
    if len(opts.Tags) > 0 {
        decision.Tags = opts.Tags
    }
    if len(opts.Links) > 0 {
        decision.Links = opts.Links
    }
}

func mergeDecision(existing, new *models.Decision) {
    if new.Title != "" {
        existing.Title = new.Title
    }
    if new.Description != "" {
        existing.Description = new.Description
    }
    if new.Decision != "" {
        existing.Decision = new.Decision
    }
    if new.Rationale != "" {
        existing.Rationale = new.Rationale
    }
    // Merge slices
    existing.Alternatives = mergeStrings(existing.Alternatives, new.Alternatives)
    existing.Files = mergeStrings(existing.Files, new.Files)
    existing.Participants = mergeStrings(existing.Participants, new.Participants)
    existing.Tags = mergeStrings(existing.Tags, new.Tags)
    existing.Links = mergeStrings(existing.Links, new.Links)
}

func deleteDecision(ctx *core.Context, id string) bool {
    for i, d := range ctx.Decisions {
        if d.ID == id {
            ctx.Decisions = append(ctx.Decisions[:i], ctx.Decisions[i+1:]...)
            return true
        }
    }
    return false
}

func filterDecisionsByStatus(decisions []*models.Decision, status string) []*models.Decision {
    result := []*models.Decision{}
    for _, d := range decisions {
        if d.Status == status {
            result = append(result, d)
        }
    }
    return result
}

func filterDecisionsByTags(decisions []*models.Decision, tags []string) []*models.Decision {
    result := []*models.Decision{}
    for _, d := range decisions {
        for _, tag := range tags {
            if contains(d.Tags, tag) {
                result = append(result, d)
                break
            }
        }
    }
    return result
}

func filterDecisionsByFiles(decisions []*models.Decision, files []string) []*models.Decision {
    result := []*models.Decision{}
    for _, d := range decisions {
        for _, file := range files {
            if contains(d.Files, file) {
                result = append(result, d)
                break
            }
        }
    }
    return result
}

func filterDecisionsByDate(decisions []*models.Decision, date string) []*models.Decision {
    filterDate, _ := time.Parse("2006-01-02", date)
    result := []*models.Decision{}
    for _, d := range decisions {
        if d.Date.Year() == filterDate.Year() &&
           d.Date.Month() == filterDate.Month() &&
           d.Date.Day() == filterDate.Day() {
            result = append(result, d)
        }
    }
    return result
}

func sortDecisions(decisions []*models.Decision, sortBy string) []*models.Decision {
    switch sortBy {
    case "date":
        sort.Slice(decisions, func(i, j int) bool {
            return decisions[i].Date.After(decisions[j].Date)
        })
    case "title":
        sort.Slice(decisions, func(i, j int) bool {
            return decisions[i].Title < decisions[j].Title
        })
    case "status":
        sort.Slice(decisions, func(i, j int) bool {
            return decisions[i].Status < decisions[j].Status
        })
    }
    return decisions
}

func printDecisionsList(decisions []*models.Decision) {
    if len(decisions) == 0 {
        color.Yellow(i18n.T("decision_none"))
        return
    }

    for _, d := range decisions {
        printDecisionSummary(d)
    }
}

func printDecisionSummary(d *models.Decision) {
    var statusColor func(format string, a ...interface{}) string
    switch d.Status {
    case "accepted":
        statusColor = color.GreenString
    case "proposed":
        statusColor = color.CyanString
    case "deprecated":
        statusColor = color.RedString
    case "rejected":
        statusColor = color.YellowString
    default:
        statusColor = color.WhiteString
    }

    fmt.Printf("%s [%s] %s\n",
        statusColor("●"),
        strings.ToUpper(d.Status),
        color.WhiteString(d.Title))
    
    fmt.Printf("  ID: %s\n", color.CyanString(d.ID))
    fmt.Printf("  Decision: %s\n", color.YellowString(d.Decision))
    fmt.Printf("  Tags: %s\n", color.GreenString(strings.Join(d.Tags, ", ")))
    fmt.Printf("  Date: %s, Participants: %d\n", 
        d.Date.Format("2006-01-02"),
        len(d.Participants))
    fmt.Println()
}

func printDecisionDetails(d *models.Decision) {
    fmt.Println(strings.Repeat("═", 60))
    fmt.Printf("⚖️ %s\n", color.WhiteString(d.Title))
    fmt.Println(strings.Repeat("─", 60))
    
    fmt.Printf("ID:         %s\n", color.CyanString(d.ID))
    fmt.Printf("Status:     %s\n", colorForStatus(d.Status))
    fmt.Printf("Date:       %s\n", d.Date.Format("2006-01-02"))
    
    if len(d.Participants) > 0 {
        fmt.Printf("Participants: %s\n", strings.Join(d.Participants, ", "))
    }
    
    fmt.Println(strings.Repeat("─", 60))
    fmt.Println("DECISION:")
    fmt.Printf("  %s\n", color.YellowString(d.Decision))
    
    fmt.Println(strings.Repeat("─", 60))
    fmt.Println("RATIONALE:")
    fmt.Println(d.Rationale)
    
    if len(d.Alternatives) > 0 {
        fmt.Println(strings.Repeat("─", 60))
        fmt.Println("ALTERNATIVES CONSIDERED:")
        for _, alt := range d.Alternatives {
            fmt.Printf("  • %s\n", alt)
        }
    }
    
    if len(d.Files) > 0 {
        fmt.Println(strings.Repeat("─", 60))
        fmt.Println("RELATED FILES:")
        for _, f := range d.Files {
            fmt.Printf("  📄 %s\n", f)
        }
    }
    
    if len(d.Tags) > 0 {
        fmt.Println(strings.Repeat("─", 60))
        fmt.Println("TAGS:")
        fmt.Printf("  %s\n", strings.Join(d.Tags, ", "))
    }
    
    if len(d.Links) > 0 {
        fmt.Println(strings.Repeat("─", 60))
        fmt.Println("LINKS:")
        for _, l := range d.Links {
            fmt.Printf("  🔗 %s\n", l)
        }
    }
    
    if d.DeprecatedBy != "" {
        fmt.Println(strings.Repeat("─", 60))
        fmt.Printf("⚠️ DEPRECATED BY: %s\n", color.RedString(d.DeprecatedBy))
    }
    
    fmt.Println(strings.Repeat("═", 60))
}

func printDecisionsTable(decisions []*models.Decision) {
    // Implementar formato tabela
}

func printDecisionsJSON(decisions []*models.Decision) error {
    data, err := json.MarshalIndent(decisions, "", "  ")
    if err != nil {
        return err
    }
    fmt.Println(string(data))
    return nil
}

func exportDecisionsMarkdown(decisions []*models.Decision) []byte {
    var md strings.Builder
    
    md.WriteString("# Architecture Decisions Record (ADR)\n\n")
    
    for _, d := range decisions {
        md.WriteString(fmt.Sprintf("## %s\n", d.Title))
        md.WriteString(fmt.Sprintf("**ID:** %s  \n", d.ID))
        md.WriteString(fmt.Sprintf("**Status:** %s  \n", d.Status))
        md.WriteString(fmt.Sprintf("**Date:** %s  \n", d.Date.Format("2006-01-02")))
        md.WriteString(fmt.Sprintf("**Decision:** %s  \n\n", d.Decision))
        md.WriteString("### Rationale\n")
        md.WriteString(d.Rationale + "\n\n")
        
        if len(d.Alternatives) > 0 {
            md.WriteString("### Alternatives Considered\n")
            for _, alt := range d.Alternatives {
                md.WriteString(fmt.Sprintf("- %s\n", alt))
            }
            md.WriteString("\n")
        }
        
        md.WriteString("---\n\n")
    }
    
    return []byte(md.String())
}

func colorForStatus(status string) string {
    switch status {
    case "accepted":
        return color.GreenString(status)
    case "proposed":
        return color.CyanString(status)
    case "deprecated":
        return color.RedString(status)
    case "rejected":
        return color.YellowString(status)
    default:
        return color.WhiteString(status)
    }
}

func mergeStrings(a, b []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    
    for _, s := range a {
        if !seen[s] {
            seen[s] = true
            result = append(result, s)
        }
    }
    
    for _, s := range b {
        if !seen[s] {
            seen[s] = true
            result = append(result, s)
        }
    }
    
    return result
}