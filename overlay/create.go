package overlay

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
)

func NewWorkSpace(rootURL string, volume string,imageName string) {
	NewUpper(rootURL)
	NewLower(imageName)
	NewWork(rootURL)
	NewMerged(rootURL,imageName)

	
	MountVolume(rootURL,volume)
}

// 镜像层
func NewLower(imageName string) {
	lowerPath := path.Join("/root/overlay",imageName)
	_, err := os.Stat(lowerPath)
	if err == nil {
		log.Println("lower is nolmal")

	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(lowerPath, 0777)
		if err != nil {
			log.Println("can't create the lower,error:", err.Error())
			return
		}
		image:= path.Join("/root", imageName+".tar")
		fmt.Println(image)
		fmt.Println(lowerPath)
		if _, err := exec.Command("tar", "-xvf", image, "-C", lowerPath).CombinedOutput(); err != nil {
			log.Println("can't tar the target file")
			return
		}
	}
	if err != nil && !os.IsNotExist(err) {
		log.Println("can't judge the lower directory'state")
	}

}

// 容器层
func NewUpper(rootURL string) {
	upperPath := path.Join(rootURL, "upper")
	_, err := os.Stat(upperPath)
	if err == nil {
		log.Println("upper already exists")

	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(upperPath, 0777)
		if err != nil {
			log.Println("can't create the upper,error:", err.Error())
			return
		}
	}
	if err != nil && !os.IsNotExist(err) {
		log.Println("can't judge the upper directory's state")

	}

}
func NewWork(rootURL string) {
	workPath := path.Join(rootURL, "work")
	_, err := os.Stat(workPath)
	if err == nil {
		log.Println("work already exists")

	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(workPath, 0777)
		if err != nil {
			log.Println("can't create the work,error:", err.Error())
			return
		}
	}
	if err != nil && !os.IsNotExist(err) {
		log.Println("can't judge the work file's state")
	}

}
func NewMerged(rootURL string,imageName string) {
	mergedPath := path.Join(rootURL, "merged")
	_, err := os.Stat(mergedPath)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(mergedPath, 0777)
			if err != nil {
				log.Println("can't create the merged ,error:", err.Error())
				return
			}
		} else {
			log.Println("can't judge the merged file's state")
			return
		}
	}

	dir := "upperdir=" + path.Join(rootURL, "upper") + ",lowerdir=" + path.Join("/root/overlay",imageName) + ",workdir=" + path.Join(rootURL, "work")
	cmd := exec.Command("mount", "-o", dir, "-t", "overlay", "overlay", path.Join(rootURL, "merged"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println("overlayFS mount failed ,error :", err)
	}
}

