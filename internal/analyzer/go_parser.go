package analyzer

import (
    "go/ast"
    "go/parser"
    "go/token"
    "os"
    "path/filepath"
    "strings"

    "github.com/ctx/ctx/internal/core"
)

type GoParser struct {
    ctx *core.Context
}

func NewGoParser(ctx *core.Context) *GoParser {
    return &GoParser{ctx: ctx}
}

func (p *GoParser) ParseFile(path string) (*core.FileInfo, error) {
    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
    if err != nil {
        return nil, err
    }
    
    info := &core.FileInfo{
        Imports:   make([]string, 0),
        Functions: make([]string, 0),
        Calls:     make([]string, 0),
        Package:   node.Name.Name,
    }
    
    // Extrai imports
    for _, imp := range node.Imports {
        impPath := strings.Trim(imp.Path.Value, "\"")
        info.Imports = append(info.Imports, impPath)
    }
    
    // Extrai funções e chamadas
    ast.Inspect(node, func(n ast.Node) bool {
        switch x := n.(type) {
        case *ast.FuncDecl:
            info.Functions = append(info.Functions, x.Name.Name)
            
        case *ast.CallExpr:
            if fun, ok := x.Fun.(*ast.SelectorExpr); ok {
                if ident, ok := fun.X.(*ast.Ident); ok {
                    info.Calls = append(info.Calls, ident.Name+"."+fun.Sel.Name)
                }
            } else if fun, ok := x.Fun.(*ast.Ident); ok {
                info.Calls = append(info.Calls, fun.Name)
            }
        }
        return true
    })
    
    // Metadata do arquivo
    fileInfo, err := os.Stat(path)
    if err == nil {
        info.Size = fileInfo.Size()
        info.ModifiedAt = fileInfo.ModTime()
        info.Hash = calculateFileHash(path)
    }
    
    return info, nil
}

func (p *GoParser) ParseProject() error {
    root := p.ctx.Path
    
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if info.IsDir() {
            // Ignora pastas comuns
            if strings.HasPrefix(info.Name(), ".") || 
               info.Name() == "vendor" || 
               info.Name() == "node_modules" {
                return filepath.SkipDir
            }
            return nil
        }
        
        if strings.HasSuffix(path, ".go") {
            relPath, err := filepath.Rel(root, path)
            if err != nil {
                return err
            }
            
            fileInfo, err := p.ParseFile(path)
            if err != nil {
                return err
            }
            
            p.ctx.Files[relPath] = fileInfo
        }
        
        return nil
    })
    
    return err
}

func (p *GoParser) AnalyzeDependencies() {
    // Constrói grafo de dependências
    for file, info := range p.ctx.Files {
        deps := make([]string, 0)
        
        for _, imp := range info.Imports {
            // Converte import path para caminho relativo
            if strings.Contains(imp, p.ctx.Module) {
                rel := strings.TrimPrefix(imp, p.ctx.Module+"/")
                deps = append(deps, rel)
            }
        }
        
        p.ctx.Dependencies[file] = deps
    }
    
    // Dependências reversas
    for file, deps := range p.ctx.Dependencies {
        for _, dep := range deps {
            reverse := p.ctx.Dependencies[dep]
            if !contains(reverse, file) {
                p.ctx.Dependencies[dep] = append(p.ctx.Dependencies[dep], file)
            }
        }
    }
}

func calculateFileHash(path string) string {
    data, _ := os.ReadFile(path)
    _ = data
    return "hash-placeholder"
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}