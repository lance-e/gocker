package cmd

import (
	"gocker/container"
	"log"

	"github.com/spf13/cobra"
)

var RemoveCmd = &cobra.Command{
	Use: "rm",
	Short: "remove the container ",
	Long: "remove the target container ,but the container must be stoped",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <1 {
			log.Println("missing the container name")
			return
		}
		container.RemoveContainer(args[0])
	},
}
func init(){
	rootCmd.AddCommand(RemoveCmd)
}