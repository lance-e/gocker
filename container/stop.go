package container

import (
	"encoding/json"
	"fmt"
	"strconv"
	"syscall"

	"log"
	"os"
)

func StopContainer(containerName string) {
	info, err := getContainerInfoByName(containerName)
	if err != nil {
		log.Println("can't get contaienr information by container name")
		return
	}
	pid, _ := strconv.Atoi(info.Pid)
	if err = syscall.Kill(pid, syscall.SIGTERM); err != nil {
		log.Println("can't send the kill sigt,error:", err)
		return
	}
	info.Status = STOP
	info.Pid = ""

	data, err := json.Marshal(info)
	if err != nil {
		log.Println("can't marshal the information ")
		return
	}

	dir := fmt.Sprintf(DefaultInfoLocation, containerName)
	allDir := dir + ConfigName
	if err = os.WriteFile(allDir, data, 0622);err != nil{
		log.Println("can't write the stoped container's information into the config file")
		return
	}
	log.Printf("%s stopping ...\n",info.ContainerName)
	
}
func getContainerInfoByName(name string) (*ContainerInfo, error) {
	dir := fmt.Sprintf(DefaultInfoLocation, name)
	allDir := dir + ConfigName
	data, err := os.ReadFile(allDir)
	if err != nil {
		log.Println("can't read the config informaiton from the file")
		return nil, err
	}
	var info = &ContainerInfo{}
	err = json.Unmarshal(data, info)
	if err != nil {
		log.Println("unmalshal the information failed")
		return nil, err
	}
	return info, nil
}
