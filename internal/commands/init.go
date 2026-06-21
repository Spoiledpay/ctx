package commands

import (
    "fmt"
    "os"
    "path/filepath"
    "time"

    "github.com/ctx/ctx/internal/core"
    "github.com/ctx/ctx/internal/i18n"
    "github.com/ctx/ctx/pkg/models"
    "github.com/fatih/color"
    "github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
    var force bool
    
    cmd := &cobra.Command{
        Use:   "init [directory]",
        Short: i18n.T("init_short"),
        Long:  i18n.T("init_long"),
        Args:  cobra.MaximumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            dir := "."
            if len(args) > 0 {
                dir = args[0]
            }
            
            absPath, err := filepath.Abs(dir)
            if err != nil {
                return err
            }
            
            color.Cyan(i18n.T("init_start"), absPath)
            
            // Verifica se já existe
            ctxPath := filepath.Join(absPath, ".ctx.yaml")
            if _, err := os.Stat(ctxPath); err == nil && !force {
                color.Yellow(i18n.T("init_exists"))
                return fmt.Errorf(i18n.T("init_exists"))
            }
            
            // Cria contexto novo
            ctx := &core.Context{
                Version:   "1.0",
                Name:      filepath.Base(absPath),
                Language:  detectLanguage(absPath),
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
                Files:     make(map[string]*core.FileInfo),
                Memories:  []*models.Memory{},
                Decisions: []*models.Decision{},
                Dependencies: make(map[string][]string),
                Tags:      make(map[string][]string),
                Path:      absPath,
            }
            
            if err := ctx.Save(); err != nil {
                return err
            }
            
            // Cria .gitignore se não existir
            gitignorePath := filepath.Join(absPath, ".gitignore")
            if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
                os.WriteFile(gitignorePath, []byte(".ctx.yaml\n.ctx.backup/\n"), 0644)
            }
            
            color.Green(i18n.T("init_success"))
            return nil
        },
    }
    
    cmd.Flags().BoolVarP(&force, "force", "f", false, "Force reinitialization")
    return cmd
}

func detectLanguage(dir string) string {
    // Detecta linguagem do projeto pelos arquivos
    files, _ := filepath.Glob(filepath.Join(dir, "*.go"))
    if len(files) > 0 {
        return "go"
    }
    
    if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
        return "javascript"
    }
    
    if _, err := os.Stat(filepath.Join(dir, "Cargo.toml")); err == nil {
        return "rust"
    }
    
    return "unknown"
}