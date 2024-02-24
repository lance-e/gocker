package cmd

import (
	"gocker/network"
	"log"

	"github.com/spf13/cobra"
)

var NetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "config the network",
	Long:  "network command is to config the container's network",
}

var (
	subnet string
	driver string
)
var CreateNetworkCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new network",
	Long:  "use the create subcommand to create a new network in the container",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Println("missint the name of the network")
			return
		}
		if err := network.Init(); err != nil {
			return
		}

		network.CreateNetWork(driver, subnet, args[0])
	},
}
var ListNetworkCmd = &cobra.Command{
	Use:   "list",
	Short: "list all network",
	Long:  "list all network",
	Run: func(cmd *cobra.Command, args []string) {
		if err := network.Init();err != nil{
			return
		}
		network.ListNetWork()
	},
}
var RemoveNetworkCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove network",
	Long:  "use the rm command to remove the target network",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Println("missing the network name")
			return
		}
		if err := network.Init();err != nil{
			return
		}
		if err := network.DeleteNetWork(args[0]);err != nil{
			log.Println("can't delete the network,error:",err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(NetworkCmd)
	NetworkCmd.AddCommand(CreateNetworkCmd)
	NetworkCmd.AddCommand(ListNetworkCmd)
	NetworkCmd.AddCommand(RemoveNetworkCmd)

	CreateNetworkCmd.Flags().StringVarP(&driver, "driver", "", "bridge", "the network driver")
	CreateNetworkCmd.Flags().StringVarP(&subnet, "subnet", "", "","the subnet for your network")
}
