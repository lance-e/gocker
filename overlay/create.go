package overlay

import (
	"log"
	"os"
	"os/exec"
	"path"
)

func NewWorkSpace() {
	_, err := os.Stat("/root/overlay")
	if os.IsNotExist(err) {
		if err := os.Mkdir("/root/overlay", 0777); err != nil {
			log.Println("can't create a new work space,error:", err.Error())
			return
		}

	}
	NewUpper("/root/overlay")
	NewLower("/root/overlay")
	NewWork("/root/overlay")
	NewMerged("/root/overlay")
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
