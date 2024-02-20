/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"gocker/container"
	"gocker/cgroup"

	"github.com/spf13/cobra"
)

var (
	tty bool
	t bool
	i bool
	resource cgroup.ResouceConfig
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("missing container command")
			return
		}
		if i && t {
			tty = true
		}
		container.Run(tty, args,&resource)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(initCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.Flags().BoolVarP(&i, "interactive", "i", true, "Keep STDIN open even if not attached")
	runCmd.Flags().BoolVarP(&t, "tty", "t", true, "Allocate a pseudo-TTY")
	runCmd.Flags().StringVarP(&resource.MemoryLimit,"memory","m","1024","set the memory limit ")
	runCmd.Flags().StringVarP(&resource.CpuSet,"cpuset","","","the cpuset subsystem")
	runCmd.Flags().StringVarP(&resource.CpuShare,"cpushare","","","the cpu share subsystem")
}
