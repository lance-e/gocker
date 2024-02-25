package network

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/tabwriter"

	"log"
	"net"
	"os"
	"path"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

var (
	defaultNewWorkPath = "/var/run/gocker/network/network"
	drivers            = map[string]NetWorkDriver{}
	networks           = map[string]*NetWork{}
)

type NetWork struct {
	Name    string     //网络名
	IpRange *net.IPNet //地址段
	Driver  string     //驱动名

}

type EndPoint struct {
	Id          string           `json:"id"`
	Device      netlink.Veth     `json:"dev"`
	IpAddress   net.IP           `json:"ip"`
	MacAddress  net.HardwareAddr `json:"mac"`
	PortMapping []string         `json:"portmapping"`
	NetWork     *NetWork
}

type NetWorkDriver interface {
	//驱动名
	Name() string
	//创建驱动
	Create(subnet string, name string) (*NetWork, error)
	//删除驱动
	Delete(network NetWork) error
	//连接
	Connect(network *NetWork, endpoint *EndPoint) error
	//断开连接
	Disconnect(network *NetWork, endpoint *EndPoint) error
}

func CreateNetWork(driver string, subnet string, name string) error {
	_, cidr, _ := net.ParseCIDR(subnet)
	//调用IPAM分配网关ip，获取到网段中的第一个ip作为网关的ip
	gatewayIp, err := ipAllocator.Allocate(cidr)
	if err != nil {
		return err
	}
	cidr.IP = gatewayIp
	//drivers字典是各个网络驱动的实例字典，然后调用create创建网络
	nw, err := drivers[driver].Create(cidr.String(), name)
	if err != nil {
		return err
	}
	return nw.dump(defaultNewWorkPath)

}
func Connect(networkName string, id string, pid string, portmapping []string) error {
	//通过networks字典获取networkName对应的network
	fmt.Println("in connect :", networks)
	network, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("no such network:%s", networkName)
	}
	//通过调用IPAM从网络的网段中获取可用ip作为容器ip地址
	ip, err := ipAllocator.Allocate(network.IpRange)
	if err != nil {
		return err
	}
	ep := &EndPoint{
		Id:          fmt.Sprintf("%s-%s", id, networkName),
		IpAddress:   ip,
		NetWork:     network,
		PortMapping: portmapping,
	}
	//调用Connect方法去连接和配置网络端点
	if err = drivers[network.Driver].Connect(network, ep); err != nil {
		return err
	}
	//配置容器网络设备的ip地址和路由
	if err = configEndpointIAddressAndRoute(ep, pid); err != nil {
		return err
	}
	//配置容器到宿主机的端口映射
	return configPortMapping(ep)
}
func Init() error {
	//加载网络驱动
	var bridgeDriver = BridgeNetworkDriver{}
	drivers[bridgeDriver.Name()] = &bridgeDriver
	_, err := os.Stat(defaultNewWorkPath)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(defaultNewWorkPath, 0644)
		} else {
			return err
		}
	}
	filepath.Walk(defaultNewWorkPath, func(nwpath string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		_, name := path.Split(nwpath)
		nw := &NetWork{
			Name: name,
		}
		if err := nw.load(nwpath); err != nil {
			log.Println("can't load the network ,error:", err)
		}
		//加入到network字典中
		networks[name] = nw
		return nil
	})
	//fmt.Println("the network is ", networks)

	return nil
}
func ListNetWork() {
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "NAME\tIpRange\tDriver\n")
	for _, nw := range networks {
		fmt.Fprintf(w, "%s\t%s\t%s\n", nw.Name, nw.IpRange.String(), nw.Driver)
	}
	if err := w.Flush(); err != nil {
		log.Println("flush error:", err)
		return
	}
}
func DeleteNetWork(networkName string) error {
	nw, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("no such network: %s", networkName)
	}
	//调用IPAM的实例ipAllocator释放网络网关的IP
	if err := ipAllocator.Release(nw.IpRange, &nw.IpRange.IP); err != nil {
		return fmt.Errorf("remove network gateway failed,error:%v", err)
	}
	if err := drivers[(*nw).Driver].Delete(*nw); err != nil {
		return fmt.Errorf("remove network driver failed ,error:%v", err)
	}
	return nw.remove(defaultNewWorkPath)
}
func (nw *NetWork) dump(dumpPath string) error {
	if _, err := os.Stat(dumpPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dumpPath, 0644)
		} else {
			return err
		}
	}
	//文件名是网络的名字
	nwPath := path.Join(dumpPath, nw.Name)
	file, err := os.OpenFile(nwPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println("open the network file failed,error:", err)
		return err
	}
	defer file.Close()
	nwJson, err := json.Marshal(nw)
	if err != nil {
		log.Println("can't marshal the network data,error:", err)
		return err
	}
	_, err = file.Write(nwJson)
	if err != nil {
		log.Println("can't write the json data into file,error:", err)
		return err
	}
	return nil
}
func (nw *NetWork) load(dumpPath string) error {
	file, err := os.Open(dumpPath)
	if err != nil {
		return err
	}
	defer file.Close()
	nwJson := make([]byte, 2000)
	n, err := file.Read(nwJson)
	if err != nil {
		return err
	}
	err = json.Unmarshal(nwJson[:n], nw)
	if err != nil {
		log.Println("load failed ,error :", err)
		return err
	}
	return nil
}
func (nw *NetWork) remove(dumpPath string) error {
	if _, err := os.Stat(path.Join(dumpPath, nw.Name)); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	} else {
		return os.Remove(path.Join(dumpPath, nw.Name))
	}
}

func configEndpointIAddressAndRoute(endpoint *EndPoint, pid string) error {
	peerLink, err := netlink.LinkByName(endpoint.Device.PeerName)
	if err != nil {
		log.Println("config endpoint failed,error:", err)
		return err
	}

	//将容器的网络端点加入到容器的网络空间中，并使这个函数下面的操作都在这个网络空间中，直到执行完函数

	defer enterContainerNetns(&peerLink, pid)()
	interfaceIP := *endpoint.NetWork.IpRange
	interfaceIP.IP = endpoint.IpAddress

	if err = setInterfaceIP(endpoint.Device.PeerName, interfaceIP.String()); err != nil {
		log.Println("set the interface ip failed,error:", err)
		return err
	}
	if err = setInterfaceUp(endpoint.Device.PeerName); err != nil {
		return err
	}
	//net namespace 中默认本地地址的127.0.0.1的lo网卡使关闭的，启动它以保证容器访问自己的请求
	if err = setInterfaceUp("lo"); err != nil {
		return err
	}
	//设置容器内的外部请求都通过容器内的veth端点访问0.0.0.0/0网段，表示所有的ip地址段
	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")

	//构建要添加的路由数据
	defaultRoute := &netlink.Route{
		LinkIndex: peerLink.Attrs().Index,
		Gw:        endpoint.NetWork.IpRange.IP,
		Dst:       cidr,
	}

	//添加路由
	if err = netlink.RouteAdd(defaultRoute); err != nil {
		return err
	}
	return nil
}

// 将容器的网络端点加入到容器的网络空间，锁定当前程序执行的线程，使当前进程进入到容器的网络空间
// 返回值是个函数指针，执行这个函数才会退回到宿主机的网络空间
func enterContainerNetns(link *netlink.Link, pid string) func() {
	//找到容器的netns
	f, err := os.OpenFile(fmt.Sprintf("/proc/%s/ns/net", pid), os.O_RDONLY, 0)
	if err != nil {
		log.Println("get container net namespace failed,error:", err)
	}
	//获取文件描述符
	nsFD := f.Fd()

	//锁定线程，因为goroutine可能会被调度到其他线程上去，所以需要锁定操作系统线程
	runtime.LockOSThread()

	if err = netlink.LinkSetNsFd(*link, int(nsFD)); err != nil {
		log.Println("set lin netns failed,error:", err)
	}
	//获取到当前网络的net namespace
	origns, err := netns.Get()
	if err != nil {
		log.Println("get current netns failed,error:", err)

	}
	// 设置当前进程到新的网络namespace，并在函数执行完成之后再恢复到之前的namespace
	if err = netns.Set(netns.NsHandle(nsFD)); err != nil {
		log.Println("can't set the netns,error:",err)
	}

	return func() {
		//恢复到之前的netns
		netns.Set(origns)
		//关闭namespace文件
		origns.Close()
		//取消线程锁定
		runtime.UnlockOSThread()
		//关闭namespace文件
		f.Close()
	}

}
func configPortMapping(endpoint *EndPoint) error {
	for _, pm := range endpoint.PortMapping {
		portmapping := strings.Split(pm, ":")
		if len(portmapping) != 2 {
			log.Println("port mapping format error")
			continue
		}

		iptableCmd := fmt.Sprintf("-t nat -A PREROUTING -p tcp -m tcp --dport %s -j DNAT --to-destination %s:%s", portmapping[0], endpoint.IpAddress.String(), portmapping[1])
		cmd := exec.Command("iptables", strings.Split(iptableCmd, " ")...)
		output, err := cmd.Output()
		if err != nil {
			log.Println("the iptables output ,", output)
			continue
		}

	}
	return nil
}
