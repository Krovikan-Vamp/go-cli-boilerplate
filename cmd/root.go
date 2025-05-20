/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"TODO/pkg/log"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	verbosity int

	outputPath string
	envPath    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zoom-recordings-vtt",
	Short: "Start the Zoom VTT Transcription Server",
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if outputPath == "" {
			outputPath = "." // default to current directory
		}
		err := log.InitLogger(outputPath, zerolog.Level(verbosity))
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to initialize logger")
			return
		}

		// Set log level based on verbosity
		switch {
		case verbosity >= 3:
			log.Logger = log.Logger.Level(zerolog.TraceLevel)
		case verbosity == 2:
			log.Logger = log.Logger.Level(zerolog.DebugLevel)
		case verbosity == 1:
			log.Logger = log.Logger.Level(zerolog.InfoLevel)
		default:
			log.Logger = log.Logger.Level(zerolog.InfoLevel)
		}

		// Load environment variables from .env file
		err = godotenv.Load(envPath)
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to load .env file. Please specify the path using the -e flag")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().CountVarP(&verbosity, "verbose", "v", "increase logging verbosity\n\t-vv for debug\n\t-vvv for trace")
	rootCmd.PersistentFlags().StringVarP(&envPath, "env", "e", ".env", "path to .env file")
	rootCmd.PersistentFlags().StringVarP(&outputPath, "output-path", "O", ".", "Path to write profiles and logs to")
}
