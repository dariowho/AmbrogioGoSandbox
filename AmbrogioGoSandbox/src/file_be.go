package main

import (
	"./common"
	"./protocol"
	"fmt"
	"io"
	"net"
	"os"
	"path"
    "strings"
)

var Preferences struct {
    showHiddenFiles bool
}

func main() {
    initialize()
	//Create socket file
	ls, err := net.Listen("unix", protocol.FileBE)
	go common.SignalCatcher(ls) // <-- needed to catch signals
	if err != nil {
		common.HandleError(err)
		return
	}
	defer ls.Close()

	for {
		//Open socket connection
		conn, err := ls.Accept()
		if conn == nil {
			return
		}
		if err != nil {
			common.HandleError(err)
			return
		}
		go processConnection(conn)
	}
}

/**
 * Init function
 * initialize the environment and load preferences from a file
 */
func initialize() {
    Preferences.showHiddenFiles = false
    return
}

/**
 * Serve proxy connection - rules:
 *  - every new incoming connection is served by a new thread
 *  - requests on an opened connection are served in sequence
 */
func processConnection(conn net.Conn) (err error) {
	var msg protocol.SocketMsg
	var i int
	for {
		//Receive SocketMsg
		i, err = msg.Recv(conn)
		if err == io.EOF {
			break
		}
		if err != nil {
			common.HandleError(err)
			return
		}

		//Debugging
		fmt.Println()
		fmt.Print("Read ", i, " bytes: ")
		fmt.Println(msg)

		//Check if this is the right backend
		if msg.BEID() != protocol.FileBEID {
			protocol.Ack(conn, protocol.ERROR, protocol.WRONG_BE)
            continue
		}

		//Deliver the SocketMsg to the requested function
		switch msg.FID() {
		case protocol.FileBE_ListDirectoryContents:
			ls(conn, msg)
		default:
			protocol.Ack(conn, protocol.ERROR, protocol.UNKNOWN_FUNCTION)
            continue
		}
	}
	return
}

/*
 * List Directory Contents
 *  - 1st arg is the path of the directory
 *  - other arguments will be ignored
 */
func ls(conn net.Conn, msg protocol.SocketMsg) {
	//Check if the first argument is a path
	if msg.ArgType(0) != protocol.PATH {
		protocol.Ack(conn, protocol.ERROR, protocol.WRONG_DATA_TYPE)
		return
	}

	//Clean received path
	dir := path.Clean("/" + string(msg.ArgData(0)))
	//Opening absolute path
	cwd, err := os.Open(protocol.HomeDir + dir)
	if err != nil {
		protocol.Ack(conn, protocol.ERROR, err.Error())
		return
	}
	defer cwd.Close()
	cwdinfo, err := cwd.Stat()
	if err != nil {
		protocol.Ack(conn, protocol.ERROR, err.Error())
		return
	}

	//Check if is a directory
	var directory protocol.Directory
	if cwdinfo.IsDir() {
		directory.Name = dir
		directory.Entry = make([]protocol.FileEntry, 0)
		//Read directory contents
		content, err := cwd.Readdir(0)
		if err != nil {
			protocol.Ack(conn, protocol.ERROR, err.Error())
			return
		}

		if dir != "/" {
            directory.Entry = append(directory.Entry, protocol.FileEntry{"..", 0, cwdinfo.ModTime().Unix(), protocol.YES})
		}
        
		for i := 0; i <= len(content)-1; i++ {
            var e protocol.FileEntry
            e.Name = content[i].Name()
            if strings.HasPrefix(e.Name, ".") && !Preferences.showHiddenFiles {
                continue
            }
            e.Size = content[i].Size()
            e.ModTime = content[i].ModTime().Unix()
            if content[i].IsDir() {
                e.IsDir = protocol.YES
            } else {
                e.IsDir = protocol.NO
            }
			directory.Entry = append(directory.Entry, e)
		}
		directory.EntriesCount = protocol.Size(len(directory.Entry))
	} else { //in case of file
		protocol.Ack(conn, protocol.ERROR, protocol.NOT_A_DIR)
		return
	}

	protocol.Ack(conn, protocol.SUCCESS, "")
	n, err := directory.Send(conn)
	if err != nil {
		common.HandleError(err)
		return
	}

	//Debugging
	fmt.Print("Sent ", n, " bytes: ")
	fmt.Println(directory)
	return
}
