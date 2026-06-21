package utils

import (
    "bytes"
    "fmt"
    "os/exec"
    "strings"
    "time"
)

type GitCommit struct {
    Hash    string
    Author  string
    Email   string
    Date    time.Time
    Message string
}

type GitHistory struct {
    Commits     []*GitCommit
    Authors     []string
    BugCount    int
    ChangeFreq  float64
}

type GitBlameLine struct {
    LineNumber int
    CommitHash string
    Author     string
    Date       time.Time
    Content    string
}

// GetGitUser retorna o usuário configurado no git
func GetGitUser() (string, error) {
    cmd := exec.Command("git", "config", "user.name")
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(output)), nil
}

// GetGitBranch retorna a branch atual
func GetGitBranch() (string, error) {
    cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(output)), nil
}

// GetGitHash retorna o hash do último commit
func GetGitHash() (string, error) {
    cmd := exec.Command("git", "rev-parse", "HEAD")
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(output)), nil
}

// GetFileHistory retorna o histórico de um arquivo
func GetFileHistory(repoPath, filePath, since string) (*GitHistory, error) {
    args := []string{"log", "--pretty=format:%H|%an|%ae|%ad|%s", "--date=iso"}
    
    if since != "" {
        args = append(args, "--since", since)
    }
    
    args = append(args, "--", filePath)
    
    cmd := exec.Command("git", args...)
    cmd.Dir = repoPath
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    history := &GitHistory{
        Commits: []*GitCommit{},
        Authors: []string{},
    }

    authors := make(map[string]bool)
    lines := strings.Split(string(output), "\n")

    for _, line := range lines {
        if line == "" {
            continue
        }

        parts := strings.Split(line, "|")
        if len(parts) >= 5 {
            date, _ := time.Parse("2006-01-02 15:04:05 -0700", parts[3])
            
            commit := &GitCommit{
                Hash:    parts[0],
                Author:  parts[1],
                Email:   parts[2],
                Date:    date,
                Message: parts[4],
            }
            
            history.Commits = append(history.Commits, commit)
            
            if !authors[commit.Author] {
                authors[commit.Author] = true
                history.Authors = append(history.Authors, commit.Author)
            }

            // Conta commits de bug (simples heurística)
            msg := strings.ToLower(commit.Message)
            if strings.Contains(msg, "fix") || 
               strings.Contains(msg, "bug") || 
               strings.Contains(msg, "hotfix") ||
               strings.Contains(msg, "issue") {
                history.BugCount++
            }
        }
    }

    // Calcula frequência de mudanças
    if len(history.Commits) > 1 {
        first := history.Commits[len(history.Commits)-1].Date
        last := history.Commits[0].Date
        days := last.Sub(first).Hours() / 24
        if days > 0 {
            history.ChangeFreq = float64(len(history.Commits)) / days
        }
    }

    return history, nil
}

// GetFileBlame retorna o blame de um arquivo
func GetFileBlame(repoPath, filePath string) ([]*GitBlameLine, error) {
    cmd := exec.Command("git", "blame", "--line-porcelain", filePath)
    cmd.Dir = repoPath
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    lines := []*GitBlameLine{}
    currentLine := &GitBlameLine{}
    lineNumber := 1

    for _, line := range bytes.Split(output, []byte("\n")) {
        if len(line) == 0 {
            continue
        }

        parts := strings.SplitN(string(line), " ", 2)
        
        switch parts[0] {
        case "author":
            currentLine.Author = parts[1]
        case "author-time":
            timestamp, _ := time.Parse("2006-01-02 15:04:05 -0700", parts[1])
            currentLine.Date = timestamp
        case "filename":
            // Ignora
        default:
            if !strings.Contains(parts[0], "-") {
                // Linha de conteúdo
                if strings.HasPrefix(parts[0], "previous") {
                    continue
                }
                
                currentLine.LineNumber = lineNumber
                currentLine.Content = parts[1]
                lines = append(lines, currentLine)
                
                currentLine = &GitBlameLine{
                    CommitHash: parts[0],
                }
                lineNumber++
            }
        }
    }

    return lines, nil
}

// GetCommitInfo retorna informações de um commit específico
func GetCommitInfo(repoPath, hash string) (*GitCommit, error) {
    cmd := exec.Command("git", "show", "--pretty=format:%H|%an|%ae|%ad|%s", "--date=iso", hash)
    cmd.Dir = repoPath
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    parts := strings.Split(string(output), "|")
    if len(parts) < 5 {
        return nil, fmt.Errorf("invalid commit format")
    }

    date, _ := time.Parse("2006-01-02 15:04:05 -0700", parts[3])
    
    return &GitCommit{
        Hash:    parts[0],
        Author:  parts[1],
        Email:   parts[2],
        Date:    date,
        Message: parts[4],
    }, nil
}

// GetCommitsBetween retorna commits entre dois hashes
func GetCommitsBetween(repoPath, from, to string) ([]*GitCommit, error) {
    cmd := exec.Command("git", "log", "--pretty=format:%H|%an|%ae|%ad|%s", "--date=iso", fmt.Sprintf("%s..%s", from, to))
    cmd.Dir = repoPath
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    commits := []*GitCommit{}
    lines := strings.Split(string(output), "\n")

    for _, line := range lines {
        if line == "" {
            continue
        }

        parts := strings.Split(line, "|")
        if len(parts) >= 5 {
            date, _ := time.Parse("2006-01-02 15:04:05 -0700", parts[3])
            
            commits = append(commits, &GitCommit{
                Hash:    parts[0],
                Author:  parts[1],
                Email:   parts[2],
                Date:    date,
                Message: parts[4],
            })
        }
    }

    return commits, nil
}

// GetStagedFiles retorna arquivos staged para commit
func GetStagedFiles(repoPath string) ([]string, error) {
    cmd := exec.Command("git", "diff", "--cached", "--name-only")
    cmd.Dir = repoPath
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    files := []string{}
    for _, line := range strings.Split(string(output), "\n") {
        if line != "" {
            files = append(files, line)
        }
    }

    return files, nil
}

// GetModifiedFiles retorna arquivos modificados (não staged)
func GetModifiedFiles(repoPath string) ([]string, error) {
    cmd := exec.Command("git", "diff", "--name-only")
    cmd.Dir = repoPath
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    files := []string{}
    for _, line := range strings.Split(string(output), "\n") {
        if line != "" {
            files = append(files, line)
        }
    }

    return files, nil
}

// HasUncommittedChanges verifica se há mudanças não commitadas
func HasUncommittedChanges(repoPath string) (bool, error) {
    cmd := exec.Command("git", "status", "--porcelain")
    cmd.Dir = repoPath
    
    output, err := cmd.Output()
    if err != nil {
        return false, err
    }

    return len(strings.TrimSpace(string(output))) > 0, nil
}

// GetCurrentBranch retorna a branch atual
func GetCurrentBranch(repoPath string) (string, error) {
    return GetGitBranch()
}

// GetRemoteURL retorna a URL do remote
func GetRemoteURL(repoPath, remote string) (string, error) {
    if remote == "" {
        remote = "origin"
    }
    
    cmd := exec.Command("git", "config", "--get", fmt.Sprintf("remote.%s.url", remote))
    cmd.Dir = repoPath
    
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }

    return strings.TrimSpace(string(output)), nil
}

// CreateCommitMessage cria uma mensagem de commit baseada nas mudanças
func CreateCommitMessage(repoPath string, files []string) (string, error) {
    if len(files) == 0 {
        return "", fmt.Errorf("no files to commit")
    }

    // Simplificado: em produção, usaria IA ou template
    if len(files) == 1 {
        return fmt.Sprintf("Update %s", files[0]), nil
    }
    
    return fmt.Sprintf("Update %d files", len(files)), nil
}

// TagExists verifica se uma tag existe
func TagExists(repoPath, tag string) (bool, error) {
    cmd := exec.Command("git", "tag", "-l", tag)
    cmd.Dir = repoPath
    
    output, err := cmd.Output()
    if err != nil {
        return false, err
    }

    return strings.TrimSpace(string(output)) == tag, nil
}

// CreateTag cria uma nova tag
func CreateTag(repoPath, tag, message string) error {
    cmd := exec.Command("git", "tag", "-a", tag, "-m", message)
    cmd.Dir = repoPath
    
    return cmd.Run()
}

// GetFileAuthor retorna o principal autor de um arquivo
func GetFileAuthor(repoPath, filePath string) (string, error) {
    cmd := exec.Command("git", "log", "--pretty=format:%an", filePath)
    cmd.Dir = repoPath
    
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }

    authors := make(map[string]int)
    for _, line := range strings.Split(string(output), "\n") {
        if line != "" {
            authors[line]++
        }
    }

    var topAuthor string
    var topCount int
    
    for author, count := range authors {
        if count > topCount {
            topCount = count
            topAuthor = author
        }
    }

    return topAuthor, nil
}

// GetFileBirth retorna a data de criação do arquivo
func GetFileBirth(repoPath, filePath string) (time.Time, error) {
    cmd := exec.Command("git", "log", "--pretty=format:%ad", "--date=iso", "--diff-filter=A", filePath)
    cmd.Dir = repoPath
    
    output, err := cmd.Output()
    if err != nil {
        return time.Time{}, err
    }

    dateStr := strings.TrimSpace(string(output))
    if dateStr == "" {
        return time.Time{}, fmt.Errorf("file not found in git history")
    }

    return time.Parse("2006-01-02 15:04:05 -0700", dateStr)
}

// GetContributors retorna todos os contribuidores do repositório
func GetContributors(repoPath string) ([]string, error) {
    cmd := exec.Command("git", "shortlog", "-sn", "--all")
    cmd.Dir = repoPath
    
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    contributors := []string{}
    for _, line := range strings.Split(string(output), "\n") {
        if line != "" {
            // Formato: " 123\tNome do Autor"
            parts := strings.Split(line, "\t")
            if len(parts) == 2 {
                contributors = append(contributors, strings.TrimSpace(parts[1]))
            }
        }
    }

    return contributors, nil
}

// IsGitRepo verifica se o diretório é um repositório git
func IsGitRepo(path string) bool {
    cmd := exec.Command("git", "rev-parse", "--git-dir")
    cmd.Dir = path
    return cmd.Run() == nil
}

// GetRepoName retorna o nome do repositório
func GetRepoName(path string) (string, error) {
    remote, err := GetRemoteURL(path, "origin")
    if err != nil {
        // Fallback para nome do diretório
        parts := strings.Split(path, "/")
        return parts[len(parts)-1], nil
    }

    // Extrai nome do remote
    parts := strings.Split(remote, "/")
    last := parts[len(parts)-1]
    return strings.TrimSuffix(last, ".git"), nil
}