package models

import "time"

type Memory struct {
    ID          string    `yaml:"id"`
    Title       string    `yaml:"title"`
    Description string    `yaml:"description"`
    Type        string    `yaml:"type"` // bug, decision, lesson, warning
    Severity    string    `yaml:"severity"` // low, medium, high, critical
    
    // Contexto
    Files       []string  `yaml:"files"`
    Lines       []int     `yaml:"lines,omitempty"`
    Commit      string    `yaml:"commit,omitempty"`
    Branch      string    `yaml:"branch,omitempty"`
    
    // Solução
    Solution    string    `yaml:"solution,omitempty"`
    Workaround  string    `yaml:"workaround,omitempty"`
    
    // Metadados
    Tags        []string  `yaml:"tags"`
    Author      string    `yaml:"author"`
    Reviewers   []string  `yaml:"reviewers,omitempty"`
    CreatedAt   time.Time `yaml:"created_at"`
    UpdatedAt   time.Time `yaml:"updated_at"`
    ResolvedAt  *time.Time `yaml:"resolved_at,omitempty"`
    
    // Links
    Links       []string  `yaml:"links,omitempty"` // issues, PRs, docs
}

type Decision struct {
    ID          string    `yaml:"id"`
    Title       string    `yaml:"title"`
    Description string    `yaml:"description"`
    
    // A decisão
    Decision    string    `yaml:"decision"`
    Rationale   string    `yaml:"rationale"`
    Alternatives []string `yaml:"alternatives,omitempty"`
    
    // Contexto
    Files       []string  `yaml:"files"`
    Participants []string `yaml:"participants"`
    Date        time.Time `yaml:"date"`
    
    // Status
    Status      string    `yaml:"status"` // proposed, accepted, deprecated, rejected
    DeprecatedBy string   `yaml:"deprecated_by,omitempty"`
    
    Tags        []string  `yaml:"tags"`
    Links       []string  `yaml:"links,omitempty"`
}

type Tag struct {
    Name        string   `yaml:"name"`
    Description string   `yaml:"description"`
    Color       string   `yaml:"color"`
    Count       int      `yaml:"count"`
}