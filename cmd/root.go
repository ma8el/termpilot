package cmd

import (
	"fmt"
	"log"
	"os"
	"termpilot/db"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "termpilot",
		Short: "Termpilot is a terminal based AI agent",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := db.InitDB(); err != nil {
				log.Fatalf("Failed to initialize database: %v", err)
			}
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.termpilot.yaml)")
	rootCmd.PersistentFlags().String("model", "llama3.2", "model to use")
	rootCmd.PersistentFlags().String("base-url", "http://localhost", "base url")
	rootCmd.PersistentFlags().String("port", "11434", "port")
	rootCmd.PersistentFlags().String("version", "v1", "version")

	viper.BindPFlag("model", rootCmd.PersistentFlags().Lookup("model"))
	viper.BindPFlag("base-url", rootCmd.PersistentFlags().Lookup("base-url"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("version", rootCmd.PersistentFlags().Lookup("version"))

	rootCmd.AddCommand(chatCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to get user home directory: %v", err)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".termpilot")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Failed to read config file: %v", err)
		}
	}
}
