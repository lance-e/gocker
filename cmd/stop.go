package cmd

import (
	"gocker/container"
	"log"

	"github.com/spf13/cobra"
)

var StopCmd = &cobra.Command{
	Use: "stop",
	Short: "stop the container ",
	Long: "use the stop command  to stop the target container ",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args)<1 {
			log.Println("missing the container name")
			return
		}
		container.StopContainer(args[0])
	},
}
func init(){
	rootCmd.AddCommand(StopCmd)
}