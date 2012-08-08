/**
 * protocol Package
 * Handle communication between Backends using socket and stardand messages
 * 
 * Copyright (c) 2012, Alessandro Falsone
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the <organization> nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

// Standard reply
const (
	SUCCESS uint8 = 1
	ERROR   uint8 = 0
)

// Error messages
const (
	WRONG_BE         string = "wrong backend"
	UNKNOWN_FUNCTION string = "function not found"
	WRONG_DATA_TYPE  string = "wrong data type"
)

// Data Types
const (
	RAW_DATA uint8 = 0
	UINT8    uint8 = 1
	/**
		 * UINT16 uint8 = 2
	     * ...
	*/
	STRING uint8 = 10
	PATH   uint8 = 11
	// ...
)

type Size uint16

/**
 * NOTA 1
 * I campi che iniziano con una lettera minuscola non vengono
 * esportati, in questo modo l'utente non può fare casini, o
 * meglio, deve impegnarsi per farli.
 * Inoltre incapsulare le informazioni nell'header ci permette di:
 *  1) poterlo trattare come un flusso unico e inviarlo in blocco
 *     con le funzioni del pkg "encoding/binary"
 *  2) nascondere i campi BEID, FID e UID all'utente che invece
 *     devono essere esportati per poter essere trattati dalle
 *     funzioni esterne
 *
 * NOTA 2
 * I campi argc e dsize sono sovrabbondanti perché si può sempre
 * risalire alla lunghezza delle slices, li ho lasciati per un maggiore
 * controllo dato che comunque sono modificati solo dalle funzioni e
 * in accordo alla quantità di dati rilevata
 */
type SocketMsg struct {
	header struct {
		BEID uint8  // ID of called BE
		FID  uint8  // ID of function called
		UID  uint64 // Unique ID of request/response
	}
	argc uint8 // number of argument in slice
	argv []arg // argument list
}

type arg struct {
	dtype uint8  // data type
	dsize Size   // data size in bytes
	data  []byte // raw data slice
}

/**
 * Create SocketMsg
 */
func CreateSocketMsg(beID, fID uint8, uID uint64) (msg SocketMsg) {
	msg.header.BEID = beID
	msg.header.FID = fID
	msg.header.UID = uID
	msg.argc = 0
	msg.argv = nil
	return
}

/**
 * Append data to msg
 */
func (msg *SocketMsg) AddData(data []byte, dtype uint8) {
	//Create new arg
	var newarg arg
	newarg.dtype = dtype
	newarg.dsize = Size(len(data))
	newarg.data = make([]byte, len(data))
	copy(newarg.data, data)
	//Append new arg and update msg
	msg.argv = append(msg.argv, newarg)
	msg.argc = uint8(len(msg.argv))
	return
}

/**
 * Get BEID
 */
func (msg *SocketMsg) BEID() (beID uint8) {
	beID = msg.header.BEID
	return
}

/**
 * Get FID
 */
func (msg *SocketMsg) FID() (fID uint8) {
	fID = msg.header.FID
	return
}

/**
 * Get UID
 */
func (msg *SocketMsg) UID() (uID uint64) {
	uID = msg.header.UID
	return
}

/**
 * Get argv[i].dtype
 */
func (msg *SocketMsg) ArgType(i uint8) (t uint8) {
	if i < 0 || i >= msg.argc {
		return RAW_DATA
	}
	t = msg.argv[i].dtype
	return
}

/**
 * Get argv[i].data
 */
func (msg *SocketMsg) ArgData(i uint8) (d []byte) {
	if i < 0 || i >= msg.argc {
		return nil
	}
	d = make([]byte, msg.argv[i].dsize)
	copy(d, msg.argv[i].data)
	return
}

/**
 * Send msg to passed connection, number of bytes sent or err on error
 */
func (msg *SocketMsg) Send(conn net.Conn) (n int, err error) {
	// Writing Header
	n, err = BinaryWrite(conn, msg.header)
	if err != nil {
		return
	}
	// Writing argc
	m, err := conn.Write([]byte{msg.argc})
	n = n + m
	if err != nil {
		return
	}
	// Writing argv
	for i := uint8(0); i < msg.argc; i++ {
		// Writing dtype
		m, err = conn.Write([]byte{msg.argv[i].dtype})
		n = n + m
		if err != nil {
			return
		}
		// Writing dsize
		m, err = BinaryWrite(conn, msg.argv[i].dsize)
		n = n + m
		if err != nil {
			return
		}
		// Writing data
		m, err = conn.Write(msg.argv[i].data)
		n = n + m
		if err != nil {
			return
		}
	}
	// Wait for be's ack
	b := make([]byte, 1)
	conn.Read(b)
	switch uint8(b[0]) {
	case SUCCESS:
		return
	case ERROR:
		var s string
		ReadString(conn, &s)
		return n, errors.New(s)
	}
	return
}

/**
 * Receive msg on passed connection, returns number of bytes read or err on error
 */
func (msg *SocketMsg) Recv(conn net.Conn) (n int, err error) {
	// Reading Header
	n, err = BinaryRead(conn, Size(binary.Size(msg.header)), &(msg.header))
	if err != nil {
		return
	}
	// Reading argc
	b := make([]byte, 1)
	m, err := conn.Read(b)
	msg.argc = b[0]
	n = n + m
	if err != nil {
		return
	}
	// Preparing the slice
	msg.argv = make([]arg, msg.argc)
	// Reading argv
	var data []byte
	for i := uint8(0); i < msg.argc; i++ {
		// Reading dtype
		m, err = conn.Read(b)
		msg.argv[i].dtype = b[0]
		n = n + m
		if err != nil {
			return
		}
		// Reading dsize
		m, err = BinaryRead(conn, Size(binary.Size(msg.argv[i].dsize)), &(msg.argv[i].dsize))
		n = n + m
		if err != nil {
			return
		}
		// Reading data
		data = make([]byte, msg.argv[i].dsize)
		m, err = conn.Read(data)
		msg.argv[i].data = make([]byte, msg.argv[i].dsize)
		copy(msg.argv[i].data, data)
		n = n + m
		if err != nil {
			return
		}
	}
	return
}

/**
 * Generic type writing function
 * send generic data over passed connection and return number of byte sent or err on error
 */
func BinaryWrite(conn net.Conn, data interface{}) (n int, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, data)
	if err != nil {
		return 0, err
	}
	n, err = conn.Write(buf.Bytes())
	return
}

/**
 * Generic type reading function
 * read generic data over passed connection and return number of byte read or err on error
 * data must be a pointer of a data structure of length equal to len
 */
func BinaryRead(conn net.Conn, len Size, data interface{}) (n int, err error) {
	buf := make([]byte, len)
	n, err = conn.Read(buf)
	if err != nil {
		return
	}
	err = binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, data)
	return
}

/**
 * Writing strings over a socket connection
 * send generic message over passed connection and return number of byte sent or err on error
 */
func WriteString(conn net.Conn, s string) (n int, err error) {
	buf := []byte(s)
	buf = append(buf, '\x00')
	n, err = conn.Write(buf)
	return
}

/**
 * Reading strings over a socket connection
 * send generic message over passed connection and return number of byte sent or err on error
 */
func ReadString(conn net.Conn, s *string) (n int, err error) {
	b := make([]byte, 1)
	buf := make([]byte, 0)
	var m int
	for {
		m, err = conn.Read(b)
		n = n + m
		if b[0] == '\x00' {
			break
		} else {
            buf = append(buf, b[0])
		}
	}
	*s = string(buf)
	return
}

/**
 * Ack the SocketMsg
 * send the ack for the received message
 */
func Ack(conn net.Conn, result uint8, err string) {
	if result == SUCCESS {
		conn.Write([]byte{byte(SUCCESS)})
	} else {
		conn.Write([]byte{byte(ERROR)})
		WriteString(conn, err)
	}
	return
}
