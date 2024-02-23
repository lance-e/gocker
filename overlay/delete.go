package overlay

import (
	"log"
	"os"
	"os/exec"
	"path"
	
)

func DeleteWorkSpace(rootURL string,volume string){
	UmountVolume(rootURL,volume)


	DeleteMerged(rootURL)
	
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
	if err := os.RemoveAll(rootURL);err != nil{
		log.Println("remove merged directory failed")
		return
	}
}
