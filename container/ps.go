package container

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

func ListContainers() {
	dir := fmt.Sprintf(DefaultInfoLocation, "")
	dir = dir[:len(dir)-1]
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Println("can't read the all directory")
		return

	}
	var containers []*ContainerInfo
	for _, file := range files {
		info, _ := file.Info()
		data, err := readInfoFromFile(info)
		if err != nil {
			log.Println("read container information failed")
			continue
		}
		containers = append(containers, data)
	}
	w:=tabwriter.NewWriter(os.Stdout,12,1,3,' ',0)
	fmt.Fprintf(w,"Container_Id\tContainer_Name\tPid\tStatus\tCommand\tCreate_At_Time\n")
	for _,con := range containers {
		fmt.Fprintf(w,"%s\t%s\t%s\t%s\t%s\t%s\t\n",con.ContainerId,con.ContainerName,con.Pid,con.Status,con.Command,con.CreateAtTime,
	)
	}
	if err := w.Flush();err != nil{
		log.Println("can't flush to stdout")
		return
	}
}
func readInfoFromFile(file os.FileInfo) (*ContainerInfo, error) {
	containerName := file.Name()
	dir := fmt.Sprintf(DefaultInfoLocation, containerName)
	dir = dir + ConfigName
	info, err := os.ReadFile(dir)
	if err != nil {
		log.Println("can't read the file")
		return nil, err
	}
	var data = &ContainerInfo{}
	err = json.Unmarshal(info, data)
	if err != nil {
		log.Println("can't unmarshal the information ")
		return nil, err
	}
	return data, nil
}
