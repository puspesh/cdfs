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
        fmt.Println("Came here .. ");
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
        filename := itemHead[fileIndex+len(lookfor):]
        filename = filename[:strings.Index(filename, "\"")]

        slashIndex := strings.LastIndex(filename, "\\")
        if slashIndex > 0 {
                filename = filename[slashIndex+1:]
        }
        slashIndex = strings.LastIndex(filename, "/")
        if slashIndex > 0 {
                filename = filename[slashIndex+1:]
        }
        _, saveToFilename := filepath.Split(filename)
        //END: work around IE sending full filepath

        //join the filename to the upload dir
        saveToFilePath := filepath.Join("/tmp", saveToFilename)

        osFile, err := os.Create(saveToFilePath)
        if err != nil {
                panic(err.Error())
        }
        defer osFile.Close()

        count, err := io.Copy(osFile, formFile)
        if err != nil {
                panic(err.Error())
        }
        fmt.Printf("ALLOW: %s SAVE: %s (%d)\n", req.RemoteAddr, saveToFilename, count)
        rw.Write([]byte("Upload Complete for " + filename))
}

func main() {
    http.HandleFunc("/upload", upload)
    log.Fatal(http.ListenAndServe(":8082", nil))
}
