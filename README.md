markdown
# CTX Documentation

Welcome to the CTX documentation! CTX is a source code diary that stores decisions, bugs, and context about your project, giving AI assistants (and your future self) the memory they need.

## 📚 Table of Contents

- [Getting Started](#getting-started)
- [Commands Reference](commands.md)
- [Memory System](memory.md)
- [Decisions (ADR)](decisions.md)
- [Explain Command](explain.md)
- [Configuration](configuration.md)
- [Internationalization](i18n.md)
- [Git Integration](git.md)
- [AI Integration](ai.md)
- [Examples](../examples/)
- [FAQ](faq.md)
- [Troubleshooting](troubleshooting.md)

## 🚀 Getting Started

### Installation

#### Windows (PowerShell)
```powershell
# Using go install
go install github.com/ctx/ctx@latest

# Or download from releases
Invoke-WebRequest -Uri "https://github.com/Spoiledpay/ctx/releases/latest/download/ctx-windows-amd64.exe" -OutFile ctx.exe
Windows (Command Prompt)
cmd
go install github.com/Spoiledpay/ctx/ctx@latest
Verify Installation
bash
ctx version
Quick Start
bash
# Initialize CTX in your project
cd your-project
ctx init

# Scan your code
ctx scan

# Add a memory (bug, decision, etc.)
ctx memory add --title "Fixed login timeout" --type bug --files auth.go

# Explain a file
ctx explain auth.go

# Check impact before changing
ctx impact auth.go
🎯 Core Concepts
What is CTX?
CTX is a source code diary that lives alongside your code. It stores:

Memories - Bugs, solutions, lessons learned

Decisions - Architecture decisions and their rationale

Context - Why code is written the way it is

History - Git history analyzed and enriched

Why CTX?
AI assistants (Cursor, Copilot, etc.) are brilliant but amnesic. They can read your code perfectly but have no idea:

Why a decision was made

What bugs happened last month

Who knows about this part of the code

What was already tried and failed

CTX fixes this by being the long-term memory for your project.

📖 Command Overview
Command	Description
ctx init	Initialize CTX in a directory
ctx scan	Scan code and update context
ctx memory	Manage memories (bugs, lessons)
ctx decision	Manage architecture decisions
ctx explain	Get context about files/functions
ctx impact	Analyze impact of changes
ctx grep	Search through context
ctx review	Code review with context
💡 Examples
Daily Workflow
bash
# Before changing code, understand it
ctx explain auth.go

# Check if there are known issues
ctx memory list --files auth.go

# See related decisions
ctx decision list --files auth.go

# Analyze impact
ctx impact auth.go

# Make your changes...

# Record what you learned
ctx memory add --title "Auth optimization" --type lesson --files auth.go
Team Onboarding
bash
# New developer joins
ctx explain . --detailed

# Shows:
# - Project structure
# - Recent decisions
# - Known bugs
# - Who to ask about what
Code Review
bash
# Before submitting PR
ctx review

# Shows:
# - Files changed
# - Related memories
# - Risk assessment
# - Tests to run
🤝 Integrating with AI
CTX is designed to be read by AI assistants. Add this to your Cursor/Copilot prompt:

text
When helping with code in this project, please check .ctx.yaml for context about decisions, bugs, and history.
The AI will then have access to:

Past decisions and their rationale

Known bugs and solutions

Code ownership and expertise

Risk areas to be careful with

🌍 Internationalization
CTX supports multiple languages:

English (en) - Default
Portuguese (pt)
Spanish (es)
Russian (ru)
Chinese (zh)
Japanese (jp)

Set language with:

bash
ctx --lang pt explain auth.go
Or via environment variable:

bash
set LANG=pt_BR
ctx explain auth.go
🔧 Configuration
CTX is configured via .ctx.yaml in your project root. See configuration.md for details.

📊 Project Health
CTX can help you understand your project's health:

bash
# Show project statistics
ctx explain .

# Find buggy files
ctx memory list --type bug --output table

# See recent decisions
ctx decision list --limit 10

# Analyze risk areas
ctx impact $(ctx grep --files-only --tags critical)
🆘 Getting Help
ctx help - Show general help

ctx help [command] - Help for specific command

FAQ - Frequently asked questions

Troubleshooting - Common issues and solutions

GitHub Issues - Report bugs and feature requests

📝 License
MIT License - see LICENSE for details.

text

---

### `docs/commands.md`

```markdown
# CTX Commands Reference

## Table of Contents

- [Global Flags](#global-flags)
- [init](#init)
- [scan](#scan)
- [memory](#memory)
- [decision](#decision)
- [explain](#explain)
- [impact](#impact)
- [grep](#grep)
- [review](#review)
- [version](#version)

## Global Flags

These flags can be used with any command:

| Flag | Description |
|------|-------------|
| `--lang, -l` | Language (auto, en, pt, es, ru, zh, jp) |
| `--verbose, -v` | Verbose output |
| `--help, -h` | Show help |
| `--version` | Show version |

## init

Initialize CTX in a directory.

```bash
ctx init [directory] [flags]
Flags:

Flag	Description
--force, -f	Force reinitialization
Examples:

bash
# Initialize current directory
ctx init

# Initialize specific directory
ctx init ./my-project

# Force reinit
ctx init --force
scan
Scan code and update context.

bash
ctx scan [directory] [flags]
Flags:

Flag	Description
--recursive, -r	Scan recursively (default true)
--include-tests, -t	Include test files
--verbose, -v	Verbose output
--force, -f	Force rescan
--watch, -w	Watch for changes
--interval, -i	Watch interval in seconds (default 5)
--output, -o	Output format (text, json)
Examples:

bash
# Scan current directory
ctx scan

# Scan with watch mode
ctx scan --watch

# Scan and output JSON
ctx scan --output json
memory
Manage memories (bugs, lessons, warnings).

Subcommands
memory add
Add a new memory.

bash
ctx memory add [flags]
Flags:

Flag	Description
--title, -t	Memory title (required)
--description, -d	Memory description
--type	Type (bug, decision, lesson, warning)
--severity, -s	Severity (low, medium, high, critical)
--files, -f	Related files
--lines, -l	Related line numbers
--commit, -c	Git commit hash
--branch, -b	Git branch
--solution	Solution description
--workaround	Workaround description
--tags	Tags for categorization
--author, -a	Author (defaults to git user)
--reviewers, -r	Reviewers
--links	Related links (issues, PRs)
--from	Import from file
Examples:

bash
# Add a bug
ctx memory add --title "Login timeout" --type bug --severity high --files auth.go

# Add a lesson
ctx memory add --title "Always check indexes" --type lesson --tags performance

# Import from file
ctx memory add --from memory.json
memory list
List memories.

bash
ctx memory list [flags]
Flags:

Flag	Description
--type, -t	Filter by type
--severity, -s	Filter by severity
--tags	Filter by tags
--files, -f	Filter by files
--output, -o	Output format (text, table, json)
--limit, -l	Limit results
--sort	Sort by (date, title, severity)
Examples:

bash
# List all memories
ctx memory list

# List bugs
ctx memory list --type bug

# List with JSON output
ctx memory list --output json
memory show
Show memory details.

bash
ctx memory show [id]
Example:

bash
ctx memory show abc-123
memory update
Update a memory.

bash
ctx memory update [id] [flags]
Flags: Same as memory add

memory delete
Delete a memory.

bash
ctx memory delete [id] [flags]
Flags:

Flag	Description
--force, -f	Force without confirmation
memory search
Search memories.

bash
ctx memory search [query]
Example:

bash
ctx memory search "login timeout"
memory export
Export memories.

bash
ctx memory export [flags]
Flags:

Flag	Description
--format, -f	Format (json, yaml, csv)
--output, -o	Output file (default: stdout)
memory import
Import memories.

bash
ctx memory import [file] [flags]
Flags:

Flag	Description
--format, -f	Format (json, yaml)
decision
Manage architecture decisions (ADRs).

Subcommands
decision add
Add a new decision.

bash
ctx decision add [flags]
Flags:

Flag	Description
--title, -t	Decision title (required)
--description, -d	Description
--decision	The decision made (required)
--rationale, -r	Rationale (required)
--alternatives, -a	Alternatives considered
--files, -f	Related files
--participants, -p	Participants
--date	Decision date (default: today)
--status, -s	Status (proposed, accepted, deprecated, rejected)
--tags	Tags
--links, -l	Related links
--from	Import from file
Examples:

bash
# Add a decision
ctx decision add --title "Database Choice" --decision "PostgreSQL" --rationale "ACID compliance"

# Add with alternatives
ctx decision add --title "API Design" --decision "REST" --alternatives "GraphQL,gRPC"
decision list
List decisions.

bash
ctx decision list [flags]
Flags:

Flag	Description
--status, -s	Filter by status
--tags, -t	Filter by tags
--files, -f	Filter by files
--date, -d	Filter by date
--output, -o	Output format (text, table, json)
--limit, -l	Limit results
--sort	Sort by (date, title, status)
decision show
Show decision details.

bash
ctx decision show [id]
decision update
Update a decision.

bash
ctx decision update [id] [flags]
decision delete
Delete a decision.

bash
ctx decision delete [id] [flags]
decision deprecate
Deprecate a decision.

bash
ctx decision deprecate [id] [flags]
Flags:

Flag	Description
--reason, -r	Deprecation reason
--replacement	Replacement decision ID
decision reactivate
Reactivate a deprecated decision.

bash
ctx decision reactivate [id]
decision export
Export decisions.

bash
ctx decision export [flags]
Flags:

Flag	Description
--format, -f	Format (json, yaml, markdown)
--output, -o	Output file
--status, -s	Filter by status
--tags, -t	Filter by tags
decision import
Import decisions.

bash
ctx decision import [file] [flags]
Flags:

Flag	Description
--format, -f	Format (json, yaml)
--merge, -m	Merge with existing
explain
Get context about files or functions.

bash
ctx explain [target] [flags]
Flags:

Flag	Description
--detailed, -d	Detailed output
--bugs, -b	Include bug information (default true)
--decisions	Include decisions (default true)
--history, -H	Include git history
--deps	Include dependencies
--tests, -t	Include tests
--format, -f	Output format (text, json)
--output, -o	Output file
--from-git, -g	Get info from git
--since	History since (e.g., "30d", "1y")
Examples:

bash
# Explain a file
ctx explain auth.go

# Explain with history
ctx explain auth.go --history

# Explain a function
ctx explain Login

# Explain entire project
ctx explain .
impact
Analyze impact of changes.

bash
ctx impact [files...] [flags]
Flags:

Flag	Description
--recursive, -r	Include indirect impact (default true)
--depth, -d	Maximum depth (default 3)
--include-tests, -t	Include tests
--show-memories, -m	Show related memories (default true)
--show-decisions	Show related decisions (default true)
--show-history, -H	Show git history
--format, -f	Output format (text, json)
Examples:

bash
# Check impact of changing a file
ctx impact auth.go

# Check multiple files
ctx impact auth.go user.go

# With JSON output
ctx impact auth.go --format json
grep
Search through context.

bash
ctx grep [pattern] [flags]
Flags:

Flag	Description
--files-only, -l	Show only file names
--memories, -m	Search memories (default true)
--decisions, -d	Search decisions (default true)
--code, -c	Search code
--tags, -t	Filter by tags
--type	Filter by type
--output, -o	Output format (text, json)
Examples:

bash
# Search for "login"
ctx grep login

# Search only in files
ctx grep login --files-only

# Search with tags
ctx grep performance --tags critical
review
Code review with context.

bash
ctx review [flags]
Flags:

Flag	Description
--staged, -s	Review staged changes
--unstaged, -u	Review unstaged changes
--pr, -p	Review PR (GitHub PR number)
--format, -f	Output format (text, json, markdown)
--output, -o	Output file
Examples:

bash
# Review staged changes
ctx review --staged

# Review PR
ctx review --pr 123

# Output as markdown
ctx review --staged --format markdown --output review.md
version
Show version information.

bash
ctx version [flags]
Flags:

Flag	Description
--verbose, -v	Detailed version info
Example:

bash
ctx version --verbose
text

---

### `docs/memory.md`

```markdown
# CTX Memory System

The memory system is the heart of CTX. It stores everything you want to remember about your code: bugs, solutions, lessons learned, warnings, and more.

## 📝 Memory Structure

Each memory has the following fields:

```yaml
id: "mem-12345678"                # Unique identifier
title: "Login timeout bug"         # Short description
description: "Users experiencing..." # Detailed description
type: "bug"                        # Type: bug, lesson, warning, note
severity: "high"                   # Severity: low, medium, high, critical

# Context
files: 
  - "internal/service/auth.go"     # Related files
  - "internal/db/user.go"
lines: [42, 87]                    # Related line numbers
commit: "abc123def"                 # Git commit hash
branch: "main"                      # Git branch

# Solution
solution: "Added database index..." # How it was solved
workaround: "Restart service..."    # Temporary workaround

# Metadata
tags: ["performance", "database"]  # Tags for categorization
author: "@joao"                     # Who created it
reviewers: ["@maria", "@pedro"]     # Who reviewed it
created_at: "2026-02-20T10:00:00Z" # When created
updated_at: "2026-02-21T15:30:00Z" # Last updated
resolved_at: "2026-02-21T15:30:00Z" # When resolved

# Links
links: 
  - "https://github.com/org/repo/issues/123"
  - "https://discord.com/channels/..."
🎯 Memory Types
Bug
Record bugs and their fixes.

bash
ctx memory add --type bug \
  --title "API timeout in production" \
  --severity critical \
  --files "internal/handler/api.go" \
  --solution "Added pagination" \
  --tags "performance,api"
Lesson
Things you learned that everyone should know.

bash
ctx memory add --type lesson \
  --title "Always check database indexes" \
  --description "Slow queries often mean missing indexes" \
  --tags "database,performance"
Warning
Dangerous areas or things to be careful with.

bash
ctx memory add --type warning \
  --title "Race condition in worker pool" \
  --severity high \
  --files "internal/worker/pool.go" \
  --description "This code is not thread-safe!" \
  --workaround "Use mutex or single-threaded"
Note
General notes about the code.

bash
ctx memory add --type note \
  --title "Why we use SHA256 instead of bcrypt" \
  --files "internal/service/auth.go" \
  --description "Legacy system compatibility" \
  --tags "security,architecture"
🏷️ Tags System
Tags help organize and find memories. Use them consistently:

Common Tags
performance - Performance issues

security - Security concerns

database - Database related

api - API related

ui - User interface

testing - Testing issues

deployment - Deployment issues

refactoring - Code cleanup

documentation - Documentation needs

critical - Critical issues

Tag Conventions
Use lowercase

Use single words or hyphenated

Be consistent across the project

🔍 Finding Memories
By Type
bash
# All bugs
ctx memory list --type bug

# All warnings
ctx memory list --type warning
By Severity
bash
# Critical issues only
ctx memory list --severity critical

# High and critical
ctx memory list --severity high --severity critical
By Tags
bash
# Performance related
ctx memory list --tags performance

# Multiple tags
ctx memory list --tags performance,database
By Files
bash
# Memories related to auth.go
ctx memory list --files auth.go

# Multiple files
ctx memory list --files auth.go,user.go
Search
bash
# Full text search
ctx memory search "login timeout"

# Search with filters
ctx memory search "database" --tags performance
📊 Memory Statistics
bash
# Count by type
ctx memory list --output stats

# Most buggy files
ctx memory list --type bug --sort count

# Recent issues
ctx memory list --sort date --limit 10
💡 Best Practices
When to Add a Memory
✅ When you fix a bug
✅ When you learn something non-obvious
✅ When you make a significant decision
✅ When you encounter a tricky part of the code
✅ When you add a workaround
✅ When you discover a dangerous area

Memory Writing Tips
Be specific - Include exact file names and line numbers
Explain why - Not just what happened, but why
Include context - What were the conditions?
Document solutions - How was it fixed?
Add tags - Make it searchable
Link to resources - Issues, PRs, discussions

Example: Good Memory
yaml
title: "Race condition in user cache"
type: bug
severity: high
files: ["internal/cache/user.go"]
lines: [42, 87]
description: |
  Concurrent access to user cache causing panic when multiple 
  goroutines try to write to the same map.
solution: |
  Replaced map with sync.Map and added mutex for complex operations.
  Also added tests for concurrent access.
tags: ["concurrency", "cache", "critical"]
links: ["https://github.com/org/repo/issues/456"]
Example: Poor Memory
yaml
title: "Bug fix"
type: bug
description: "Fixed the thing"
# Missing: files, severity, solution, tags
🔄 Memory Lifecycle
Create - Memory is added when issue discovered

Update - Add more info as you learn
Resolve - Mark as resolved when fixed
Reference - Used in explain/impact commands
Archive - Old memories can be archived

🤖 AI Integration
Memories are automatically used by AI when:
Explaining files (ctx explain)
Analyzing impact (ctx impact)
Reviewing code (ctx review)

The AI will warn you about related memories before you make similar mistakes.

bash
ctx explain auth.go
# Output includes:
# ⚠️ Warning: This file has 3 related bugs:
#   - Login timeout (2026-02-15) - Fixed with index
#   - Session leak (2026-01-10) - Fixed with cleanup
#   - Race condition (2026-01-05) - Fixed with mutex
📈 Memory Health
Good signs:

✅ Growing number of memories
✅ Decreasing bugs over time
✅ Active updates and resolutions
✅ Cross-referenced with decisions

Warning signs:

❌ No memories (you're not learning)
❌ Only bugs (no lessons/warnings)
❌ Never resolved (issues linger)
❌ Outdated (no recent updates)

text

---

### `docs/decisions.md`

```markdown
# CTX Decisions System (ADR)

The decisions system implements Architecture Decision Records (ADRs) - a way to document important architectural decisions and their context.

## 📖 What is an ADR?

An Architecture Decision Record (ADR) is a document that captures an important architectural decision, along with its context and consequences.

Each ADR is:
- **Atomic** - One decision per record
- **Immutable** - Decisions can be deprecated but not changed
- **Structured** - Consistent format for easy reading

## 📝 Decision Structure

```yaml
id: "adr-001"                       # Unique identifier
title: "Database Selection"           # Short title
description: "Choosing primary..."    # Detailed context

# The Decision
decision: "Use PostgreSQL"            # What was decided
rationale: "ACID compliance..."       # Why
alternatives:                         # What was considered
  - "MySQL - Rejected due to..."
  - "MongoDB - Rejected due to..."

# Context
files: 
  - "internal/db/connection.go"      # Related files
  - "docs/architecture.md"
participants:                         # Who decided
  - "@joao"
  - "@maria"
  - "@pedro"
date: "2026-01-15"                    # When

# Status
status: "accepted"                    # proposed, accepted, deprecated, rejected
deprecated_by: "adr-042"               # If deprecated, what replaced it

# Metadata
tags: ["database", "infrastructure"]  # Tags
links:                                 # Related links
  - "https://github.com/org/repo/issues/123"
🎯 Decision Status
Proposed
A decision being considered.

bash
ctx decision add --status proposed \
  --title "Use Kubernetes" \
  --decision "Migrate to k8s" \
  --rationale "Better scaling" \
  --participants "@joao,@maria"
Accepted
A decision that has been made and is in effect.

bash
ctx decision update adr-001 --status accepted
Deprecated
A decision that is no longer valid, replaced by another.

bash
ctx decision deprecate adr-001 \
  --reason "New infrastructure" \
  --replacement adr-042
Rejected
A decision that was considered but not adopted.

bash
ctx decision add --status rejected \
  --title "Use MongoDB" \
  --decision "Rejected" \
  --rationale "Consistency requirements" \
  --alternatives "PostgreSQL"
🏷️ Common Decision Tags
architecture - Overall architecture

database - Database choices

api - API design

security - Security decisions

scalability - Scaling decisions

performance - Performance tradeoffs

tech-stack - Technology choices

process - Development process

deployment - Deployment strategy

testing - Testing approach

🔍 Finding Decisions
By Status
bash
# All accepted decisions
ctx decision list --status accepted

# Proposed decisions needing review
ctx decision list --status proposed
By Tags
bash
# Database decisions
ctx decision list --tags database

# Multiple tags
ctx decision list --tags database,performance
By Date
bash
# Recent decisions (last 30 days)
ctx decision list --since 30d

# Decisions from specific date
ctx decision list --date 2026-01-15
By File
bash
# Decisions affecting auth.go
ctx decision list --files auth.go
📊 Decision Workflow
1. Proposal
bash
ctx decision add \
  --title "Use GraphQL" \
  --status proposed \
  --decision "Adopt GraphQL for new APIs" \
  --rationale "Flexible queries for frontend" \
  --alternatives "REST (less flexible), gRPC (overkill)" \
  --participants "@joao,@maria" \
  --tags "api,architecture"
2. Discussion
Team discusses the proposal. Update with more context:

bash
ctx decision update adr-023 \
  --description "Added performance considerations"
3. Decision
bash
# If accepted
ctx decision update adr-023 --status accepted

# If rejected
ctx decision update adr-023 --status rejected
4. Implementation
Link to implementation files:

bash
ctx decision update adr-023 \
  --files "internal/graphql/schema.go" \
  --links "https://github.com/org/repo/pull/456"
5. Deprecation
If a decision is later replaced:

bash
ctx decision deprecate adr-023 \
  --reason "Moved to BFF pattern" \
  --replacement adr-089
💡 Best Practices
When to Create an ADR
✅ Choosing a technology (database, framework, etc.)
✅ Defining architecture patterns
✅ Making significant design decisions
✅ Changing existing decisions
✅ Rejecting alternatives
✅ Documenting non-negotiable requirements

Good ADR Example
yaml
title: "Cache Strategy for User Sessions"
status: accepted
date: 2026-02-15

decision: "Use Redis for session storage"
rationale: |
  - Need sub-millisecond latency
  - Session data must persist through restarts
  - Team already familiar with Redis
  - Built-in expiration mechanisms

alternatives:
  - "In-memory cache - Rejected: Doesn't survive restarts"
  - "PostgreSQL - Rejected: Too slow for this use case"
  - "Memcached - Rejected: Less features than Redis"

consequences:
  - Add Redis to infrastructure
  - Update deployment scripts
  - Add cache invalidation strategy

files: ["internal/session/redis_store.go"]
participants: ["@joao", "@maria", "@pedro"]
tags: ["caching", "performance", "infrastructure"]
Poor ADR Example
yaml
title: "Use Redis"
decision: "Use Redis"
# Missing: rationale, alternatives, consequences
🔄 Integration with Memories
Decisions and memories work together:

Decisions capture intentional choices

Memories capture learned experiences

bash
# Decision leading to a bug
ctx decision add --title "Use in-memory cache" --status accepted
# Later...
ctx memory add --type bug --title "Cache causes memory leak" --decision adr-001
📈 Decision Health
Good signs:

✅ Clear rationale for each decision
✅ Alternatives considered and documented
✅ Active decisions are implemented
✅ Deprecated decisions have replacements

Warning signs:

❌ No decisions (random architecture)
❌ Decisions without rationale
❌ Never revisited or updated
❌ Contradictory decisions

🤖 AI Integration
Decisions are automatically used by AI when:

Explaining why code is written a certain way

Suggesting alternatives

Warning about deprecated approaches

bash
ctx explain auth.go
# Output includes:
# ⚖️ Decisions affecting this file:
#   - adr-001: Use JWT tokens (2026-01-15)
#   - adr-023: Redis for sessions (2026-02-01)
📋 ADR Template
markdown
# [ID] Title

## Status
[proposed | accepted | deprecated | rejected]

## Context
What is the issue that we're seeing that is motivating this decision or change?

## Decision
What is the change that we're proposing and/or doing?

## Rationale
Why are we doing this? What are the benefits?

## Consequences
What becomes easier or more difficult to do?

## Alternatives
What alternatives did we consider?

## Participants
Who was involved in this decision?

## References
- Links to related issues, PRs, documents
Use the template with:

bash
ctx decision add --from template.md
text

---

Estes arquivos completam a documentação principal. Quer que eu gere também:

- `docs/configuration.md` - Configuração detalhada
- `docs/explain.md` - Detalhes do comando explain
- `docs/git.md` - Integração com Git
- `docs/ai.md` - Integração com IAs
- `docs/faq.md` - Perguntas frequentes
