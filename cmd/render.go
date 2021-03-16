package cmd

import (
	"os"
	"strings"

	"github.com/TakeScoop/steve/pkg/helm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "steve",
	Short: "SONAR stevedore for deploying apps to Kubernetes",
	RunE: func(cmd *cobra.Command, args []string) error {
		hc, err := helm.New()
		if err != nil {
			return err
		}
		hc.Install(args)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	f := rootCmd.PersistentFlags()
	_ = viper.BindPFlags(f)
}
