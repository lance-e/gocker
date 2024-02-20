package cgroup

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path"

	"strings"
)

// 1.cgroup hierarchy 中的节点，用于管理进程和 subsystem 的控制关系。
// 2.subsystem 作用于 hierarchy 上的 cgroup 节 点，并控制节点中进程的资源占用。
// 3.hierarchy 将 cgroup 通过树状结构串起来，并通过虚拟文件系统的方式暴露给用户。

// ResouceConfig 定义了资源配置，包括内存限制，cpu时间片权重，cpu数
type ResouceConfig struct {
	MemoryLimit string
	CpuSet      string
	CpuShare    string
}

type Subsystem interface {
	//subsystem的名称，例如cpu，memory
	Name() string
	//设置cgroup在subsystem的资源限制
	Set(cgroupPath string, resouce *ResouceConfig) error
	//添加某个进程
	Apply(cgroupPath string, pid int) error
	//移除某个cgroup
	Remove(cgroupPath string) error
}

var (
	SubsystemInstance = []Subsystem{
		&MemorySubsystem{},
		&CpusetSubsystem{},
		&CpuSubsystem{},
	}
)

// FindCgroupMountPoint 获取cgroup挂载点
func FindCgroupMountPoint(subsystem string) string {
	file, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		log.Println("can't open /proc/self/mountinfo")
		return "error :can't open file"
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		fileds := strings.Split(txt, " ")
		for _, opt := range strings.Split(fileds[len(fileds)-1], ",") {
			if opt == subsystem {
				return fileds[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}
	return ""
}

// GetCgroupAbsolutePath 获取cgroup绝对路径
// 找到对应subsystem挂载的 hierarchy 相对路径对应的 cgroup 在虚拟文件系统中的路径,
// 然后通过这个目录的读写去操作 cgroup
// Cgroups 的 hierarchy 的虚拟文件系统是通过 cgroup类型文件系统的 mount挂载上去的,
// option 中加上 subsystem，代表挂载 的 subsystem 类型 ,
// 这样就可 以在 mountinfo 中找到对应的 subsystem 的挂载目录了 ，比如 memory。
func GetCgroupAbsolutePath(subsys string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountPoint(subsys)
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755)
			if err != nil {
				return "", errors.New("cgroup path error :" + err.Error())
			}
			return path.Join(cgroupRoot, cgroupPath), nil
		}
		return path.Join(cgroupRoot, cgroupPath), nil

	} else {
		return "", errors.New("cgroup path error :" + err.Error())
	}

}
