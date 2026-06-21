package i18n

import (
    "sync"
)

var (
    currentLang string = "en"
    messages    = make(map[string]map[string]string)
    mu          sync.RWMutex
)

func init() {
    // Inglês (padrão)
    messages["en"] = map[string]string{
        "ctx_short":     "CTX - Code Context & Memory",
        "ctx_long":      "CTX is a source code diary that stores decisions, bugs, and context about your project",
        "flag_lang":     "Language (auto, en, pt, es, ru, zh, jp)",
        "flag_verbose":  "Verbose output",
        "error":         "❌ Error: %v",
        "warning":       "⚠️ Warning: %v",
        "success":       "✅ Success: %v",
        "info":          "ℹ️ %v",
        
        // Init command
        "init_start":    "Initializing CTX in %s...",
        "init_success":  "CTX initialized successfully",
        "init_exists":   "CTX already exists in this directory",
        
        // Scan command
        "scan_start":    "Scanning project...",
        "scan_files":    "Found %d Go files",
        "scan_deps":     "Analyzing dependencies...",
        "scan_complete": "Scan complete. Updated %d files",
        
        // Memory command
        "memory_add":    "Adding memory...",
        "memory_added":  "Memory added with ID: %s",
        "memory_list":   "Found %d memories",
        "memory_show":   "Memory: %s",
        
        // Impact command
        "impact_title":  "Impact Analysis for %s",
        "impact_files":  "Affected files:",
        "impact_tests":  "Tests to run:",
        "impact_docker": "Docker services affected:",
    }

    // Português
    messages["pt"] = map[string]string{
        "ctx_short":     "CTX - Contexto e Memória de Código",
        "ctx_long":      "CTX é um diário de código fonte que armazena decisões, bugs e contexto sobre seu projeto",
        "flag_lang":     "Idioma (auto, en, pt, es, ru, zh, jp)",
        "flag_verbose":  "Saída detalhada",
        "error":         "❌ Erro: %v",
        "warning":       "⚠️ Aviso: %v",
        "success":       "✅ Sucesso: %v",
        "info":          "ℹ️ %v",
        
        "init_start":    "Inicializando CTX em %s...",
        "init_success":  "CTX inicializado com sucesso",
        "init_exists":   "CTX já existe neste diretório",
        
        "scan_start":    "Escaneando projeto...",
        "scan_files":    "Encontrados %d arquivos Go",
        "scan_deps":     "Analisando dependências...",
        "scan_complete": "Scan completo. %d arquivos atualizados",
        
        "memory_add":    "Adicionando memória...",
        "memory_added":  "Memória adicionada com ID: %s",
        "memory_list":   "Encontradas %d memórias",
        "memory_show":   "Memória: %s",
        
        "impact_title":  "Análise de Impacto para %s",
        "impact_files":  "Arquivos afetados:",
        "impact_tests":  "Testes a rodar:",
        "impact_docker": "Serviços Docker afetados:",
    }

    // Russo
    messages["ru"] = map[string]string{
        "ctx_short":     "CTX - Контекст и память кода",
        "ctx_long":      "CTX - это дневник исходного кода, который хранит решения, ошибки и контекст о вашем проекте",
        "error":         "❌ Ошибка: %v",
        "warning":       "⚠️ Предупреждение: %v",
        "success":       "✅ Успешно: %v",
        "info":          "ℹ️ %v",
        // ... resto das traduções
    }

    // Chinês
    messages["zh"] = map[string]string{
        "ctx_short":     "CTX - 代码上下文和记忆",
        "ctx_long":      "CTX 是一个源代码日记，存储有关项目的决策、错误和上下文",
        "error":         "❌ 错误: %v",
        "warning":       "⚠️ 警告: %v",
        "success":       "✅ 成功: %v",
        "info":          "ℹ️ %v",
        // ... resto das traduções
    }
}

func T(key string) string {
    mu.RLock()
    defer mu.RUnlock()
    
    if msg, ok := messages[currentLang][key]; ok {
        return msg
    }
    // Fallback para inglês
    if msg, ok := messages["en"][key]; ok {
        return msg
    }
    return key
}

func SetLanguage(lang string) {
    mu.Lock()
    defer mu.Unlock()
    currentLang = lang
}