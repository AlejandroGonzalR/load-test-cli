/*
Copyright © 2020 ALEJANDRO GONZALEZ alejandrogonzalr@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type options struct {
	Concurrency int
	Method      string
	Headers     string
	Body        string
}

var (
	cfgFile    string
	ClientOpts options

	// RootCmd represents the base command when called without any subcommands
	RootCmd = &cobra.Command{
		Use:   "load-test-cli",
		Short: "Tracing HTTP GET request latency",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), args[0]+"\n")

			u, err := url.ParseRequestURI(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			ExecRequest(u.String())
			return nil
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.load-test-cli.yaml)")
	RootCmd.Flags().IntVarP(&ClientOpts.Concurrency, "concurrency", "c", 1, "Number of requests to perform")
	RootCmd.Flags().StringVarP(&ClientOpts.Method, "method", "m", "GET", "Method to url")
	RootCmd.Flags().StringVarP(&ClientOpts.Headers, "headers", "H", "", "Send a header in request")
	RootCmd.Flags().StringVarP(&ClientOpts.Body, "data", "", "", "Send data body")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".load-test-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".load-test-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
