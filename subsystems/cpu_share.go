package subsystems

import (
	"errors"
	"os"
	"path"
	"strconv"
)

type CpuSubsystem struct {

}
func (c *CpuSubsystem)Name()string{
	return "cpu"
}
func (c *CpuSubsystem)Set(cgroupPath string ,resource *ResouceConfig)error{
	absolutePath,err :=GetCgroupAbsolutePath(c.Name(),cgroupPath,true)
	if err != nil{
		return errors.New("cpu set failed,error:"+err.Error())
	}
	if resource.CpuShare!=""{
		err = os.WriteFile(path.Join(absolutePath,"cpu.shares"),[]byte(resource.CpuShare),0644)
		if err!= nil{
			return errors.New("cpu set failed,error:"+err.Error())
		}
	}
	
	return nil
}
func (c *CpuSubsystem)Apply(cgroupPath string ,pid int)error{
	absolutePath,err :=GetCgroupAbsolutePath(c.Name(),cgroupPath,false)
	if err != nil{
		return errors.New("cpu apply failed,error:"+err.Error())
	}
	err = os.WriteFile(path.Join(absolutePath,"tasks"),[]byte(strconv.Itoa(pid)),0644)
	if err != nil{
		return errors.New("cpu apply failed,error:"+err.Error())
	}
	return nil

}
func (c *CpuSubsystem)Remove(cgroupPath string)error{
	absolutePath,err :=GetCgroupAbsolutePath(c.Name(),cgroupPath,false)
	if err != nil{
		return errors.New("cpu remove failed,error:"+err.Error())
	}
	err = os.Remove(absolutePath)
	if err != nil{
		return errors.New("cpu remove failed,error:"+err.Error())
	}
	return nil
}