package commands

import (
    "fmt"
    "runtime"
    "strings"

    "github.com/ctx/ctx/internal/i18n"
    "github.com/fatih/color"
    "github.com/spf13/cobra"
)

var (
    Version = "dev"
    Commit  = "none"
    Date    = "unknown"
)

func NewVersionCmd() *cobra.Command {
    var verbose bool

    cmd := &cobra.Command{
        Use:   "version",
        Short: i18n.T("version_short"),
        Run: func(cmd *cobra.Command, args []string) {
            if verbose {
                printVerboseVersion()
            } else {
                fmt.Printf("ctx version %s\n", Version)
            }
        },
    }

    cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, i18n.T("version_flag_verbose"))
    return cmd
}

func printVerboseVersion() {
    color.Cyan("📦 CTX - Code Context & Memory")
    fmt.Println(strings.Repeat("─", 40))
    fmt.Printf("Version:    %s\n", color.GreenString(Version))
    fmt.Printf("Commit:     %s\n", color.YellowString(Commit))
    fmt.Printf("Built:      %s\n", Date)
    fmt.Printf("Go:         %s\n", runtime.Version())
    fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
}