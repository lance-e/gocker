package overlay

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func NewWorkSpace(rootURL string, volume string) {
	_, err := os.Stat(rootURL)
	if os.IsNotExist(err) {
		if err := os.Mkdir(rootURL, 0777); err != nil {
			log.Println("can't create a new work space,error:", err.Error())
			return
		}

	}
	NewUpper(rootURL)
	NewLower(rootURL)
	NewWork(rootURL)
	NewMerged(rootURL)
	if volume != "" {
		volumeslice := strings.Split(volume, ":")
		if len(volumeslice) == 2 && volumeslice[0] != "" && volumeslice[1] != "" {
			VolumeMount(rootURL, volumeslice)
		} else {
			log.Println("the params are not correct")
		}
	}
}

// 镜像层
func NewLower(rootURL string) {
	lowerPath := path.Join(rootURL, "lower")
	_, err := os.Stat(lowerPath)
	if err == nil {
		log.Println("lower is nolmal")

	}
	if os.IsNotExist(err) {
		err = os.Mkdir(lowerPath, 0777)
		if err != nil {
			log.Println("can't create the lower,error:", err.Error())
			return
		}
		if _, err := exec.Command("tar", "-xvf", path.Join("/root", "busybox.tar"), "-C", lowerPath).CombinedOutput(); err != nil {
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
		err := os.Mkdir(upperPath, 0777)
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
		err := os.Mkdir(workPath, 0777)
		if err != nil {
			log.Println("can't create the work,error:", err.Error())
			return
		}
	}
	if err != nil && !os.IsNotExist(err) {
		log.Println("can't judge the work file's state")
	}

}
func NewMerged(rootURL string) {
	mergedPath := path.Join(rootURL, "merged")
	_, err := os.Stat(mergedPath)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(mergedPath, 0777)
			if err != nil {
				log.Println("can't create the merged ,error:", err.Error())
				return
			}
		} else {
			log.Println("can't judge the merged file's state")
			return
		}
	}

	dir := "upperdir=" + path.Join(rootURL, "upper") + ",lowerdir=" + path.Join(rootURL, "lower") + ",workdir=" + path.Join(rootURL, "work")
	log.Println("the dir is :--->", dir)
	cmd := exec.Command("mount", "-o", dir, "-t", "overlay", "overlay", path.Join(rootURL, "merged"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println("overlayFS mount failed ,error :", err)
	}
}

func VolumeMount(rootURL string, volume []string) {
	//创建宿主机目录
	if err := os.Mkdir(volume[0], 0755); err != nil {
		log.Println("the state of the ", volume[0], "in host is :", err.Error())
	}

	mergedPath := path.Join(rootURL, "merged")
	target := path.Join(mergedPath, volume[1])
	//创建容器目录
	if err := os.Mkdir(target, 0755); err != nil {
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
