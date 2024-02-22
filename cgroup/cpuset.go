package cgroup

import (
	"errors"
	"log"
	"os"
	"path"
	"strconv"
)

type CpusetSubsystem struct {
}

func (c *CpusetSubsystem) Name() string {
	return "cpuset"
}

func (c *CpusetSubsystem) Set(cgroupPath string, resource *ResouceConfig) error {
	absolutePath, err := GetCgroupAbsolutePath(c.Name(), cgroupPath, true)
	if err != nil {
		return errors.New("get absolute path falied,error:" + err.Error())
	}
	log.Println("the absolute path is :",absolutePath)
	//一些子系统有强制参数，在您将任务移至使用这些子系统的 cgroup 中之前，这些参数必须被设定。
	// 例如:一个使用 cpuset子系统的 cgroup, 在您将任务移至此 cgroup 前，
	//cpuset.cpus 和 cpuset.mems 参数必须被设定。
	if resource.CpuSet != "" {
		if err = os.WriteFile(path.Join(absolutePath, "cpuset.cpus"), []byte(resource.CpuSet), 0644); err != nil {
			return errors.New("cpuset set falied,error:" + err.Error())
		}
		if err = os.WriteFile(path.Join(absolutePath, "cpuset.mems"), []byte(resource.CpuSet), 0644); err != nil {
			return errors.New("cpuset set falied,error:" + err.Error())
		}
	}
	return nil

}
func (c *CpusetSubsystem) Apply(cgroupPath string, pid int) error {
	absolutePath, err := GetCgroupAbsolutePath(c.Name(), cgroupPath, false)
	if err != nil {
		return errors.New("get absolute path falied,error:" + err.Error())
	}
	log.Println("the absolute path is :",absolutePath)
	err = os.WriteFile(path.Join(absolutePath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return err
	}
	return nil
}
func (c *CpusetSubsystem) Remove(cgroupPath string) error {
	absolutePath, err := GetCgroupAbsolutePath(c.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	log.Println("the absolute path is :",absolutePath)
	err = os.RemoveAll(absolutePath)
	if err != nil {
		return err
	}
	return nil
}
