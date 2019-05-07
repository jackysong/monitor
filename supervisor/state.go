package supervisor

import (
	"log"
	"strings"
)

type State string

const (
	ADDED    State = "PROCESS_GROUP_ADDED"    //新增Supervisord的进程组
	REMOVED  State = "PROCESS_GROUP_REMOVED"  //删除Supervisord的进程组
	UNKNOWN  State = "PROCESS_STATE_UNKNOWN"  //进程状态未知
	STARTING State = "PROCESS_STATE_STARTING" //进程状态从其他状态转换为正在启动
	RUNNING  State = "PROCESS_STATE_RUNNING"  //进程状态由正在启动转换为正在运行
	STOPPING State = "PROCESS_STATE_STOPPING" //进程状态由正在运行转换为正在停止
	STOPPED  State = "PROCESS_STATE_STOPPED"  //进程状态由正在停止转换为已经停止(人为控制退出)
	EXITED   State = "PROCESS_STATE_EXITED"   //进程状态由正在运行转换为退出(程序自行退出)
	BACKOFF  State = "PROCESS_STATE_BACKOFF"  //进程状态由正在启动转换为失败
	FATAL    State = "PROCESS_STATE_FATAL"    //进程状态由正在运行转换为失败
)

func (s State) Source() string {
	return string(s)
}

func (s State) String() string {
	temp := strings.Split(string(s), "_")
	return temp[len(temp)-1]
}

func getState(s string) State {
	var State State
	switch strings.ToUpper(s) {
	case ADDED.String():
		State = ADDED
	case REMOVED.String():
		State = REMOVED
	case STARTING.String():
		State = STARTING
	case RUNNING.String():
		State = RUNNING
	case STOPPING.String():
		State = STOPPING
	case STOPPED.String():
		State = STOPPED
	case EXITED.String():
		State = EXITED
	case BACKOFF.String():
		State = BACKOFF
	case FATAL.String():
		State = FATAL
	default:
		State = UNKNOWN
		log.Println(UNKNOWN, ":", s)
	}
	return State
}
