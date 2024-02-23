package overlay

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func MountVolume(rootURL string,volume string){
	if volume != "" {
		
		volumeslice := strings.Split(volume, ":")
		if len(volumeslice) == 2 && volumeslice[0] != "" && volumeslice[1] != "" {
			volumeMount(rootURL, volumeslice)
		} else {
			log.Println("the params are not correct")
		}
	}
}

func volumeMount(rootURL string, volume []string) {
	//创建宿主机目录
	if err := os.Mkdir(volume[0], 0755); err != nil {
		log.Println("the state of the ", volume[0], "in host is :", err.Error())
	}

	mergedPath := path.Join(rootURL, "merged")
	target := path.Join(mergedPath, volume[1])
	//创建容器目录
	if err := os.MkdirAll(target, 0755); err != nil {
		log.Println("the state of the ", volume[1], "in container is :", err.Error())
	}
	cmd := exec.Command("mount", "-o", "bind", volume[0], target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println("mount volume failed,error:", err.Error())
		return
	}
	log.Printf("mount %s to %s successful!\n", volume[0], target)
}
func UmountVolume(rootURL string,volume string){
	if volume != ""{
		volumeSlice:=strings.Split(volume,":")
		if len(volumeSlice)==2  && volumeSlice[0] !=""&&volumeSlice[1] != ""{
			if err := volumeUnmount(path.Join(rootURL,"merged",volumeSlice[1]));err != nil{
				log.Println("umount volume failed,error:",err)
				return
			}
		}
	}
}

func volumeUnmount(target string)error{
	cmd := exec.Command("umount",target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run();err != nil{
		log.Println("umount volume failed")
		return err
	}
	log.Println("umount volume successful")
	return nil
}