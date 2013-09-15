package models

import (
    "fmt"
//    "strconv"
    "encoding/json"
    "code.google.com/p/go-uuid/uuid"
    "github.com/garyburd/redigo/redis"
)

type UserConfigData struct {
    Id      string  `json:"id"`
    Service int     `json:"service"`
    Token   map[string]string  `json:"token"`
    Valid   bool    `json:"valid"`
    Files   map[string]string   `json:"files"`
}

type FileMappingData struct {
    Id      string `json:"id"`
    Parts   map[string][]string  `json:"token"`
}

func getConn() redis.Conn {
    c, err := redis.Dial("tcp", ":6379")
    if err != nil {
        fmt.Println(err)
    }
    return c
}
/*
func main() {
    c := getConn()
    defer c.Close()

    tokens := make(map[string]string)
    tokens["1"] = "adasfasfag@#!24123jlkasdjaklsjd"
    tokens["2"] = "323@#!@3ag@#!24123jlkasdjaklsjd"
    uid := RegisterUser(tokens)
    fmt.Println("New user created - "+uid);

    tokens = make(map[string]string)
    tokens["4"] = "323@#!@3ag@#!24123jlkasdjaklsjd"
    UpdateUser(uid, tokens)
}*/

func GetFileId() string {
  return "f" + uuid.NewRandom().String()
}

func getFileMappingKey(fid string) string {
  return "file_"+fid+"_map" 
}

func FileUploaded(fid string, partsInfo map[string][]string) {
    c := getConn()
    defer c.Close()
        
    _, err := redis.String(c.Do("GET", getFileMappingKey(fid)))
    if err != nil {
      fData := new(FileMappingData)
      fData.Id = fid
      fData.Parts = partsInfo

      fDataSer, err := json.Marshal(fData)
      if err != nil {
        fmt.Println("ERROR! not able to marshal the file metadata ...")
      }
      c.Do("SET", getFileMappingKey(fid), fDataSer)

    } else {
      fmt.Println("Error! File mapping data already exists.");
    }
}

func GetFileMetadata(fid string, part int) string {
    c := getConn()
    defer c.Close()
  
    return ""
}

func UpdateUser(uid string, tokens map[string]string) {
    c := getConn()
    defer c.Close()

    res, err := redis.String(c.Do("GET", "user_"+uid+"_data"))
    if err != nil {
      fmt.Println("Error! Key doesnt exist ... yet.");
      RegisterUser(tokens)
    } else {
      // fmt.Println("Before: "+res)
      ts := new(UserConfigData)
      err := json.Unmarshal([]byte(res), ts)
      if err != nil {
          fmt.Println("Error! not able to read data from redis.")
          return
      }
      for k, v := range tokens {
        if _, ok := ts.Token[k]; !ok {
          ts.Token[k] = v
        }
      }
      b, err := json.Marshal(ts)
      if err != nil {
          fmt.Println(err)
          return
      }
      c.Do("SET", "user_"+uid+"_data", b)
      // fmt.Println("After: "+string(b))
    }
}

func RegisterUser(tokens map[string]string) string {
    c := getConn()
    defer c.Close()

    uid := uuid.NewRandom().String()

    data :=  new(UserConfigData)
    data.Id = uid
    data.Valid = true
    data.Token = tokens

    b, err := json.Marshal(data)
    if err != nil {
        fmt.Println(err)
        return  ""
    }
    c.Send("SET", "user_" + uid + "_data", b)
    c.Flush()

    return uid
}

func GetFileListForUser(userId string) map[string]string {
    c := getConn()
    defer c.Close()

    res, err := redis.String(c.Do("GET", "user_"+userId+"_data"))
    if err != nil {
      fmt.Println("Error! Key doesnt exist ... yet.");
      return nil
    }
    
    ts := new(UserConfigData)
    err = json.Unmarshal([]byte(res), ts)
    if err != nil {
        fmt.Println("Error! not able to read data from redis.")
        return nil
    }
    return ts.Files
}


