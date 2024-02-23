package container

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	_"gocker/setns"
)

const ENV_EXEC_PID = "gocker_pid"
const ENV_EXEC_CMD = "gocker_cmd"

func ExecContainer(containerName string, command []string) {
	pid, err := getPidByContainerName(containerName)
	if err != nil {
		log.Println("can't get pid by container name")
		return
	}
	cmdstr := strings.Join(command, " ")
	log.Printf("pid :%s,cmd :%s\n", pid, cmdstr)
	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	os.Setenv(ENV_EXEC_PID, pid)
	os.Setenv(ENV_EXEC_CMD, cmdstr)
	cmd.Env = append( os.Environ(),(pid))
	cmd.Dir = "/root/overlay/merged"

	if err := cmd.Run(); err != nil {
		log.Println("exec container failed, error:", err)
	}
}
func getPidByContainerName(name string) (string, error) {
	dir := fmt.Sprintf(DefaultInfoLocation, name)
	allDir := dir + ConfigName
	data, err := os.ReadFile(allDir)
	if err != nil {
		log.Println("can't read the config informaiton from the file")
		return "", err
	}
	var info = &ContainerInfo{}
	err = json.Unmarshal(data, info)
	if err != nil {
		log.Println("unmalshal the information failed")
		return "", err
	}
	return info.Pid, nil
}
func getEnvByPid(pid string)[]string {
	dir := fmt.Sprintf("/proc/%s/environ", pid)
	data, err := os.ReadFile(dir)
	if err != nil {
		log.Println("can't read the environ file")
		return nil
	}
	return strings.Split(string(data), "\u0000")

}
