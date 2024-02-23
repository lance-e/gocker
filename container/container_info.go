package container

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"

	"time"
)

type ContainerInfo struct {
	Pid           string `json:"pid"`
	ContainerName string `json:"container_name"`
	ContainerId   string ` json:"container_id"`
	CreateAtTime  string `json:"create_at_time"`
	Command       string `json:"command"`
	Status        string `json:"status"`
}

var (
	RUNING              = "runing"
	STOP                = "stopped"
	EXIT                = "exited"
	DefaultInfoLocation = "/var/run/gocker/%s/"
	ConfigName          = "config.json"
)

func RecordContainerInformation(containerName string, containerPid int, command []string) (string, error) {
	id := randStringBytes(10)
	createTime := time.Now().Format("2006-01-02 15:04:05")
	if containerName == "" {
		containerName = id
	}
	cmd := strings.Join(command, " ")
	info := &ContainerInfo{
		Pid:           strconv.Itoa(containerPid),
		ContainerName: containerName,
		ContainerId:   id,
		CreateAtTime:  createTime,
		Status:        RUNING,
		Command:       cmd,
	}
	infoByte, err := json.Marshal(info)
	if err != nil {
		log.Println("can't marshal the information")
		return "", err
	}

	location := fmt.Sprintf(DefaultInfoLocation, containerName)
	log.Println("the location is ", location)
	if err := os.MkdirAll(location, 0622); err != nil {
		return "", err
	}
	file, err := os.Create(path.Join(location, ConfigName))
	if err != nil {
		log.Println("can't create the config file")
		return "", err
	}
	defer file.Close()

	_, err = file.Write(infoByte)
	if err != nil {
		log.Println("can't write the byte of information into the target file")
		return "", err
	}
	log.Println("record the informaiton successful!")
	return containerName, nil
}
func randStringBytes(n int) string {
	num := "01234567890123456789"
	by := make([]byte, n)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i, _ := range by {
		by[i] = num[r.Intn(n)]
	}
	return string(by)
}

func DeleteContainerInfo(containerName string){
	location := fmt.Sprintf(DefaultInfoLocation,containerName)
	dirPath := location+ConfigName
	if err := os.RemoveAll(dirPath);err != nil{
		log.Println("can't delete the config file")
	}

}