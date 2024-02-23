package container

import (
	"fmt"
	"io"
	"log"
	"os"
)


var (
	ContainerLogFile = "container.logs"
)
func ShowLogs(containerName string){
	dir := fmt.Sprintf(DefaultInfoLocation,containerName)
	path := dir + ContainerLogFile
	file ,err := os.Open(path)
	if err != nil{
		log.Println("can't open the logs file")
		return
	}
	defer file.Close()
	data,err := io.ReadAll(file)
	if err != nil{
		log.Println("can't read from thelogs file")
		return
	}
	fmt.Fprint(os.Stdout,string(data))
}