package controllers

import (
  "github.com/robfig/revel"
  "cdfs/app/models"
  "strconv"
)

type UploadData struct {
      numParts      int
      parts         map[string]string
}

type DownloadData struct {
      filename     string
      parts        int
      urls         map[string]string
}

type App struct {
	*revel.Controller
}

func (c App) Up() revel.Result {
	return c.Render()
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Register() revel.Result {
        
      return c.Render()
}

func (c App) Upload(fname string, fid string, n string) revel.Result {
    user := new(models.User)
    user.Id = "6964943e-535a-4736-85ad-4baaa9656709"

    // Put the upload logic here 
    // and then call FileUploaded() to add the metadata
    u := models.UserConfigData{Token:make(map[string]string),}
    f := models.FileMappingData{Parts:make(map[string]map[string]string),}
    u.Token["google"]="/tmp/token"

    nn, _ := strconv.Atoi(n)
    models.UploadFiles("/tmp/"+fid, fname, nn, &u, &f)

    models.FileUploaded(user.Id, f.Parts, fid, fname)

    return c.Render()
}

func getGoogleUrl(u string) string {
  return ""
}

func (c App) Download(fid string) revel.Result {
  // return info to which nodes do this needs to download 
  // file blocks from and merge
  d := models.GetFileMetadata(fid)

  returnData := new(DownloadData)
  returnData.filename = ""
  returnData.parts = 1
  returnData.urls = make(map[string]string)
  for k, v := range d {
    u := ""
    for l, m := range v {
      if(l == "google") {
        u = getGoogleUrl(m)
        break;
      }
    }
    returnData.urls[k] = u
  }
	return c.RenderJson(returnData)
}

func (c App) List() revel.Result {
  user := new(models.User)
  user.Id = "6964943e-535a-4736-85ad-4baaa9656709"
	return c.Render(user)
}

func (c App) Login(user *models.User) revel.Result {
        user.Validate(c.Validation)

        // Handle errors
        if c.Validation.HasErrors() {
                c.Validation.Keep()
                c.FlashParams()
                return c.Redirect(App.Index)
        }
        // Ok, display the created user
        return c.Redirect(App.List)
        // return c.Render(user)	
}
