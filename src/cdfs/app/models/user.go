package models

import "github.com/robfig/revel"

type User struct {
        Id              string
        Password        string
        Email           string
}

type FileListData struct {
  name    string
  fid     string 
}

func (user *User) GetFiles() map[string]string {
      ret := GetFileListForUser(user.Id)
      return ret
}

func (user *User) Validate(v *revel.Validation) {
        v.Required(user.Password)
        //v.MinSize(user.Password, 6)
        v.Required(user.Email)
        v.Email(user.Email)

        if user.Email == "test@cdfs.com" && user.Password == "test" {
          user.Id = "6964943e-535a-4736-85ad-4baaa9656709" 
        } else {
          v.Error("Sorry! invalid account ...")
        }
}
