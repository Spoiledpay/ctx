package analyzer

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"

    "github.com/ctx/ctx/internal/core"
    "github.com/ctx/ctx/pkg/models"
    "github.com/fatih/color"
)

type ImpactAnalyzer struct {
    ctx     *core.Context
    options *ImpactOptions
}

type ImpactOptions struct {
    Recursive     bool
    IncludeTests  bool
    ShowMemories  bool
    ShowDecisions bool
    ShowHistory   bool
    Depth         int
    Format        string
}

type ImpactResult struct {
    File           string          `json:"file"`
    DirectImpact   []ImpactItem    `json:"direct_impact"`
    IndirectImpact []ImpactItem    `json:"indirect_impact"`
    Tests          []string        `json:"tests"`
    DockerServices []string        `json:"docker_services"`
    Memories       []*models.Memory `json:"memories,omitempty"`
    Decisions      []*models.Decision `json:"decisions,omitempty"`
    History        *FileHistory     `json:"history,omitempty"`
    Risk           string           `json:"risk"`
    RiskScore      int              `json:"risk_score"`
}

type ImpactItem struct {
    File     string   `json:"file"`
    Path     string   `json:"path"`
    Type     string   `json:"type"` // direct, indirect
    Depth    int      `json:"depth"`
    Reasons  []string `json:"reasons"`
    Package  string   `json:"package,omitempty"`
}

type FileHistory struct {
    Commits     int       `json:"commits"`
    Authors     []string  `json:"authors"`
    LastChange  time.Time `json:"last_change"`
    BugCount    int       `json:"bug_count"`
    ChangeFreq  float64   `json:"change_frequency"`
}

func NewImpactAnalyzer(ctx *core.Context, opts *ImpactOptions) *ImpactAnalyzer {
    if opts == nil {
        opts = &ImpactOptions{
            Recursive:    true,
            IncludeTests: false,
            ShowMemories: true,
            ShowDecisions: true,
            ShowHistory:  true,
            Depth:        3,
            Format:       "text",
        }
    }
    return &ImpactAnalyzer{
        ctx:     ctx,
        options: opts,
    }
}

func (a *ImpactAnalyzer) Analyze(files []string) (map[string]*ImpactResult, error) {
    results := make(map[string]*ImpactResult)

    for _, file := range files {
        result, err := a.analyzeFile(file)
        if err != nil {
            return nil, err
        }
        results[file] = result
    }

    return results, nil
}

func (a *ImpactAnalyzer) analyzeFile(file string) (*ImpactResult, error) {
    result := &ImpactResult{
        File:           file,
        DirectImpact:   []ImpactItem{},
        IndirectImpact: []ImpactItem{},
        Tests:          []string{},
        DockerServices: []string{},
        Memories:       []*models.Memory{},
        Decisions:      []*models.Decision{},
        Risk:           "low",
        RiskScore:      0,
    }

    // Encontra dependências diretas (arquivos que importam este)
    direct := a.findDirectDependents(file)
    for _, d := range direct {
        item := ImpactItem{
            File:    d,
            Path:    d,
            Type:    "direct",
            Depth:   1,
            Reasons: a.findReasons(file, d),
        }
        
        // Tenta encontrar package
        if info, ok := a.ctx.Files[d]; ok {
            item.Package = info.Package
        }
        
        result.DirectImpact = append(result.DirectImpact, item)
    }

    // Encontra dependências indiretas
    if a.options.Recursive {
        indirect := a.findIndirectDependents(file, 1)
        for _, d := range indirect {
            item := ImpactItem{
                File:    d.File,
                Path:    d.File,
                Type:    "indirect",
                Depth:   d.Depth,
                Reasons: a.buildIndirectReasons(file, d.File),
            }
            result.IndirectImpact = append(result.IndirectImpact, item)
        }
    }

    // Encontra testes relacionados
    result.Tests = a.findRelatedTests(file)

    // Encontra serviços Docker relacionados
    result.DockerServices = a.findDockerServices(file)

    // Adiciona memórias relacionadas
    if a.options.ShowMemories {
        result.Memories = a.findRelatedMemories(file)
    }

    // Adiciona decisões relacionadas
    if a.options.ShowDecisions {
        result.Decisions = a.findRelatedDecisions(file)
    }

    // Análise de histórico
    if a.options.ShowHistory {
        history, err := a.analyzeHistory(file)
        if err == nil {
            result.History = history
        }
    }

    // Calcula risco
    result.RiskScore = a.calculateRiskScore(result)
    result.Risk = a.riskLevel(result.RiskScore)

    return result, nil
}

func (a *ImpactAnalyzer) findDirectDependents(file string) []string {
    dependents := []string{}
    
    for path, deps := range a.ctx.Dependencies {
        for _, dep := range deps {
            if dep == file {
                dependents = append(dependents, path)
                break
            }
        }
    }
    
    return dependents
}

type IndirectDependent struct {
    File  string
    Depth int
}

func (a *ImpactAnalyzer) findIndirectDependents(file string, currentDepth int) []IndirectDependent {
    if currentDepth > a.options.Depth {
        return []IndirectDependent{}
    }

    result := []IndirectDependent{}
    seen := make(map[string]bool)

    // Encontra dependentes diretos primeiro
    direct := a.findDirectDependents(file)
    
    for _, d := range direct {
        if !seen[d] {
            seen[d] = true
            result = append(result, IndirectDependent{File: d, Depth: currentDepth})
            
            // Recursivamente encontra dependentes indiretos
            indirect := a.findIndirectDependents(d, currentDepth+1)
            for _, id := range indirect {
                if !seen[id.File] {
                    seen[id.File] = true
                    result = append(result, id)
                }
            }
        }
    }

    return result
}

func (a *ImpactAnalyzer) findReasons(from, to string) []string {
    reasons := []string{}
    
    if info, ok := a.ctx.Files[from]; ok {
        // Verifica imports
        for _, imp := range info.Imports {
            if strings.Contains(imp, to) {
                reasons = append(reasons, fmt.Sprintf("imports %s", to))
            }
        }
        
        // Verifica chamadas de função
        for _, call := range info.Calls {
            if strings.Contains(call, filepath.Base(to)) {
                reasons = append(reasons, fmt.Sprintf("calls function from %s", to))
            }
        }
    }
    
    return reasons
}

func (a *ImpactAnalyzer) buildIndirectReasons(original, dependent string) []string {
    reasons := []string{}
    path := a.findDependencyPath(original, dependent)
    
    if len(path) > 0 {
        reasons = append(reasons, fmt.Sprintf("dependency chain: %s", strings.Join(path, " → ")))
    }
    
    return reasons
}

func (a *ImpactAnalyzer) findDependencyPath(from, to string) []string {
    // BFS para encontrar caminho de dependência
    visited := make(map[string]bool)
    queue := [][]string{{from}}
    
    for len(queue) > 0 {
        path := queue[0]
        queue = queue[1:]
        
        last := path[len(path)-1]
        
        if last == to {
            return path
        }
        
        if visited[last] {
            continue
        }
        visited[last] = true
        
        for _, dep := range a.findDirectDependents(last) {
            newPath := make([]string, len(path))
            copy(newPath, path)
            newPath = append(newPath, dep)
            queue = append(queue, newPath)
        }
    }
    
    return []string{}
}

func (a *ImpactAnalyzer) findRelatedTests(file string) []string {
    tests := []string{}
    
    base := strings.TrimSuffix(file, ".go")
    testFile := base + "_test.go"
    
    // Verifica se o arquivo de teste existe
    if _, ok := a.ctx.Files[testFile]; ok {
        tests = append(tests, testFile)
    }
    
    // Procura testes em pacotes relacionados
    for path := range a.ctx.Files {
        if strings.Contains(path, "_test.go") {
            // Verifica se o teste importa este arquivo
            if info, ok := a.ctx.Files[path]; ok {
                for _, imp := range info.Imports {
                    if strings.Contains(imp, file) {
                        tests = append(tests, path)
                        break
                    }
                }
            }
        }
    }
    
    return tests
}

func (a *ImpactAnalyzer) findDockerServices(file string) []string {
    services := []string{}
    
    // Procura por docker-compose.yml
    dockerPath := filepath.Join(a.ctx.Path, "docker-compose.yml")
    if _, err := os.Stat(dockerPath); err == nil {
        // Simplificado: em produção, faria parse do YAML
        services = append(services, "api", "worker") // exemplo
    }
    
    return services
}

func (a *ImpactAnalyzer) findRelatedMemories(file string) []*models.Memory {
    memories := []*models.Memory{}
    
    for _, m := range a.ctx.Memories {
        for _, f := range m.Files {
            if f == file {
                memories = append(memories, m)
                break
            }
        }
    }
    
    return memories
}

func (a *ImpactAnalyzer) findRelatedDecisions(file string) []*models.Decision {
    decisions := []*models.Decision{}
    
    for _, d := range a.ctx.Decisions {
        for _, f := range d.Files {
            if f == file {
                decisions = append(decisions, d)
                break
            }
        }
    }
    
    return decisions
}

func (a *ImpactAnalyzer) analyzeHistory(file string) (*FileHistory, error) {
    // Usa git para analisar histórico
    cmd := exec.Command("git", "log", "--pretty=format:%H|%an|%ad", "--date=iso", file)
    cmd.Dir = a.ctx.Path
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    history := &FileHistory{
        Authors: []string{},
    }
    
    lines := strings.Split(string(output), "\n")
    authors := make(map[string]bool)
    bugCount := 0
    
    for _, line := range lines {
        if line == "" {
            continue
        }
        
        parts := strings.Split(line, "|")
        if len(parts) >= 3 {
            history.Commits++
            authors[parts[1]] = true
            
            // Conta bugs por palavras-chave
            // Nota: em produção, isso seria mais sofisticado
            if strings.Contains(strings.ToLower(line), "fix") ||
               strings.Contains(strings.ToLower(line), "bug") {
                bugCount++
            }
        }
    }
    
    for author := range authors {
        history.Authors = append(history.Authors, author)
    }
    
    history.BugCount = bugCount
    history.ChangeFreq = float64(history.Commits) / 30.0 // commits por dia (aproximado)
    
    return history, nil
}

func (a *ImpactAnalyzer) calculateRiskScore(result *ImpactResult) int {
    score := 0
    
    // Quantidade de dependentes diretos
    score += len(result.DirectImpact) * 2
    
    // Quantidade de dependentes indiretos
    score += len(result.IndirectImpact)
    
    // Presença de memórias (bugs passados)
    score += len(result.Memories) * 3
    
    // Histórico de mudanças
    if result.History != nil {
        score += result.History.BugCount * 4
        if result.History.ChangeFreq > 1.0 {
            score += 5 // arquivo muda muito
        }
    }
    
    // Serviços Docker afetados
    score += len(result.DockerServices) * 3
    
    return score
}

func (a *ImpactAnalyzer) riskLevel(score int) string {
    switch {
    case score >= 20:
        return "critical"
    case score >= 10:
        return "high"
    case score >= 5:
        return "medium"
    default:
        return "low"
    }
}

func (a *ImpactAnalyzer) PrintResults(results map[string]*ImpactResult) {
    for file, result := range results {
        a.printResult(file, result)
    }
}

func (a *ImpactAnalyzer) printResult(file string, result *ImpactResult) {
    // Cabeçalho
    fmt.Println()
    fmt.Println(strings.Repeat("═", 80))
    
    // Risco com cor
    var riskColor func(format string, a ...interface{}) string
    switch result.Risk {
    case "critical":
        riskColor = color.RedString
    case "high":
        riskColor = color.YellowString
    case "medium":
        riskColor = color.CyanString
    default:
        riskColor = color.GreenString
    }
    
    fmt.Printf("📊 %s\n", color.WhiteString(file))
    fmt.Printf("   Risco: %s (score: %d)\n", riskColor(strings.ToUpper(result.Risk)), result.RiskScore)
    fmt.Println(strings.Repeat("─", 80))

    // Impacto direto
    if len(result.DirectImpact) > 0 {
        fmt.Printf("\n%s\n", color.YellowString("📁 IMPACTO DIRETO:"))
        for _, item := range result.DirectImpact {
            fmt.Printf("  ├─ %s", color.WhiteString(item.File))
            if item.Package != "" {
                fmt.Printf(" (%s)", color.CyanString(item.Package))
            }
            fmt.Println()
            for _, reason := range item.Reasons {
                fmt.Printf("  │  └─ %s\n", color.BlackString(reason))
            }
        }
    }

    // Impacto indireto
    if len(result.IndirectImpact) > 0 && a.options.Recursive {
        fmt.Printf("\n%s\n", color.YellowString("🔄 IMPACTO INDIRETO:"))
        for _, item := range result.IndirectImpact {
            indent := strings.Repeat("  ", item.Depth)
            fmt.Printf("%s├─ [%d] %s\n", indent, item.Depth, color.WhiteString(item.File))
        }
    }

    // Testes
    if len(result.Tests) > 0 {
        fmt.Printf("\n%s\n", color.GreenString("🧪 TESTES RELACIONADOS:"))
        for _, test := range result.Tests {
            fmt.Printf("  ├─ %s\n", test)
        }
    }

    // Docker
    if len(result.DockerServices) > 0 {
        fmt.Printf("\n%s\n", color.BlueString("🐳 SERVIÇOS DOCKER:"))
        for _, svc := range result.DockerServices {
            fmt.Printf("  ├─ %s\n", svc)
        }
    }

    // Memórias
    if len(result.Memories) > 0 && a.options.ShowMemories {
        fmt.Printf("\n%s\n", color.RedString("📝 MEMÓRIAS RELACIONADAS:"))
        for _, m := range result.Memories {
            fmt.Printf("  ├─ [%s] %s (%s)\n", 
                colorForSeverity(m.Severity),
                m.Title,
                m.CreatedAt.Format("2006-01-02"))
        }
    }

    // Decisões
    if len(result.Decisions) > 0 && a.options.ShowDecisions {
        fmt.Printf("\n%s\n", color.BlueString("⚖️ DECISÕES RELACIONADAS:"))
        for _, d := range result.Decisions {
            fmt.Printf("  ├─ %s (%s)\n", d.Title, d.Date.Format("2006-01-02"))
        }
    }

    // Histórico
    if result.History != nil && a.options.ShowHistory {
        fmt.Printf("\n%s\n", color.CyanString("📜 HISTÓRICO:"))
        fmt.Printf("  ├─ Commits: %d\n", result.History.Commits)
        fmt.Printf("  ├─ Autores: %s\n", strings.Join(result.History.Authors, ", "))
        fmt.Printf("  ├─ Bugs passados: %d\n", result.History.BugCount)
        fmt.Printf("  └─ Frequência: %.2f commits/dia\n", result.History.ChangeFreq)
    }

    // Recomendações baseadas no risco
    fmt.Printf("\n%s\n", color.MagentaString("💡 RECOMENDAÇÕES:"))
    switch result.Risk {
    case "critical":
        fmt.Println("  ├─ 🚨 Necessária revisão de arquitetura")
        fmt.Println("  ├─ 👥 Envolver time inteiro")
        fmt.Println("  └─ ⏸️  Considere pausar para planejamento")
    case "high":
        fmt.Println("  ├─ ⚠️  Code review obrigatório")
        fmt.Println("  ├─ 🧪 Testes extensivos necessários")
        fmt.Println("  └─ 📋 Documentar decisões")
    case "medium":
        fmt.Println("  ├─ ✅ Code review recomendado")
        fmt.Println("  └─ 🧪 Testes unitários suficientes")
    default:
        fmt.Println("  ├─ 👍 Baixo risco, pode prosseguir")
        fmt.Println("  └─ 🔍 Ainda assim, revise as mudanças")
    }
    
    fmt.Println(strings.Repeat("═", 80))
}

// Função auxiliar para cor da severidade
func colorForSeverity(severity string) string {
    switch severity {
    case "critical":
        return color.RedString("●")
    case "high":
        return color.YellowString("●")
    case "medium":
        return color.CyanString("●")
    default:
        return color.GreenString("●")
    }
}