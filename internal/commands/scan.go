package commands

import (
    "fmt"
    "path/filepath"
    "strings"
    "time"

    "github.com/ctx/ctx/internal/analyzer"
    "github.com/ctx/ctx/internal/core"
    "github.com/ctx/ctx/internal/i18n"
    "github.com/fatih/color"
    "github.com/spf13/cobra"
)

type ScanOptions struct {
    Recursive    bool
    IncludeTests bool
    Verbose      bool
    Force        bool
    Watch        bool
    Interval     int
    Output       string
}

func NewScanCmd() *cobra.Command {
    opts := &ScanOptions{}
    
    cmd := &cobra.Command{
        Use:   "scan [directory]",
        Short: i18n.T("scan_short"),
        Long:  i18n.T("scan_long"),
        Example: `  ctx scan                    # Scan current directory
  ctx scan ./internal         # Scan specific directory
  ctx scan --watch            # Watch for changes
  ctx scan --include-tests    # Include test files
  ctx scan --output json      # Output as JSON`,
        Args: cobra.MaximumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            dir := "."
            if len(args) > 0 {
                dir = args[0]
            }

            absPath, err := filepath.Abs(dir)
            if err != nil {
                return fmt.Errorf(i18n.T("error_path"), err)
            }

            // Carrega contexto existente ou cria novo
            ctx, err := core.LoadContext(absPath)
            if err != nil {
                return err
            }

            if ctx == nil {
                color.Yellow(i18n.T("scan_no_ctx"))
                return nil
            }

            if opts.Watch {
                return watchScan(ctx, opts)
            }

            return runScan(ctx, opts)
        },
    }

    cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", true, i18n.T("scan_flag_recursive"))
    cmd.Flags().BoolVarP(&opts.IncludeTests, "include-tests", "t", false, i18n.T("scan_flag_tests"))
    cmd.Flags().BoolVarP(&opts.Verbose, "verbose", "v", false, i18n.T("scan_flag_verbose"))
    cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, i18n.T("scan_flag_force"))
    cmd.Flags().BoolVarP(&opts.Watch, "watch", "w", false, i18n.T("scan_flag_watch"))
    cmd.Flags().IntVarP(&opts.Interval, "interval", "i", 5, i18n.T("scan_flag_interval"))
    cmd.Flags().StringVarP(&opts.Output, "output", "o", "text", i18n.T("scan_flag_output"))

    return cmd
}

func runScan(ctx *core.Context, opts *ScanOptions) error {
    color.Cyan(i18n.T("scan_start"), ctx.Path)

    startTime := time.Now()
    stats := &ScanStats{
        FilesScanned: 0,
        FilesChanged: 0,
        FilesNew:     0,
        Dependencies: 0,
        Memories:     len(ctx.Memories),
        Decisions:    len(ctx.Decisions),
    }

    // Inicializa parser apropriado
    if ctx.Language != "go" {
        return fmt.Errorf(i18n.T("scan_unsupported_lang"), ctx.Language)
    }

    oldFiles := ctx.Files
    ctx.Files = make(map[string]*core.FileInfo)

    parser := analyzer.NewGoParser(ctx)

    // Escaneia arquivos
    if err := parser.ParseProject(); err != nil {
        return err
    }

    stats.FilesScanned = len(ctx.Files)

    // Analisa mudanças
    for path, info := range ctx.Files {
        oldInfo, exists := oldFiles[path]
        
        if !exists {
            stats.FilesNew++
            color.Green("  + %s", path)
        } else if info.Hash != oldInfo.Hash {
            stats.FilesChanged++
            color.Yellow("  ~ %s", path)
        }
    }

    // Arquivos removidos
    for path := range oldFiles {
        if _, exists := ctx.Files[path]; !exists {
            color.Red("  - %s", path)
        }
    }

    // Analisa dependências
    if opts.Verbose {
        color.Cyan(i18n.T("scan_deps"))
    }
    
    parser.AnalyzeDependencies()
    stats.Dependencies = len(ctx.Dependencies)

    // Atualiza timestamps
    ctx.UpdatedAt = time.Now()

    // Salva contexto
    if err := ctx.Save(); err != nil {
        return err
    }

    // Mostra estatísticas
    duration := time.Since(startTime)
    printScanStats(stats, duration, opts)

    return nil
}

func watchScan(ctx *core.Context, opts *ScanOptions) error {
    color.Cyan(i18n.T("scan_watch_start"), opts.Interval)

    // Cria ticker para scan periódico
    ticker := time.NewTicker(time.Duration(opts.Interval) * time.Second)
    defer ticker.Stop()

    // Mapa para controle de mudanças
    lastHash := make(map[string]string)
    for path, info := range ctx.Files {
        lastHash[path] = info.Hash
    }

    for {
        select {
        case <-ticker.C:
            if err := watchScanIteration(ctx, opts, lastHash); err != nil {
                color.Red(i18n.T("error"), err)
            }
        }
    }
}

func watchScanIteration(ctx *core.Context, opts *ScanOptions, lastHash map[string]string) error {
    oldFiles := ctx.Files
    ctx.Files = make(map[string]*core.FileInfo)

    parser := analyzer.NewGoParser(ctx)
    
    if err := parser.ParseProject(); err != nil {
        return err
    }

    changes := false

    // Verifica mudanças
    for path, info := range ctx.Files {
        oldHash, exists := lastHash[path]
        
        if !exists {
            color.Green("  + %s", path)
            lastHash[path] = info.Hash
            changes = true
        } else if info.Hash != oldHash {
            color.Yellow("  ~ %s", path)
            lastHash[path] = info.Hash
            changes = true
        }

        // Remove old file entries that no longer exist
        delete(oldFiles, path)
    }

    // Arquivos removidos (still in oldFiles after iterating ctx.Files)
    for path := range oldFiles {
        color.Red("  - %s", path)
        delete(lastHash, path)
        delete(ctx.Files, path)
        changes = true
    }

    if changes {
        ctx.UpdatedAt = time.Now()
        if err := ctx.Save(); err != nil {
            return err
        }
        color.Green(i18n.T("scan_watch_updated"), time.Now().Format("15:04:05"))
    }

    return nil
}

type ScanStats struct {
    FilesScanned int
    FilesChanged int
    FilesNew     int
    Dependencies int
    Memories     int
    Decisions    int
}

func printScanStats(stats *ScanStats, duration time.Duration, opts *ScanOptions) {
    fmt.Println()
    color.Cyan("📊 " + i18n.T("scan_stats"))
    fmt.Println(strings.Repeat("─", 40))
    
    fmt.Printf("  %-20s %d\n", i18n.T("stats_files_scanned")+":", stats.FilesScanned)
    fmt.Printf("  %-20s %d\n", i18n.T("stats_files_new")+":", stats.FilesNew)
    fmt.Printf("  %-20s %d\n", i18n.T("stats_files_changed")+":", stats.FilesChanged)
    fmt.Printf("  %-20s %d\n", i18n.T("stats_dependencies")+":", stats.Dependencies)
    fmt.Printf("  %-20s %d\n", i18n.T("stats_memories")+":", stats.Memories)
    fmt.Printf("  %-20s %d\n", i18n.T("stats_decisions")+":", stats.Decisions)
    fmt.Printf("  %-20s %v\n", i18n.T("stats_duration")+":", duration.Round(time.Millisecond))
    
    if opts.Output == "json" {
        // Implementar output JSON
    }
}