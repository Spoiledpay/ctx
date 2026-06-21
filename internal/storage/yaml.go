package storage

import (
    "io/ioutil"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

type Storage struct {
    path string
}

func NewStorage(path string) *Storage {
    return &Storage{path: path}
}

func (s *Storage) Load(data interface{}) error {
    content, err := ioutil.ReadFile(s.path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil // Arquivo não existe ainda
        }
        return err
    }

    return yaml.Unmarshal(content, data)
}

func (s *Storage) Save(data interface{}) error {
    content, err := yaml.Marshal(data)
    if err != nil {
        return err
    }

    // Cria diretório se não existir
    dir := filepath.Dir(s.path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    return ioutil.WriteFile(s.path, content, 0644)
}

func (s *Storage) Exists() bool {
    _, err := os.Stat(s.path)
    return err == nil
}

func (s *Storage) Backup() error {
    if !s.Exists() {
        return nil
    }

    backupPath := s.path + ".backup"
    input, err := ioutil.ReadFile(s.path)
    if err != nil {
        return err
    }

    return ioutil.WriteFile(backupPath, input, 0644)
}