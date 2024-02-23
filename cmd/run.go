/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	
	"gocker/cgroup"
	"gocker/container"
	"log"

	"github.com/spf13/cobra"
)

var (
	tty bool
	t bool
	i bool
	resource cgroup.ResouceConfig
	volume string
	detach bool
	name string
	environment []string
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
			log.Println("missing container command")
			return
		}
		if i && t {
			tty = true
		}
		if tty && detach{
			log.Println("the tty and detach can't exist at the same time")
			return
		}
		var id = ""
		if name == ""{
			name = container.RandStringBytes(10)
			id =name
		}
		imageName := args[0]
		container.Run(tty, imageName,args[1:],&resource,volume,name,id,environment)
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
	runCmd.Flags().BoolVarP(&i, "interactive", "i", false, "Keep STDIN open even if not attached")
	runCmd.Flags().BoolVarP(&t, "tty", "t", false, "Allocate a pseudo-TTY")
	runCmd.Flags().StringVarP(&resource.MemoryLimit,"memory","m","1024m","set the memory limit ")
	runCmd.Flags().StringVarP(&resource.CpuSet,"cpuset","","0-2","the cpuset subsystem")
	runCmd.Flags().StringVarP(&resource.CpuShare,"cpushare","","1024","the cpu share subsystem")
	runCmd.Flags().StringVarP(&volume,"volume","v","","the data volume")
	runCmd.Flags().BoolVarP(&detach,"detach","d",false,"make the container detach")
	runCmd.Flags().StringVarP(&name,"name","","","set the container name ")
	runCmd.Flags().StringSliceVarP(&environment,"environment","e",[]string{},"set the environment variable")
}
