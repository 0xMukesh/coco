package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/0xmukesh/coco/internal/driver"
	"github.com/spf13/cobra"
)

var (
	outputFilePath string
	emitIr         bool
)

var rootCmd = &cobra.Command{
	Use:   "coco",
	Short: "a statically typed compiled scripting language",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

var buildCmd = &cobra.Command{
	Use:   "build [file]",
	Short: "Builds and outputs and executable binary",
	Run:   buildSource,
	Args:  cobra.ExactArgs(1),
}

var typeCheckCmd = &cobra.Command{
	Use:     "typecheck [file]",
	Aliases: []string{"tc"},
	Short:   "Statically typechecks source code",
	Run:     typeCheckSource,
	Args:    cobra.ExactArgs(1),
}

func Run() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	buildCmd.Flags().StringVarP(&outputFilePath, "output", "o", "", "path where the executable binary needs to be saved")
	buildCmd.Flags().BoolVarP(&emitIr, "emit-ir", "", false, "whether to emit llvm ir or not")
	rootCmd.AddCommand(buildCmd, typeCheckCmd)
}

func buildSource(cmd *cobra.Command, args []string) {
	sourceFilePath := args[0]
	d, err := driver.NewDriverFromFile(sourceFilePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := d.Pipeline(outputFilePath, emitIr); err != nil {
		log.Fatal(err)
	}
}

func typeCheckSource(cmd *cobra.Command, args []string) {
	sourceFilePath := args[0]
	d, err := driver.NewDriverFromFile(sourceFilePath)
	if err != nil {
		log.Fatal(err)
	}

	tks, err := d.Lex()
	if err != nil {
		log.Fatal(err)
	}

	ast, err := d.Parse(tks)
	if err != nil {
		log.Fatal(err)
	}

	if err := d.TypeCheck(ast); err != nil {
		log.Fatal(err)
	}
}
