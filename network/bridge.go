package network

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"

	"github.com/vishvananda/netlink"
)

type BridgeNetworkDriver struct {

}
func (b *BridgeNetworkDriver)Name()string{
	return "bridge"
}
func (b *BridgeNetworkDriver)Create(subnet string ,name string)(*NetWork,error){
	ip,ipRange , _ := net.ParseCIDR(subnet)
	ipRange.IP=ip
	n := &NetWork{
		Name: name,
		IpRange: ipRange,
		Driver: b.Name(),
	}
	err := b.initBridge(n)
	if err != nil{
		log.Println("can't init the bridge,error:",err)
	}
	return n,err
}

func (b *BridgeNetworkDriver)Delete(network NetWork)error{
	name := network.Name
	br,err := netlink.LinkByName(name)
	if err != nil{
		return err
	}
	return netlink.LinkDel(br)
}

func (b *BridgeNetworkDriver)Connect(network *NetWork,endpoint *EndPoint)error{
	name := network.Name
	br,err := netlink.LinkByName(name)
	if err != nil{
		log.Println("can't get the target link")
		return err
	}
	//创建接口的配置
	link := netlink.NewLinkAttrs()
	//linux接口名限制，取endpoint ID前五位
	link.Name = endpoint.Id[:5]
	//通过设置veth接口的master属性，设置一端挂载到linux bridge上
	link.MasterIndex = br.Attrs().Index
	//创建veth对象，通过peername配置另一端接口名
	endpoint.Device = netlink.Veth{
		LinkAttrs: link,
		PeerName: "cif-"+endpoint.Id[:5],
	}
	//创建
	if err := netlink.LinkAdd(&endpoint.Device);err != nil{
		log.Println("can't add this link,error:",err)
		return err
	}

	//设置veth启动
	if err = netlink.LinkSetUp(&endpoint.Device);err != nil{
		log.Println("can't set up this link,error:",err)
		return err
	}
	return nil

}
func (b *BridgeNetworkDriver)Disconnect(network *NetWork,endpoint *EndPoint)error{
	return nil
}
func (b *BridgeNetworkDriver)initBridge(n *NetWork)error{
	name := n.Name
	//创建设备
	if err := createBridgeInterface(name);err != nil{
		return fmt.Errorf("add bridge %s failed,error:%v",name,err)

	}
	//设置bridge设备地址和路由
	gatewayIp := *n.IpRange
	gatewayIp.IP = n.IpRange.IP
	if err := setInterfaceIP(name,gatewayIp.String());err != nil{
		return fmt.Errorf("assigning address :%s on bridge: %s with an error:%v",&gatewayIp,name,err)
	}
	//启动bridge设备
	if err := setInterfaceUp (name);err != nil{
		return fmt.Errorf("set bridge up :%s ,error :%v",name,err)
	}
	//设置iptables的SNAT规则
	if err := setupIPTables(name,n.IpRange);err != nil{
		return fmt.Errorf("set iptable for :%s failed,error: %v",name,err)
	}
	return  nil
}

func createBridgeInterface(name string)error{
	//检查是否存在
	_,err := net.InterfaceByName(name)
	if err ==nil || !strings.Contains(err.Error(),"no such network interface"){
		return err
	}
	//初始化一个link对象
	link := netlink.NewLinkAttrs()
	link.Name = name
	//创建bridge对象
	br := &netlink.Bridge{LinkAttrs: link}
	//创建虚拟网络设备
	if err := netlink.LinkAdd(br);err != nil{
		return fmt.Errorf("create bridge failed ,error:%v",err)
	}
	return nil

}




func setInterfaceIP(name string ,ip string)error{
	i ,err := netlink.LinkByName(name)
	if err != nil{
		log.Println(err)
		return err
	}
	ipNet ,err := netlink.ParseIPNet(ip)
	if err != nil{
		return err
	}

	//给网络接口配置地址
	addr := &netlink.Addr{IPNet: ipNet,Label: "",Flags: 0,Scope: 0}
	return netlink.AddrAdd(i,addr)
}

func setInterfaceUp(name string)error{
	i ,err := netlink.LinkByName(name)
	if err != nil{
		log.Println("can't find the target link ,error:",err)
		return err
	}
	if err := netlink.LinkSetUp(i);err != nil{
		log.Println("can't set up the interface,error:",err)
		return err
	}
	return nil
}

func setupIPTables (name string ,subnet *net.IPNet)error{
	//-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE
	iptableCmd := fmt.Sprintf("-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE",subnet.String(),name )
	cmd := exec.Command("iptables",strings.Split(iptableCmd," ")...)
	output,err := cmd.Output()
	if err != nil{
		log.Println("iptables output :",output)
	}
	return nil
}