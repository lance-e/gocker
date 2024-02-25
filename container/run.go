package container

import (
	"fmt"
	"gocker/cgroup"
	"gocker/network"
	"gocker/overlay"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
)
const RootUrl = "/root/overlay"
const MergedPath = "/root/overlay/%s/merged"

func Run(tty bool,imageName string, command []string, resource *cgroup.ResouceConfig, volume string, containerName string,containerId string,environment []string,net string,port []string) {
	//原本启动init进程是通过在参数中添加一个init，来进行init命令，改成通过匿名管道进行父子进程间通信
	cmd, write := newparentProcess(tty, volume,containerName,imageName,environment)
	if cmd ==nil{
		log.Println("create parent proccess failed")
		return
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err.Error())
	}

	name, err := RecordContainerInformation(containerName,containerId, cmd.Process.Pid, command,volume)
	if err != nil {
		log.Println("can't record container information ")
		return
	}
	//新建一个cgroup manager来进行资源管理
	manager := cgroup.NewCgroupManager("gocker-cgroup")
	defer manager.Destory()
	manager.Set(resource)
	manager.Apply(cmd.Process.Pid)
	
	if net != ""{
		network.Init()
		if err := network.Connect(net,containerId,strconv.Itoa(cmd.Process.Pid),port);err != nil{
			log.Println("can't connect the network,error:",err)
		}
	}

	sendInitCommand(command, write)
	if tty {
		cmd.Wait() //父进程等待子进程
		DeleteContainerInfo(name)
		overlay.DeleteWorkSpace(path.Join(RootUrl,containerName), volume)

	}

}
func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}
func readUserCommand() []string {
	//通过提供的文件描述符，返回一个新文件
	//uintptr(3)指的是index为3的文件描述符，也就是管道文件读取端，前三个是stdin，stdout，stderr
	readPipe := os.NewFile(uintptr(3), "pipe")
	defer readPipe.Close()
	data, err := io.ReadAll(readPipe)
	if err != nil {
		log.Println("init read pipe failed,error :" + err.Error())
		return nil
	}
	log.Println("this is the data of pipe to translate :", string(data))
	return strings.Split(string(data), " ")
}
func sendInitCommand(cmdArray []string, write *os.File) {
	command := strings.Join(cmdArray, " ")
	//log.Println("send init command ,all is :", command)
	write.WriteString(command)
	write.Close()
}

// 1. 这里的/ proc/self/exe 调用中，/ proc/self/指的是当前运行进程自己的环境， exe 其实就是自己 调用了自己，使用这种方式对创建出来的进程进行初始化
// 2. 后面的 args 是参数，其中 init 是传递给本进程的第 一个参数，在本例中，其实就是会去调用 initCornmand 去初始化进程的 一些环境和资源
// 3. 下面的 clone 参数就是去 fork 出来一个新进程，并且使用了 namespace 隔离新创建的进程和外部环境 。
// 4. 如果用户指定了- ti 参数，就需要把当前进程的输入输出导入到标准输入输出上
func newparentProcess(tty bool, volume string,containerName string,imageName string,environment []string) (*exec.Cmd, *os.File) {
	read, write, err := NewPipe()
	if err != nil {
		log.Println("can't create a new pipe")
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUTS,
		//syscall.CLONE_NEWUSER,
	}
	//系统调用，简单实现docker的六种namespace隔离
	//log.Println("the tty is ", tty)
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		dir := fmt.Sprintf(DefaultInfoLocation,containerName)
		if err := os.MkdirAll(dir,0622);err!= nil{
			log.Println("can't make all directory")
			return nil ,nil
		}
		 file,err := os.Create(dir+ContainerLogFile)
		 if err != nil{
			log.Println("can't create the log file")
			return nil, nil
		}
		cmd.Stdout=file
	}
	//!!!在这里传入管道读取端的句柄
	cmd.ExtraFiles = []*os.File{read} //该属性的意思是外带着其他的文件句柄来创建子进程,因为一个进程默认带着三个文件描述符，stdin，stdout，stderr
	cmd.Env = append(os.Environ(),environment... )
	root := path.Join(RootUrl,containerName)
	overlay.NewWorkSpace(root, volume,imageName)
	cmd.Dir = fmt.Sprintf(MergedPath,containerName) //使用cmd.Dir设置初始化后的工作目录
	return cmd, write
}
