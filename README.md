# 从零开始手写docker
[项目地址](https://github.com/lance-e/gocker)
欢迎folk，star，follow🥰
## docker 三大核心技术：


##### Namespace:



- UTS Namespace 主要用来隔离 nodename 和 domainname 两个系统标识

- IPC Namespace 用来隔离 System V IPC 和 POSIX message queues(进程间通信)

- PID Namespace是用来隔离进程 ID的

- MountNamespace用来隔离各个进程看到的挂载点视图

- User Namespace 主要是隔离用户的用户组 ID

- Network Namespace 是用来隔离网络设备、 IP地址端口 等网络械的 Namespace。



##### Cgroups:

四个重要的概念：tasks，cgroup，hierarchy，subsystem

- tasks:就是一个进程
- cgroup:控制族群，也就是一个按某种规定划分的进程组，cgroups使用的资源都是以控制族群为单位划分的
- hierarchy:层级，是由多个cgroup组成的树状关系，
- subsystem:资源子系统，相当于资源控制器，如cpu，memory子系统，必须附加到某个层级hierarchy才能起作用，



##### union file system:

把其他文件系统联合到一个联合挂载点的文件系统服务，具有写时复制和联合挂载的特性

- 现在docker多采用的是overlay2或者aufs


## Docker详细构建教程：

### 一.构建容器：
通过实现一个简易的run命令，来构建容器
run命令的实现流程：

- 通过newparentProcess函数，构建一个父进程，此时已经进行了namespace
- 使用cgroup来对资源的限制，此时容器就已经创建完毕
- 创建init子进程(容器内第一个进程)，mount到/proc文件系统(方便ps命令)，同时使用syscall.Exec来覆盖之前的进程信息，堆栈信息(保证第一个进程是我们规定的进程)。

在容器中简单实现namespace和cgroup：
namespace的实现之间进行系统调用：
~~~go
cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUTS,
		//syscall.CLONE_NEWUSER,
	}
~~~
我们这里暂时不实现user namespace，因为较复杂会牵涉到权限等问题
cgroup则是通过读取/proc/self/mountinfo文件，获取当前进程的mount 情况，再根据我们需要限制的subsystem来获取到cgroup的挂载点，例如：/sys/fs/cgroup/memory，此时的subsystem为memory。
(本项目中仅支持了memory,cpuset ,cpushare进行了资源限制，本质上都是对文件进行读写操作)
~~~go
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
				log.Println("the mount point : ",fileds[4])
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
// option 中加上 subsystem，代表挂载的 subsystem 类型 ,
// 这样就可以在 mountinfo 中找到对应的 subsystem 的挂载目录了，比如 memory。
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
~~~
找到了cgroup挂载的绝对路径，就可以通过操作文件来进行资源限制,这里以memory为例：
通过set方法，将设置的内存资源限制写入memory.limit_in_bytes文件中；

通过apply方法，把pid写入到tasks，目标进程加入到该cgroup中(！！！这里必须要注意，如果你在set中写入的数据格式不争取，是无法将pid写入tasks的)；

通过remove方法，则是取消该cgroup
~~~go
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
~~~

### 二.构建镜像：
补充知识：
rootfs(root filesystem):
- rootfs是分层文件树的顶端，包含对系统运行至关重要的文件和目录，包括设备目录和用于启动系统的程序。系统启动时，初始化进程会将rootfs挂载到/目录，之后再挂载其他文件系统到其子目录。

mount namespace 工作原理：
- 每个进程可以创建属于自己的 mount table，但前提是必须先复制父进程的 mount table，之后再调用 mount 发生的更改都只会影响当前进程的 mount table
  
pivot_root系统调用介绍：
- pivot_root 是由 Linux 提供的一种系统调用，它能够将一个 mount namespace 中的所有进程的根目录和当前工作目录切换到一个新的目录。pivot_root 的主要用途是在系统启动时，先挂载一个临时的 rootfs 完成特定功能，然后再切换到真正的 rootfs。
- 可以将当前root文件系统移动到put_old文件夹中，然后将new_root成为新的root文件系统(注：new_root和put_old不能同时存在当前root的同一个文件系统中)

pivot_root与chroot区别：
- chroot只改变某个进程的根目录，系统的其他部分依旧运行于旧的root目录。 pivot_root把整个系统切换到一个新的root目录中，然后去掉对之前rootfs的依赖，以便于可以umount之前的文件系统。

##### 1.使用busybox来构建极简镜像
项目中使用privot_root系统调用的函数：
~~~go
// pivotRoot 进行pivot_root系统调用
func pivotRoot(root string )error{
	
	//为了使当前root文件系统的老root文件系统和新root文件系统不在同一个文件系统下，这里把root重新mount一次
	//bind mount 就是把相同的内容换一个挂载点的挂载方式
	err := syscall.Mount(root,root,"bind",syscall.MS_BIND |syscall.MS_REC,"")
	if err != nil{
		return errors.New("Mount rootfs to itself failed,error:"+err.Error())
	}
	//存储旧root文件系统
	pivotDir := filepath.Join(root,".pivot_root")
	err = os.Mkdir(pivotDir,0777)
	if err != nil{
		return err
	}
	//root 为新root文件系统，pivotDir代表put_old文件夹，将旧root文件系统放在pivotDir文件夹中
	err = syscall.PivotRoot(root,pivotDir)
	if err != nil {
		return errors.New("pivot_root,error :"+err.Error())
	}
	err = syscall.Chdir("/")
	if err != nil{
		return errors.New("chdir,error:"+err.Error())
	}
	//此时的pivotDir就是刚刚存放旧的root文件系统的文件夹
	pivotDir = filepath.Join("/",".pivot_root")
	err = syscall.Unmount(pivotDir,syscall.MNT_DETACH)
	if err != nil{
		return errors.New("umount pivot_roo5t directory failed ,error :"+err.Error())
	}
	return os.Remove(pivotDir)
}
~~~

再将原先在InitProcess函数中进行mount操作移到setUpMount函数中，同时进行pivot_root：
~~~go
func setUpMount(){
	pwd ,err := os.Getwd()
	if err != nil{
		log.Println("get current location error:"+err.Error())
		return 
	}
	log.Println("the current location is "+pwd)
	
	err = syscall.Mount("","/","",syscall.MS_REC | syscall.MS_PRIVATE,"")
	if err != nil{
		log.Println("the first mount failed,error:",err.Error())
	}
	err = pivotRoot(pwd)
	if err != nil{
		log.Println("pivot_root system call failed")
	}
	//mount proc
	defaultMountFlag := syscall.MS_NODEV | syscall.MS_NOSUID | syscall.MS_NOEXEC 
	// 这里的 MountFlag 的意思如下。
	// 1. MS_NOEXEC在本文件系统中不允许运行其他程序。
	// 2. MS_NOSUID在本系统中运行程序的时候， 不允许 set-user-ID或 set-group-ID。
	// 3. MS_NODEV这个参数是自从Linux2.4以来，所有mount的系统都会默认设定的参数。
	syscall.Mount("proc","/proc","proc",uintptr(defaultMountFlag),"")
	//mount tmpfs
	syscall.Mount("tmpfs","/dev","tmpfs",syscall.MS_NOSUID | syscall.MS_STRICTATIME,"mode=755")
}
~~~
这里挂载到/，可以使后面挂载的/proc在退出容器时自动umount /proc,因为这样可以声明这个新的mount namespace独立
(注!!!:这个/挂载，必须要在所有挂载之前)


现在InitProcess函数就是这个样子：
~~~go
func InitProcess() error {

	data := readUserCommand()
	if  len(data) ==0 {
		return errors.New("Run container get command failed")
	}
	setUpMount() //将mount封装
	
	log.Println("mount success")
	//通过exec.LookPath找到命令在环境变量中路径
	cmdpath, err := exec.LookPath(data[0])
	if err != nil {
		log.Println(data[0]," look path in PATH environment variable failed")
		return err
	}

	err = syscall.Exec(cmdpath, data[0:], os.Environ())
	//!!!最重要的操作
	//syscall.Exec这个方法,
	//其实最终调用了Kernel的intexecve(const char *filename,char *const argv[], char *const envp[]);
	//这个系统函数。它的作用是执行当前 filename对应的程序。
	//它会覆盖当前进程的镜像、数据和堆栈等信息，包括 PID， 这些都会被将要运行的进程覆盖掉。
	//保证了我们进入容器之后，我们看到的第一个进程是我们指定的进程，因为之前的信息都被覆盖掉了
	if err != nil {
		log.Fatal("error :", err.Error())
	}
	log.Println("exec  success")
	return nil
}
~~~

完成后的结果就是这样：
~~~bash
/ # ps
PID   USER     TIME  COMMAND
    1 root      0:00 sh
    7 root      0:00 ps
/ # mount
/dev/sdc on / type ext4 (rw,relatime,discard,errors=remount-ro,data=ordered)
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
tmpfs on /dev type tmpfs (rw,nosuid,mode=755)
/ # ls
bin   dev   etc   home  proc  root  sys   tmp   usr   var
~~~
此时子进程就看不到父进程的mount信息，且rootfs切换成了我们设置的busybox

##### 2.使用union filesystem来包装镜像
这里选择使用overlayFS，作为unionfs

相较于aufs，overlayFS的优势：
- 速度更快，aufs层数更多，性能损耗更大
- 简单，overlay2只有两层，容器层upper和镜像层lower
- overlay2加入了linux kernel

相较于overlay驱动，overlay2驱动的优势：
- overlay驱动只在一个lower overlayFS层之上，所以为了实现多层镜像需要大量的硬链接
- overlay2驱动原生支持多个lower overlayFS

组成：
- lower：镜像层，存储镜像文件，且只能读不能写
- upper：容器层，可读可写，写时复制时，就是将需要写入的文件复制到upper中，进行修改，后续就直接在复制的文件中修改
- merged：挂载的文件，展示参与联合挂载的目录的文件
- work：主要是用来保证操作的原子性

通过命令来展示overlayFS的挂载,挂载前只有lower和upper中有文件，merged为空，挂载后merged展示了参与联合挂载的目录文件：
~~~bash
root@localhost:~/test# mount -t overlay overlay ./merged -o upperdir=./upper,lowerdir=./lower,workdir=./work
root@localhost:~/test# ls
lower  merged  upper  work
root@localhost:~/test# tree
.
├── lower
│   └── lll
├── merged
│   ├── lll
│   └── uuu
├── upper
│   └── uuu
└── work
~~~
(这里我们还是使用busybox来作为镜像)
创建upper，lower，work，merged目录，同时挂载overlayFS:
~~~go
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
~~~

取消挂载，并删除其他目录，保留镜像目录：
~~~go
func DeleteWorkSpace(){
	DeleteMerged("/root/overlay")
	DeleteUpper("/root/overlay")
	DeleteWork("/root/overlay")
}

func DeleteMerged(rootURL string){
	mergedPath := path.Join(rootURL,"merged")
	cmd := exec.Command("umount",mergedPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run();err != nil{
		log.Println("umount overlay failed")
	}
	log.Println("umount overlay successful...")
	if err := os.RemoveAll(mergedPath);err != nil{
		log.Println("remove merged directory failed")
		return
	}
}
func DeleteUpper(rootURL string){
	upperPath := path.Join(rootURL,"upper")
	if err := os.RemoveAll(upperPath);err != nil{
		log.Println("remove upper diroctory failed")
		return
	}

}

func DeleteWork(rootURL string){
	workPath := path.Join(rootURL,"work")
	if err := os.RemoveAll(workPath);err != nil{
		log.Println("remove work diroctory failed")
		return
	}
}
~~~
通过cmd.Dir把overlay的挂载点作为容器的根目录

##### 3.使用volume数据卷挂载数据持久化

docker中的volume挂载时，已经启动了容器进程，此时已经有了mount namespace，但是现在未进行pivot_root系统调用或者chroot，在容器中还是可以观察到宿主机的全部文件系统，我们只需要在进行pivot_root系统调用前，将宿主机上的目录挂载到指定容器目录



对于volume的参数校验就不展示了，这里仅展示挂载数据卷：
~~~go
func VolumeMount(rootURL string, volume []string) {
	//创建宿主机目录
	if err := os.Mkdir(volume[0], 0755); err != nil {
		log.Println("the state of the ", volume[0], "in host is :", err.Error())
	}

	mergedPath := path.Join(rootURL, "merged")
	target := path.Join(mergedPath, volume[1])
	//创建容器目录
	if err := os.Mkdir(target, 0755); err != nil {
		log.Println("the state of the ", volume[1], "in container is :", err.Error())
	}
	cmd := exec.Command("mount", "-o", "bind", volume[0], target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println("mount volume failed,error:", err.Error())
		return
	}
	log.Printf("mount %s to %s successful!\n", volume[0], target)
}
~~~

注意数据卷和文件系的顺序：
- 挂载：先进行文件系统的挂载，最后进行数据卷的挂载。
- 取消挂载：先进行数据卷的取消挂载，再进行文件系统的取消挂载



## 参考文档：
https://blog.csdn.net/qq_53267860/article/details/131729601

https://blog.csdn.net/qq_31960623/article/details/120242671

https://blog.csdn.net/qq_31960623/article/details/120260769

https://www.cnblogs.com/crazymakercircle/p/15400946.html#autoid-h3-2-2-0

https://tech.meituan.com/2015/03/31/cgroups.html

https://www.cnblogs.com/charlieroro/p/10281469.html

https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/7/html/resource_management_guide/index

https://waynerv.com/posts/container-fundamentals-filesystem-isolation-and-sharing/

https://cloud.tencent.com/developer/article/1681523

https://www.cnblogs.com/FengZeng666/p/14173906.html

https://blog.csdn.net/luckyapple1028/article/details/78075358

https://zhuanlan.zhihu.com/p/374924046

https://www.cnblogs.com/istitches/p/18011539