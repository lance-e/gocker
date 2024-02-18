package container
import(
	"errors"
	"os/exec"
	"log"
	"os"
	"syscall"
)

// 此时容器已经创建，这是容器内的第一个进程，先去mount到proc文件系统，方便之后的ps等操作
func InitProcess() error {
	data := readUserCommand()
	if  len(data) ==0 {
		return errors.New("Run container get command failed")
	}
	defaultMountFlag := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_RELATIME  
	// 这里的 MountFlag 的意思如下。
	// 1. MS_NOEXEC在本文件系统中不允许运行其他程序。
	// 2. MS_NOSUID在本系统中运行程序的时候， 不允许 set-user-ID或 set-group-ID。
	// 3. MS_NODEV这个参数是自从Linux2.4以来，所有mount的系统都会默认设定的参数。
	err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlag), "")
	if err != nil {
		log.Println("proc mount happened error: ", err.Error())
	}
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
