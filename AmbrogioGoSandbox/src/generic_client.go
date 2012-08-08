package main

import (
	"./protocol"
	"fmt"
	"math/rand"
	"net"
    "time"
)

func main() {
	// Opening connection with a backend
	conn, err := net.Dial("unix", protocol.FileBE)
	if err != nil {
		fmt.Println("[CLIENT] error:", err)
		return
	}
	defer conn.Close()
	/*
	 * NOTE
	 * the client is responsible of creating a unique ID
	 * for the opened connection
	 */
	uID := uint64(rand.Int63())

	// Create SocketMsg (Correct BEID and FID, with data)
	msg := protocol.CreateSocketMsg(protocol.FileBEID, protocol.FileBE_ListDirectoryContents, uID)
	data := []byte("dir2")
	msg.AddData(data, protocol.PATH)
	var dir protocol.Directory
	// Send SocketMsg
	n, err := msg.Send(conn)
	fmt.Print("Sent ", n, " bytes... ")
	if err != nil {
		fmt.Println("[CLIENT] error: ", err)
	} else {
		fmt.Println("[ OK ]")
        // Receive reply
		n, err = dir.Recv(conn)
		fmt.Print("Read ", n, " bytes... ")
		if err != nil {
			fmt.Println("[CLIENT] error: ", err)
		}
		fmt.Println()
		fmt.Println(dir.Name)
        for _,e := range dir.Entry {
            fmt.Print(e.Name, "\t\t")
            if e.Name == ".." {
                fmt.Println()
                continue
            }
            fmt.Print(e.Size, "\t")
            fmt.Println(time.Unix(e.ModTime,0))
        }
	}

	// Create SocketMsg (Correct BEID, incorrect FID, with data)
	msg = protocol.CreateSocketMsg(protocol.FileBEID, 0, uID)
	data = []byte{0xFF, 0xAA, 0x88, 0x55, 0x00}
	msg.AddData(data, protocol.RAW_DATA)
	// Send SocketMsg
	n, err = msg.Send(conn)
	fmt.Print("Sent ", n, " bytes... ")
	if err != nil {
		fmt.Println("[CLIENT] error: ", err)
	} else {
		fmt.Println("[ OK ]")
	}

	// Create SocketMsg (Incorrect BEID and FID, without data)
	msg = protocol.CreateSocketMsg(0, 0, uID)
	// Send SocketMsg
	n, err = msg.Send(conn)
	fmt.Print("Sent ", n, " bytes... ")
	if err != nil {
		fmt.Println("[CLIENT] error: ", err)
	} else {
		fmt.Println("[ OK ]")
	}
}
