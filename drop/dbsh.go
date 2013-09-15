package main

import (
	"bufio"
	"drop/dropbox"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	gpath "path"
	"strings"
	"text/tabwriter"
	"time"
	"unicode"
)

func Ls(db *dropbox.Client, args []string) error {
	md, e := db.GetMetadata(Cwd, true)
	if e != nil {
		return e
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 1, ' ', 0)
	defer w.Flush()
	for _, f := range md.Contents {
		fmt.Fprintf(w, "%d\t%s\t%s\t\n", f.Bytes, f.ModTime().Format(time.Stamp), gpath.Base(f.Path))
	}
	return nil
}

func Cd(db *dropbox.Client, args []string) error {
	dest := args[0]
	if dest == ".." {
		Cwd = gpath.Dir(Cwd)
		return nil
	}
	dest = mkabs(dest)
	md, e := db.GetMetadata(dest, false)
	if e != nil {
		return e
	}
	if md.IsDir {
		Cwd = dest
		return nil
	}
	return fmt.Errorf("No such dir: %s", dest)
}

func Cat(db *dropbox.Client, args []string) error {
	rc, e := db.GetFile(mkabs(args[0]))
	if e != nil {
		return e
	}
	defer rc.Close()
	if !strings.HasPrefix(rc.ContentType, "text/") {
		return fmt.Errorf("Not a content type you should cat: %s", rc.ContentType)
	}
	_, e = io.Copy(os.Stdout, rc)
	return e
}

func Put(db *dropbox.Client, args []string) error {
	srcfile := args[0]
	if !gpath.IsAbs(srcfile) {
		srcdir, e := os.Getwd()
		if e != nil {
			return e
		}
		srcfile = gpath.Join(srcdir, srcfile)
	}
	src, e := os.Open(srcfile)
	if e != nil {
		return e
	}
	defer src.Close()
	fi, e := src.Stat()
	if e != nil {
		return e
	}
	destpath := gpath.Join(Cwd, gpath.Base(srcfile))
	fmt.Printf("Uploading to %s\n", destpath)
	_, e = db.AddFile(destpath, src, fi.Size())
	return e
}

func Get(db *dropbox.Client, args []string) error {
	fname := mkabs(args[0])
	destdir, e := os.Getwd()
	if e != nil {
		return e
	}
	destfile := gpath.Join(destdir, gpath.Base(fname))
	r, e := db.GetFile(fname)
	if e != nil {
		return e
	}
	defer r.Close()
	fmt.Printf("Saving to %s\n", destfile)
	dest, e := os.Create(destfile)
	if e != nil {
		return e
	}
	defer dest.Close()
	_, e = io.Copy(dest, r)
	return e
}

func Share(db *dropbox.Client, args []string) error {
	link, e := db.GetLink(mkabs(args[0]))
	if e != nil {
		return e
	}
	fmt.Println(link.URL)
	return nil
}

func mkabs(path string) string {
	if !gpath.IsAbs(path) {
		return gpath.Join(Cwd, path)
	}
	return path
}

func Mv(db *dropbox.Client, args []string) error {
	from, to := mkabs(args[0]), mkabs(args[1])
	_, e := db.Move(from, to)
	return e
}

func Cp(db *dropbox.Client, args []string) error {
	from, to := mkabs(args[0]), mkabs(args[1])
	_, e := db.Copy(from, to)
	return e
}

func Rm(db *dropbox.Client, args []string) error {
	_, e := db.Delete(mkabs(args[0]))
	return e
}

func Mkdir(db *dropbox.Client, args []string) error {
	_, e := db.CreateDir(mkabs(args[0]))
	return e
}

func Whoami(db *dropbox.Client, args []string) error {
	ai, e := db.GetAccountInfo()
	if e != nil {
		return e
	}
	b, e := json.MarshalIndent(ai, "", "   ")
	if e != nil {
		return e
	}
	fmt.Println(string(b))
	return nil
}

func Find(db *dropbox.Client, args []string) error {
	r, e := db.Search(Cwd, args[0], 0)
	if e != nil {
		return e
	}
	for _, m := range r {
		fmt.Println(m.Path)
	}
	return nil
}

type Cmd struct {
	Fn       func(*dropbox.Client, []string) error
	ArgCount int
}

func tryCmd(c *dropbox.Client, cname string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Last argument needs to be a dropbox path.")
	}
	fname := args[len(args)-1]
	args = args[:len(args)-1]
	f, e := c.GetFile(mkabs(fname))
	if e != nil {
		return e
	}
	defer f.Close()
	cmd := exec.Command(cname, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = f
	return cmd.Run()
}

var cmds = map[string]Cmd{
	"pwd": {func(*dropbox.Client, []string) error {
		fmt.Println(Cwd)
		return nil
	}, 0},
	"mv":     {Mv, 2},
	"cp":     {Cp, 2},
	"cd":     {Cd, 1},
	"rm":     {Rm, 1},
	"mkdir":  {Mkdir, 1},
	"share":  {Share, 1},
	"cat":    {Cat, 1},
	"ls":     {Ls, 0},
	"find":   {Find, 1},
	"whoami": {Whoami, 0},
	"put": {Put, 1},
	"get": {Get, 1},
	"exit": {func(*dropbox.Client, []string) error {
		os.Exit(0)
		return nil
	}, 0},
}

// Global mutable var. Oh noes.
var Cwd = "/"

func tokenize(s string) []string {
	curr := make([]rune, 0, 32)
	res := make([]string, 0, 10)
	var in_word, last_slash bool
	for _, r := range s {
		if last_slash || (!unicode.IsSpace(r) && (r != '\\')) {
			curr = append(curr, r)
			in_word = true
		} else {
			if in_word {
				res = append(res, string(curr))
				curr = curr[0:0]
				in_word = false
			}
		}
		last_slash = r == '\\'
	}
	if len(curr) > 0 {
		res = append(res, string(curr))
	}
	return res
}

var AppToken = dropbox.AppToken{
	Key:    "y9864ds1g65pug9",
	Secret: "318z58t5z5afjs5",}
var AccessToken = dropbox.AccessToken{
	Secret:    "y9864ds1g65pug9",
	Key: "JkZO0FZyPa8AAAAAAAAAAQpj67CBMgtHTB8Vl2F0pMJqLSFwv5saqLBKa2wZ71Wj",}

func main() {
	db := &dropbox.Client{
		AppToken:    AppToken,
		AccessToken: AccessToken,
		Config: dropbox.Config{
			Access: dropbox.Dropbox,
			Locale: "us",
		}}
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s > ", gpath.Base(Cwd))
		lineb, _, e := in.ReadLine()
		if e != nil {
			if e == io.EOF {
				break
			}
			panic(e)
		}
		tokens := tokenize(string(lineb))
		if len(tokens) == 0 {
			continue
		}

		cmd, ok := cmds[tokens[0]]
		if !ok {
			if e := tryCmd(db, tokens[0], tokens[1:]); e != nil {
				fmt.Printf("ERROR: %v\n", e)
			}
		} else {
			args := tokens[1:]
			if len(args) != cmd.ArgCount && cmd.ArgCount != -1 {
				fmt.Printf("ERROR: %s expected %d args, got %d.\n", tokens[0], cmd.ArgCount, len(args))
			} else if e := cmd.Fn(db, args); e != nil {
				fmt.Printf("ERROR: %v\n", e)
			}
		}
	}
}
