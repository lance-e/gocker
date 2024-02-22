package overlay

import (
	"log"
	"os"
	"os/exec"
	"path"
)

func DeleteWorkSpace(){
	DeleteMerged("/root/overlay")
	DeleteUpper("/root/overlay")
	DeleteWork("/root/overlay")
}

func DeleteMerged(rootURL string){
	mergedPath := path.Join(rootURL,"merged")
	cmd := exec.Command("umount",mergedPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run();err != nil{
		log.Println("umount overlay failed")
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
