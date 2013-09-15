package models

import "strconv"
import "fmt"

func GetApiHandler(store string, auth string) *GoogleDrive {
//	if store == "google" {
		return &GoogleDrive{Store: store, authString: auth}
//	}
}

func UploadFiles(location string, fileName string, parts int, u *UserConfigData, f *FileMappingData) {
	path := location + "/" + fileName
	clients := []GoogleDrive{}
	for key, value := range u.Token {
		clients = append(clients, *GetApiHandler(key, value))
	}
	totClients := len(clients)
	for i := 0; i < parts; i++ {
		fi := path + "_" + strconv.Itoa(i)
		fn := fileName + "_" + strconv.Itoa(i)
		cop := 0
		loc := []string{}
		f.Parts[strconv.Itoa(i)] = loc
		for j := 0; j < totClients && cop < 2; j++ {
			index := (i + j) % totClients
			client := clients[index]
                fmt.Println("file is %s", fi)
			if err,val := client.CheckSize(); err == nil && val {
                fmt.Println("file is %s", fi)
				client.Upload(fi, fn)
				loc = append(loc, client.Store)
				cop++
			} else {
                if err != nil {
                    fmt.Println("upload error", err)
                }
            }
		}
	}
}

/*
func main() {
    u := UserConfigData{Token:make(map[string]string),}
    f := FileMappingData{Parts:make(map[string][]string),}
    u.Token["google"]="token"
    UploadFiles("test", "output.txt", 9, &u, &f)
}*/
