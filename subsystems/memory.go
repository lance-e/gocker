package subsystems

import (
	"errors"
	"os"
	"path"
	"strconv"

)

// memory subsystem 的实现
type MemorySubsystem struct {
}

func (m *MemorySubsystem) Name() string {
	return "memory"
}

// 设置对应的cgroup memory限制
func (m *MemorySubsystem) Set(cgroupPath string, resource *ResouceConfig) error {
	//获取subsystem路径(cgroup绝对路径)
	absolutePath, err := GetCgroupAbsolutePath(m.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	//设置内存限制
	if resource.MemoryLimit != ""{
		err = os.WriteFile(path.Join(absolutePath, "memory.limit_in_bytes"), []byte(resource.MemoryLimit), 0644)
		if err != nil {
			return errors.New("set memory failed,error:" + err.Error())
		}
	}
	
	return nil

}
func (m *MemorySubsystem) Apply(cgroupPath string, pid int) error {
	absolutePath,err := GetCgroupAbsolutePath(m.Name(),cgroupPath,false)
	if err != nil{
		return errors.New("apply cgroup faield ,error:"+err.Error())
	}
	err = os.WriteFile(path.Join(absolutePath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return errors.New("apply cgroup faield ,error:" + err.Error())
	}
	return nil

}
func (m *MemorySubsystem) Remove(cgroupPath string) error {
	absolutePath ,err := GetCgroupAbsolutePath(m.Name(),cgroupPath,false)
	if err != nil{
		return err
	}
	err = os.Remove(absolutePath)
	if err != nil{
		return err
	}
	return nil
}
