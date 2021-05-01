package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/joyrex2001/kubedock/internal"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "kubedock",
	Short: "kubedock is a docker on kubernetes service.",
	Long:  ``,
	Run:   internal.Main,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().String("listen-addr", ":8080", "Webserver listen address")
	rootCmd.PersistentFlags().String("namespace", "default", "Namespace in which containers should be orchestrated")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose mode")
	rootCmd.PersistentFlags().BoolP("logrequest", "r", false, "Log requests and responses (can contain credentials)")
	viper.BindPFlag("generic.verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("generic.logrequest", rootCmd.PersistentFlags().Lookup("logrequest"))
	viper.BindPFlag("server.listen-addr", rootCmd.PersistentFlags().Lookup("listen-addr"))

	// kubeconfig
	if home := homeDir(); home != "" {
		rootCmd.PersistentFlags().String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		rootCmd.PersistentFlags().String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	viper.BindPFlag("kubernetes.kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))

	viper.BindPFlag("kubernetes.namespace", rootCmd.PersistentFlags().Lookup("namespace"))
	viper.BindEnv("server.listen-addr", "SERVER_LISTEN_ADDR")
	viper.BindEnv("kubernetes.namespace", "NAMESPACE")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
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
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		// fmt.Printf("not using config file: %s\n", err)
	} else {
		fmt.Printf("using config: %s\n", viper.ConfigFileUsed())
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
