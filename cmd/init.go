/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gocker/container"
	"log"

	"github.com/spf13/cobra"
)

// private command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init",
	Long: "Init container process run user's process in container . Do not call it outside.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("init command begin")
		container.InitProcess()
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
