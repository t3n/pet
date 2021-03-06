package cmd

import (
	"fmt"
	"io/ioutil"

	petSync "github.com/knqyf263/pet/sync"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit snippet file",
	Long:  `Edit snippet file (default: opened by vim)`,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("color", cmd.Flags().Lookup("color"))
		viper.BindPFlag("query", cmd.Flags().Lookup("query"))
	},
	RunE: edit,
}

func edit(cmd *cobra.Command, args []string) (err error) {
	editor := viper.GetString("general.editor")
	snippetFile := viper.GetString("general.snippetfile")

	var options []string
	if viper.GetString("query") != "" {
		options = append(options, fmt.Sprintf("--query %s", viper.GetString("query")))
	}

	if len(viper.GetStringSlice("general.snippetdirs")) > 0 {
		snippetFile, err = selectFile(options)
		if err != nil {
			return err
		}
	}

	if snippetFile == "" {
		return errors.New("No sippet file seleted")
	}

	// file content before editing
	before := fileContent(snippetFile)

	err = editFile(editor, snippetFile)
	if err != nil {
		return
	}

	// file content after editing
	after := fileContent(snippetFile)

	// return if same file content
	if before == after {
		return nil
	}

	if viper.GetBool("gist.auto_sync") {
		return petSync.AutoSync()
	}

	return nil
}

func fileContent(fname string) string {
	data, _ := ioutil.ReadFile(fname)
	return string(data)
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().BoolP("color", "", false, `Enable colorized output (only fzf)`)
	editCmd.Flags().StringP("query", "q", "", `Initial value for query`)
}
