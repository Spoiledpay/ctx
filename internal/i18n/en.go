package i18n

func init() {
    messages["en"] = map[string]string{
        // Comandos gerais
        "ctx_short": "CTX - Code Context & Memory",
        "ctx_long": "CTX is a source code diary that stores decisions, bugs, and context about your project",
        
        // Init
        "init_short": "Initialize CTX in a directory",
        "init_success": "✅ CTX initialized successfully",
        
        // Scan
        "scan_short": "Scan code and update context",
        "scan_start": "Scanning project in %s...",
        
        // Memory
        "memory_short": "Manage memories",
        "memory_add_short": "Add a new memory",
        
        // Decision
        "decision_short": "Manage architecture decisions",
        
        // Explain
        "explain_short": "Get context about files or functions",
        
        // Impact
        "impact_short": "Analyze impact of changes",
        
        // Grep
        "grep_short": "Search through context",
        
        // Version
        "version_short": "Show version information",
        
        // Erros
        "error_no_ctx": "❌ No CTX context found. Run 'ctx init' first",
        "error_not_found": "❌ %s not found",
    }
}