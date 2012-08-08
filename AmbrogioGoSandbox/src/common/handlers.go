package common

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

/**
 * Handle errors
 */
func HandleError(err error) {
	//Do something...
	fmt.Println("[BACKEND] error:", err)
	return
}

/**
 * Handle SIGINT signal ONLY, and then delete socket file
 */
func SignalCatcher(ls net.Listener) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch
	fmt.Println(" SIGINT catched; exiting...")
	if ls != nil {
		ls.Close()
	}
	os.Exit(0)
}
