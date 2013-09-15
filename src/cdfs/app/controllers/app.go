package controllers

import "github.com/robfig/revel"
import "cdfs/app/models"
import "fmt"

type UploadData struct {
      numParts      int
      parts         map[string]string
}

type DownloadData struct {
      numParts      int
      parts         map[string]string
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

func (c App) Upload(fid string, n string) revel.Result {
    user := new(models.User)
    user.Id = "6964943e-535a-4736-85ad-4baaa9656709"

    // Put the upload logic here 
    // and then call FileUploaded() to add the metadata

    
    return c.Render()
}
func (c App) Download(fid string) revel.Result {
  // return info to which nodes do this needs to download 
  // file blocks from and merge
  
	return c.Render()
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
