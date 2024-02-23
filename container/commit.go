package container

import (
	"fmt"
	"log"
	"os/exec"
)


func CommitContainer (contaienrName string,imageName string){
	merged:=fmt.Sprintf(MergedPath,contaienrName)
	if _,err := exec.Command("tar","-czf","/root/"+imageName+".tar","-C",merged,".").CombinedOutput();err != nil{
		log.Println("commit failed")
	}
}