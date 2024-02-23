package cmd

import (
	"gocker/container"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var ExecCmd = &cobra.Command{
	Use:   "exec",
	Short: "enter a container ",
	Long:  "use the exec command to enter the target container",
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv(container.ENV_EXEC_PID) != "" {
			log.Printf("pid callback is %d\n", os.Getpid())
			return
		}
		//格式为./gocker exec 容器名 后续命令
		if len(args) < 2 {
			log.Println("missing the contianer name or command ")
			return
		}
		container.ExecContainer(args[0],args[1:])
	},
}

func init() {
	rootCmd.AddCommand(ExecCmd)
}
