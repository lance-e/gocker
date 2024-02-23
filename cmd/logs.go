package cmd

import (
	"gocker/container"
	"log"

	"github.com/spf13/cobra"
)

var LogsCmd = &cobra.Command{
	Use: "logs",
	Short: "logs the detach container",
	Long: "use the logs commond to show the detach container print out ",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args)<1{
			log.Println("missing the container name")
			return
		}
		container.ShowLogs(args[0])
	},
}
func init(){
	rootCmd.AddCommand(LogsCmd)
}