package cmd

import (
	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed help/filtering.md
var filteringHelp string

func filteringTopic() *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "filtering",
		Short: "Pattern based filtering",
		Long:  "Pattern based filtering\n\n" + filteringHelp,
	}
}

//go:embed help/regions.md
var regionsHelp string

func regionsTopic() *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "regions",
		Short: "Handling file regions",
		Long:  "Handling file regions\n\n" + regionsHelp,
	}
}

//go:embed help/metadata.md
var metadataHelp string

func metadataTopic() *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "metadata",
		Short: "Code block metadata",
		Long:  "Code block metadata\n\n" + metadataHelp,
	}
}

//go:embed help/outline.md
var outlineHelp string

func outlineTopic() *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "outline",
		Short: "Embedding the file structure",
		Long:  "Embedding the file structure\n\n" + outlineHelp,
	}
}

//go:embed help/invisible.md
var invisibleHelp string

func invisibleTopic() *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "invisible",
		Short: "Invisible code blocks",
		Long:  "Invisible code blocks\n\n" + invisibleHelp,
	}
}
