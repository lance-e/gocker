package overlay

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func DeleteWorkSpace(rootURL string,volume string){
	if volume != ""{
		volumeSlice:=strings.Split(volume,":")
		if len(volumeSlice)==2  && volumeSlice[0] !=""&&volumeSlice[1] != ""{
			if err := VolumeUnmount(path.Join(rootURL,"merged",volumeSlice[1]));err != nil{
				log.Println("umount volume failed,error:",err)
				return
			}
		}
	}
	DeleteMerged(rootURL)
	DeleteUpper(rootURL)
	DeleteWork(rootURL)
}

func DeleteMerged(rootURL string){
	mergedPath := path.Join(rootURL,"merged")
	cmd := exec.Command("umount",mergedPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run();err != nil{
		log.Println("umount overlay failed")
		return
	}
	log.Println("umount overlay successful...")
	if err := os.RemoveAll(mergedPath);err != nil{
		log.Println("remove merged directory failed")
		return
	}
}
func DeleteUpper(rootURL string){
	upperPath := path.Join(rootURL,"upper")
	if err := os.RemoveAll(upperPath);err != nil{
		log.Println("remove upper diroctory failed")
		return
	}

}

func DeleteWork(rootURL string){
	workPath := path.Join(rootURL,"work")
	if err := os.RemoveAll(workPath);err != nil{
		log.Println("remove work diroctory failed")
		return
	}
}
func VolumeUnmount(target string)error{
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