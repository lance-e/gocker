package cmd

import (
	"gocker/container"

	"github.com/spf13/cobra"
)

var PsCmd = &cobra.Command{
	Use: "ps",
	Short: "list the container ",
	Long: "use the ps command to show the all container's status",
	Run: func(cmd *cobra.Command, args []string) {
		container.ListContainers()
	},
}
func init(){
	rootCmd.AddCommand(PsCmd)
}