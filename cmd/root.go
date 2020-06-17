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
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-httpstat"
)

var cfgFile string
var requestNum int

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "load-test-cli",
	Short: "Tracing HTTP GET request latency",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintf(cmd.OutOrStdout(), args[0]+"\n")
		req, _ := cmd.Flags().GetInt("request")

		u, err := url.ParseRequestURI(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		execRequest(u.String(), req)
		return nil
	},
}

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
	RootCmd.Flags().IntVarP(&requestNum, "request", "n", 1, "Number of sequential HTTP request")
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

func execRequest(url string, request int) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create a httpstat powered context
	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	// Send request by default HTTP client
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	res.Body.Close()

	// Show the results
	fmt.Printf("DNS lookup: %d ms\n", int(result.DNSLookup/time.Millisecond))
	fmt.Printf("TCP connection: %d ms\n", int(result.TCPConnection/time.Millisecond))
	fmt.Printf("TLS handshake: %d ms\n", int(result.TLSHandshake/time.Millisecond))
	fmt.Printf("Server processing: %d ms\n", int(result.ServerProcessing/time.Millisecond))
	fmt.Printf("Content transfer: %d ms\n", int(result.ContentTransfer(time.Now())/time.Millisecond))
}
