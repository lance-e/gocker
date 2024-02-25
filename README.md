## gocker: A implement of container like RunC

#### Status:
namespace :
- ip namespace: **OK**
- ipc namespace: **OK**
- uts namespace: **OK**
- pid namespace: **OK**
- mnt namespace: **OK**
- user namespace: **TODO**

cgroup:
- memory: **OK**
- cpushare: **OK**
- cpuset: **OK**

unionFS:
- overlayFS: **OK**

some command:
- run:**OK**
    

    run's flags:
    - i(interactive,used with t)
    - t(tty,used with i)
    - d(detach)
    - name
    - m(memory)
    - cpuset
    - cpushare
    - e(environment)
    - network
    - p(port)
    - v(volume)
- commit: **OK**
- exec: **OK**
- stop: **OK**
- rm: **OK**
- logs: **OK**
- ps: **OK**
- restart: **TODO**
- network:**YES OR NO**  

    (this command was contain some bug ,could create the virtually network device ,but could't communicate with the outside of container.)
  
    network's subcommand:
    - create:**OK**
        
        create's flags:
        - driver
        - subnet
    - list:**OK**
    - rm:**OK**
#### install:
~~~bash
git clone https://github.com/lance-e/gocker.git
cd gocker
mv ./image/busybox.tar /root/
~~~
(attention:should put the image's tar file in the **/root/**.you could try use docker export the image into **/root/**,and use gocker to build container by this image)
#### simply display
~~~bash
go build .
root@localhost:/home/longxu/gocker# go build .
root@localhost:/home/longxu/gocker# ./gocker run -it --name=testcontainer busybox sh
2024/02/25 12:41:44 cgroup set...
2024/02/25 12:41:44 cgroup apply...
2024/02/25 12:41:44 init command begin
2024/02/25 12:41:44 this is the data of pipe to translate : sh
2024/02/25 12:41:44 the current location is /root/overlay/testcontainer/merged
/ # ps
PID   USER     TIME  COMMAND
    1 root      0:00 sh
    7 root      0:00 ps
/ # mount
overlay on / type overlay (rw,relatime,lowerdir=/root/overlay/busybox,upperdir=/root/overlay/testcontainer/upper,workdir=/root/overlay/testcontainer/work)
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
tmpfs on /dev type tmpfs (rw,nosuid,mode=755)
/ # ls
bin   dev   etc   home  proc  root  sys   tmp   usr   var
/ # 
~~~

#### tutorial :
[the detail tutorial](https://lance-e.github.io/%E4%BB%8E0%E5%88%B01%E5%86%99docker%E4%B9%8B%E4%B8%80%E6%A6%82%E5%BF%B5%E7%90%86%E8%A7%A3/)