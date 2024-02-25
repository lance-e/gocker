package container

import (
	"errors"
	"log"
	"os"
	"os/exec"
	
	"path/filepath"
	"syscall"
)

// 此时容器已经创建，这是容器内的第一个进程，先去mount到proc文件系统，方便之后的ps等操作
func InitProcess() error {

	data := readUserCommand()
	if  len(data) ==0 {
		return errors.New("Run container get command failed")
	}
	setUpMount() //将mount封装
	
	//log.Println("mount success")
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
	//log.Println("exec  success")
	return nil
}


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

func setUpMount(){
	pwd ,err := os.Getwd()
	if err != nil{
		log.Println("get current location error:"+err.Error())
		return 
	}
	log.Println("the current location is "+pwd)
	//注：这个/挂载，非常非常重要！！！
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