package container

import (
	"fmt"
	"gocker/overlay"
	"log"
	"os"
	"path"
)

func RemoveContainer(contaienrName string){
	info,err := getContainerInfoByName(contaienrName)
	if err != nil{
		log.Println("can't get the information ,error:",err)
		return
	}
	if info.Status != STOP{
		log.Println("can't remove the running container")
		return
	}
	
	dir := fmt.Sprintf(DefaultInfoLocation, contaienrName)
	if err = os.RemoveAll(dir);err != nil{
		log.Println("can't remove all the file and directory")
		return
	}
	overlay.DeleteWorkSpace(path.Join(RootUrl,info.ContainerName),info.Volume)
	log.Printf("remove %s successful",info.ContainerName)
	
}