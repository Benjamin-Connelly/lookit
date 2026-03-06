package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lookit [path]",
	Short: "Dual-mode markdown navigator with inter-document link navigation",
	Long: `Lookit is a dual-mode markdown navigator that provides both TUI and web
interfaces for browsing code, markdown, and files. Features inter-document
link navigation with history, backlinks, and broken link detection.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root := "."
		if len(args) > 0 {
			root = args[0]
		}
		_ = root
		// Default: launch TUI mode
		fmt.Println("lookit: TUI mode (not yet implemented)")
		return nil
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve [path]",
	Short: "Start the web server",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root := "."
		if len(args) > 0 {
			root = args[0]
		}
		_ = root
		fmt.Println("lookit serve: web mode (not yet implemented)")
		return nil
	},
}

var catCmd = &cobra.Command{
	Use:   "cat <file>",
	Short: "Render markdown to terminal",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("lookit cat: %s (not yet implemented)\n", args[0])
		return nil
	},
}

var exportCmd = &cobra.Command{
	Use:   "export [path]",
	Short: "Export markdown files to HTML or PDF",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("lookit export: (not yet implemented)")
		return nil
	},
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check environment and diagnose issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("lookit doctor: (not yet implemented)")
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file path")
	rootCmd.PersistentFlags().String("theme", "auto", "color theme (light|dark|auto)")

	rootCmd.Flags().String("keymap", "default", "keybinding preset (default|vim|emacs)")

	serveCmd.Flags().IntP("port", "p", 7777, "server port")
	serveCmd.Flags().Bool("open", false, "open browser after starting")
	serveCmd.Flags().Bool("no-https", false, "disable HTTPS even if certs exist")

	exportCmd.Flags().StringP("format", "f", "html", "export format (html|pdf)")
	exportCmd.Flags().StringP("output", "o", "", "output directory")

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(catCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(doctorCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
