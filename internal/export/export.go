package export

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Benjamin-Connelly/lookit/internal/index"
)

// Format identifies the export output format.
type Format int

const (
	FormatHTML Format = iota
	FormatPDF
)

// Options configures the export operation.
type Options struct {
	Format    Format
	OutputDir string
	Files     []string // specific files, or empty for all markdown
}

// Export converts markdown files to the specified output format.
func Export(idx *index.Index, opts Options) error {
	files := opts.Files
	if len(files) == 0 {
		for _, e := range idx.MarkdownFiles() {
			files = append(files, e.RelPath)
		}
	}

	if len(files) == 0 {
		return fmt.Errorf("no markdown files found")
	}

	outputDir := opts.OutputDir
	if outputDir == "" {
		outputDir = "lookit-export"
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	for _, file := range files {
		if err := exportFile(idx, file, outputDir, opts.Format); err != nil {
			return fmt.Errorf("exporting %s: %w", file, err)
		}
	}

	fmt.Printf("Exported %d files to %s\n", len(files), outputDir)
	return nil
}

func exportFile(idx *index.Index, relPath, outputDir string, format Format) error {
	entry := idx.Lookup(relPath)
	if entry == nil {
		return fmt.Errorf("file not found in index: %s", relPath)
	}

	source, err := os.ReadFile(entry.Path)
	if err != nil {
		return err
	}

	var outExt string
	switch format {
	case FormatHTML:
		outExt = ".html"
	case FormatPDF:
		outExt = ".pdf"
	}

	outName := replaceExt(filepath.Base(relPath), outExt)
	outPath := filepath.Join(outputDir, outName)

	// Placeholder: actual rendering via Goldmark
	return os.WriteFile(outPath, source, 0o644)
}

func replaceExt(name, newExt string) string {
	ext := filepath.Ext(name)
	return name[:len(name)-len(ext)] + newExt
}
