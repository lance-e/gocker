package setns
/*
#define _GNU_SOURCE
#include "errno.h"
#include "string.h"
#include "stdlib.h"
#include "stdio.h"
#include "sched.h"
#include "fcntl.h"
#include "unistd.h"
//__attribute__((constructor)) 相当于init，让函数提前运行
__attribute__((constructor))  void enter_namespace(void ){
    char *pid;
    pid  = getenv("gocker_pid");
    if (pid){
        //fprintf(stdout,"get the gocker_pid %s\n",pid);
    }else {
        //fprintf(stdout,"get the gocker_pid failed\n");
        return;
    }
    char * cmd;
    cmd = getenv("gocker_cmd");
    if (cmd){
        //fprintf(stdout,"get the gocker_cmd %s\n",cmd);
    }else {
        //fprintf(stdout,"get the gocker_cmd failed\n");
        return;
    }
    int i ;
    char nspath [1024];
    char *namespace []={"pid","net","ipc","uts","mnt"};
    for (i =0;i<5;i++){
        sprintf(nspath,"/proc/%s/ns/%s",pid,namespace[i]);
        int fd =open(nspath,O_RDONLY);
        if( (setns(fd,0))== -1){
            fprintf(stderr,"setns on %s namespace failed,error:%s\n",namespace[i],strerror(errno));
        }else {
            fprintf(stdout,"setns on %s namespace successful\n",namespace[i]);
        }
        close(fd);
    }
    int res = system(cmd);
    exit(0);
    return;
    
}
*/
import "C"