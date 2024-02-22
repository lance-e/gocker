package cgroup
import "log"
type CgroupManager struct {
	Path     string // cgroup 在hierarchy中的路径，就是创建的cgroup相当于root cgroup的相对路径
	Resource *ResouceConfig
}

func NewCgroupManager(p string) *CgroupManager {
	return &CgroupManager{
		Path: p,
	}
}

func (c *CgroupManager) Apply(pid int) error {
	log.Println("cgroup apply...")
	for _, instance := range SubsystemInstance {
		instance.Apply(c.Path, pid)

	}
	return nil
}
func (c *CgroupManager) Set(res *ResouceConfig) error {
	log.Println("cgroup set...")
	for _, instance := range SubsystemInstance {
		instance.Set(c.Path, res)
	}
	return nil
}
func (c *CgroupManager) Destory() error {
	log.Println("cgroup destory...")
	for _, instance := range SubsystemInstance {
		instance.Remove(c.Path)

	}
	return nil
}
