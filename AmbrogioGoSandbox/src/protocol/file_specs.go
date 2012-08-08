package protocol

import (
    "encoding/binary"
	"net"
)

// General informations
const HomeDir string = "test_home"
const YES = 1
const NO = 0

// Socket informations
const FileBE string = "file"
const FileBEID uint8 = 3

// Error messages
const (
	NOT_A_DIR string = "not a directory"
)

// Exported functions
const (
	/**
	 * FileBE_NewFile uint8 = 0
	 * FileBE_NewFolder uint8 = 1
	 * ...
	 */
	FileBE_ListDirectoryContents uint8 = 8
)

// Returned data types
type Directory struct {
	Name         string
	EntriesCount Size
	Entry        []FileEntry
}

type FileEntry struct {
	Name string
    Size    int64
    ModTime int64
    IsDir   uint8
}


// Functions to return data over socket connection

/**
 * Write the directory contents to the passed connection
 */
func (dir *Directory) Send(conn net.Conn) (n int, err error) {
	// Writing directory name
	n, err = WriteString(conn, dir.Name)
	if err != nil {
		return
	}
	// Writing number of entries
	m, err := BinaryWrite(conn, dir.EntriesCount)
	n = n + m
	if err != nil {
		return
	}
	// Writing each entry
	for _,e := range dir.Entry {
		// Writing entry name
		m, err = WriteString(conn, e.Name)
		n = n + m
		if err != nil {
			return
		}
		// Writing entry size
		m, err = BinaryWrite(conn, e.Size)
		n = n + m
		if err != nil {
			return
		}
		// Writing entry time
		m, err = BinaryWrite(conn, e.ModTime)
		n = n + m
		if err != nil {
			return
		}
		// Writing entry isdir
		m, err = BinaryWrite(conn, e.IsDir)
		n = n + m
		if err != nil {
			return
		}
	}
	return
}

/**
 * Read a directory contents from the passed connection
 */
func (dir *Directory) Recv(conn net.Conn) (n int, err error) {
	// Reading directory name
	n, err = ReadString(conn, &(dir.Name))
	if err != nil {
		return
	}
	// Reading number of entries
	m, err := BinaryRead(conn, Size(binary.Size(dir.EntriesCount)), &(dir.EntriesCount))
	n = n + m
	if err != nil {
		return
	}
    dir.Entry = make([]FileEntry, dir.EntriesCount)
	// Reading each entry
	for i,e := range dir.Entry {
		// Reading entry name
		m, err = ReadString(conn, &(e.Name))
		n = n + m
		if err != nil {
			return
		}
		// Reading entry size
		m, err = BinaryRead(conn, Size(binary.Size(e.Size)), &(e.Size))
		n = n + m
		if err != nil {
			return
		}
		// Reading entry time
		m, err = BinaryRead(conn, Size(binary.Size(e.ModTime)), &(e.ModTime))
		n = n + m
		if err != nil {
			return
		}
		// Reading entry isdir
		m, err = BinaryRead(conn, Size(binary.Size(e.IsDir)), &(e.IsDir))
		n = n + m
		if err != nil {
			return
		}
        dir.Entry[i] = e
	}
	return
}
