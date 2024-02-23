package container

import (
	"log"
	"os/exec"
)


func CommitContainer (imageName string){
	if _,err := exec.Command("tar","-czf","/root/"+imageName+".tar","-C","/root/overlay/merged",".").CombinedOutput();err != nil{
		log.Println("commit failed")
	}
}