package network

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"path"
	"strings"
)

//IPAM 用于ip地址分配信息
type IPAM struct {
	//分配文件存放位置
	SubnetAllocatorPath string
	//key是网段，value是位图数组
	Subnets *map[string]string
}
//默认存储位置
const ipamDefaultAllocatorPath = "/var/run/gocker/network/ipam/subnet.json"
var ipAllocator = &IPAM{
	SubnetAllocatorPath: ipamDefaultAllocatorPath,
}

func (ipam *IPAM)Allocate(subnet *net.IPNet)(ip net.IP,err error){
	ipam.Subnets = &map[string]string{}
	err = ipam.load()
	if err != nil{
		log.Println("load allocation information failed,error:",err)
		return 
	}

	_ , subnet,_ = net.ParseCIDR(subnet.String())
	//返回网段的子网掩码的总长度和网段前面固定位的长度
	one,size := subnet.Mask.Size()


	//如果之前没有分配过这个网段，则先初始化网段的分配配置
	if _,exist := (*ipam.Subnets)[subnet.String()];!exist{
		//size - one 代表后面的网络位数 , 2^(size-one)表示可用ip数
		//全部用0填满
		(*ipam.Subnets)[subnet.String()]=strings.Repeat("0", 1<< uint8(size-one))
	}
	for c := range ((*ipam.Subnets)[subnet.String()]){
		if (*ipam.Subnets)[subnet.String()][c]=='0'{
			//分配这个ip；string无法修改，先转换为byte数组，修改后再切换
			ipalloc := []byte((*ipam.Subnets)[subnet.String()])
			ipalloc[c]= '1'
			(*ipam.Subnets)[subnet.String()] = string(ipalloc)
			ip = subnet.IP
			for t := uint(4);t >0;t -=1{
				[]byte(ip)[4-t] += uint8(c >> ((t-1)*8))
				
			}
			ip[3]+=1
			break
				
		}
	}
	

	//调用dump，保存到文件中
	ipam.dump()
	return
}
func (ipam *IPAM)Release(subnet *net.IPNet,ipaddr *net.IP)error{
	ipam.Subnets = &map[string]string{}
	 _,subnet,_ = net.ParseCIDR(subnet.String())
	err := ipam.load()
	if err != nil{
		log.Println("can't load the ipam information")
		return err
	}

	c :=0 
	releaseIp := ipaddr.To4()
	releaseIp[3]-=1
	for t := uint(4);t >0;t-=1{
		c += int(releaseIp[t-1] - subnet.IP[t-1]) <<((4-t) * 8)
	}
	ipalloc := []byte((*ipam.Subnets)[subnet.String()])
	log.Println(c," ",len(ipalloc))
	ipalloc[c]= '0'
	(*ipam.Subnets)[subnet.String()] = string(ipalloc)
	//调用dump，保存到文件中
	ipam.dump()
	return nil
}


func (ipam *IPAM)load()error{
	if _,err := os.Stat(ipam.SubnetAllocatorPath);err != nil{
		if os.IsNotExist(err){
			return nil
		}else {
			return err
		}

	}
	file,err := os.Open(ipam.SubnetAllocatorPath)
	if err != nil{
		return err
	}
	defer file.Close()
	subnetJson := make([]byte,2000)
	n,err := file.Read(subnetJson)
	if err != nil{
		return err
	}
	err = json.Unmarshal(subnetJson[:n],ipam.Subnets)
	if err != nil{
		log.Println("can't unmarshal the data,error:",err)
		return err
	}
	return nil
}
func (ipam *IPAM)dump()error{
	configDir,_:=path.Split(ipam.SubnetAllocatorPath)
	if _,err := os.Stat(configDir);err != nil{
		if os.IsNotExist(err){
			if err = os.MkdirAll(configDir,0744);err != nil{
				log.Printf("can't create the directory %s,error:%v\n",configDir,err)
				return err
			}
		}else {
			return err
		}
	}
	configFile ,err := os.OpenFile(ipam.SubnetAllocatorPath,os.O_TRUNC | os.O_WRONLY | os.O_CREATE,0644)
	if err != nil{
		return err
	}
	defer configFile.Close()
	data,err := json.Marshal(ipam.Subnets)
	if err != nil{
		return err
	}
	if _, err = configFile.Write(data);err != nil{
		return err
	}
	return nil

}