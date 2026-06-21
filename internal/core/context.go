package core

import (
    "crypto/rand"
    "fmt"
    "os"
    "path/filepath"
    "time"

    "github.com/ctx/ctx/pkg/models"
    "gopkg.in/yaml.v3"
)

type Context struct {
    Version     string                  `yaml:"version"`
    Name        string                  `yaml:"name"`
    Language    string                  `yaml:"language"`
    Module      string                  `yaml:"module,omitempty"`
    CreatedAt   time.Time               `yaml:"created_at"`
    UpdatedAt   time.Time               `yaml:"updated_at"`
    Files       map[string]*FileInfo     `yaml:"files"`
    Memories    []*models.Memory         `yaml:"memories"`
    Decisions   []*models.Decision       `yaml:"decisions"`
    Dependencies map[string][]string     `yaml:"dependencies"`
    Tags        map[string][]string      `yaml:"tags"`
    Path        string                   `yaml:"-"`
}

type FileInfo struct {
    Imports     []string    `yaml:"imports"`
    Functions   []string    `yaml:"functions"`
    Calls       []string    `yaml:"calls"`
    Package     string      `yaml:"package"`
    Size        int64       `yaml:"size"`
    ModifiedAt  time.Time   `yaml:"modified_at"`
    Hash        string      `yaml:"hash"`
    Memories    []string    `yaml:"memories,omitempty"`
    Decisions   []string    `yaml:"decisions,omitempty"`
}

func LoadContext(path string) (*Context, error) {
    ctxPath := filepath.Join(path, ".ctx.yaml")
    data, err := os.ReadFile(ctxPath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, nil
        }
        return nil, err
    }

    var ctx Context
    if err := yaml.Unmarshal(data, &ctx); err != nil {
        return nil, err
    }
    ctx.Path = path
    return &ctx, nil
}

func (c *Context) Save() error {
    c.UpdatedAt = time.Now()
    data, err := yaml.Marshal(c)
    if err != nil {
        return err
    }

    ctxPath := filepath.Join(c.Path, ".ctx.yaml")
    return os.WriteFile(ctxPath, data, 0644)
}

func (c *Context) AddMemory(memory *models.Memory) {
    memory.ID = generateID()
    memory.CreatedAt = time.Now()
    memory.UpdatedAt = time.Now()
    c.Memories = append(c.Memories, memory)
    
    // Atualiza referências nos arquivos
    for _, file := range memory.Files {
        if info, ok := c.Files[file]; ok {
            info.Memories = append(info.Memories, memory.ID)
        }
    }
}

func (c *Context) FindMemories(tags []string, files []string) []*models.Memory {
    var result []*models.Memory
    
    for _, memory := range c.Memories {
        // Filtra por tags
        if len(tags) > 0 {
            tagMatch := false
            for _, tag := range tags {
                if contains(memory.Tags, tag) {
                    tagMatch = true
                    break
                }
            }
            if !tagMatch {
                continue
            }
        }
        
        // Filtra por arquivos
        if len(files) > 0 {
            fileMatch := false
            for _, file := range files {
                if contains(memory.Files, file) {
                    fileMatch = true
                    break
                }
            }
            if !fileMatch {
                continue
            }
        }
        
        result = append(result, memory)
    }
    
    return result
}

func generateID() string {
    b := make([]byte, 16)
    rand.Read(b)
    return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}