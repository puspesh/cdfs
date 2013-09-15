package main

import (
    "log"
    "net/http"
    "fmt"
        "io"
        "os"
        "path/filepath"
        "strings"
)

func upload(rw http.ResponseWriter, req *http.Request) {
        rw.Header().Set("Access-Control-Allow-Origin", "*")
        fmt.Println("Came here .. ");

        vals := req.URL.Query()
        filenamesent := vals["file"][0]
        fileParts := vals["num"][0]
        fileToken := vals["token"][0]

        formFile, formHead, err := req.FormFile("TheFile")
        if err != nil {
          fmt.Println(err)
          return
        }
        defer formFile.Close()
        fmt.Println("Yay!")

        itemHead := formHead.Header["Content-Disposition"][0]
        lookfor := "filename=\""
        fileIndex := strings.Index(itemHead, lookfor)
        if fileIndex < 0 {
                panic("runUp: no filename")
        }
        /*
        filename := itemHead[fileIndex+len(lookfor):]
        filename = filename[:strings.Index(filename, "\"")]
        */
        filename := fileToken +"_"+filenamesent+"_"+fileParts
        fmt.Println("filename == "+filename);

        //END: work around IE sending full filepath

        //join the filename to the upload dir
        saveToFilePath := filepath.Join("/tmp/", filename)

        osFile, err := os.Create(saveToFilePath)
        if err != nil {
                panic(err.Error())
        }
        defer osFile.Close()

        count, err := io.Copy(osFile, formFile)
        if err != nil {
                panic(err.Error())
        }
        fmt.Printf("ALLOW: %s SAVE: %s (%d)\n", req.RemoteAddr, filename, count)
        rw.Write([]byte("Upload Complete for " + filename))
}

func main() {
    http.HandleFunc("/upload", upload)
    log.Fatal(http.ListenAndServe(":8082", nil))
}
