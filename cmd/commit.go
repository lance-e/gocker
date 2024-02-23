package cmd

import (
	"gocker/container"
	"log"

	"github.com/spf13/cobra"
)


var CommitCmd = &cobra.Command{
	Use: "commit",
	Short: "commit a container into image",
	Long: "commit a container into image",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args)<2 {
			log.Println("miss the container name or the image name")
			return
		}
		container.CommitContainer(args[0],args[1])
	},
}

func init(){
	rootCmd.AddCommand(CommitCmd)
}