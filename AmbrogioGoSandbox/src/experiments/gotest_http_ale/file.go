package main

import (
	"net/http"
	"fmt"
	"os"
	"path"
)

func ls_dir_cont(w http.ResponseWriter, req *http.Request) {
	val := req.URL.Query()					//processing URL's parameters and validating path
	dir := path.Clean("/"+val.Get("cd"))


	cwd, err := os.Open(root_dir+dir)			//opening absolute path
	defer cwd.Close()
	if err != nil {
		fmt.Fprintln(w, "Error: ", err)
		return
	}
	
	fmt.Fprintln(w, top_template)
	cwdinfo, err := cwd.Stat()
	if err != nil {
		fmt.Fprintln(w, "Error: ", err)
		return
	}
	if cwdinfo.IsDirectory() {					//checking if the file is a directory
		content, err := cwd.Readdir(0)				//reading directory content
		if err != nil {
			fmt.Fprintln(w, "Error: ", err)
			return
		}
		
		if dir != "/" {
			fmt.Fprintf(w, "<a href=\"list?cd=%s\">..</a></br>", dir+"/..")
		}
		for i:=0; i<=len(content)-1; i++ {
			if dir != "/" {
				fmt.Fprintf(w, "<a href=\"list?cd=%s\">%s</a></br>", dir+"/"+content[i].Name, content[i].Name)
			} else {
				fmt.Fprintf(w, "<a href=\"list?cd=%s\">%s</a></br>", dir+content[i].Name, content[i].Name)
			}
		}
	} else {							//in case of file
		http.ServeFile(w, req, root_dir+dir)
	}
	fmt.Fprintln(w, bot_template)

	return
}

func main() {
	http.HandleFunc("/list", ls_dir_cont)
	http.Handle("/", http.FileServer(http.Dir("/var/www/gotest")))
	http.ListenAndServe(":8080", nil)
}

const top_template = "<html><head><title>List Directory Content</title></head><body><div align=\"center\">"
const bot_template = "</div></body></html>"

const root_dir = "mydir"
