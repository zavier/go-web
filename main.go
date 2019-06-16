package main

import (
	"./session"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var globalSessions *session.Manager

func init() {
	globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
	go globalSessions.GC()
}

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("Key:", k)
		fmt.Println("val:", strings.Join(v, ","))
	}
	cookie, err := r.Cookie("username")
	if err != nil && cookie != nil {
		fmt.Fprintf(w, "Hello %s!", cookie.Value)
	} else {
		fmt.Fprintln(w, "Hello welcome!")
	}
}

func generateToken() string {
	curTime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(curTime, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token
}

// 登陆
func login(w http.ResponseWriter, r *http.Request) {
	session := globalSessions.SessionStart(w, r)
	if r.Method == "GET" {
		token := generateToken()

		t, err := template.ParseFiles("login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, token)
	} else {
		username := r.FormValue("username")
		password := r.FormValue("password")
		fmt.Println("username:", username)
		fmt.Println("password", password)

		session.Set("username", username)

		expiration := time.Now()
		expiration = expiration.AddDate(1, 0, 0)
		cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
		http.SetCookie(w, &cookie)
		fmt.Fprintf(w, "welcome, %s", r.FormValue("username"))
	}
}

func count(w http.ResponseWriter, r *http.Request) {
	session := globalSessions.SessionStart(w, r)
	createtime := session.Get("createtime")
	if createtime == nil {
		session.Set("createtime", time.Now().Unix())
	} else if createtime.(int64)+360 < (time.Now().Unix()) {
		globalSessions.SessionDestroy(w, r)
		session = globalSessions.SessionStart(w, r)
	}
	ct := session.Get("countnum")
	if ct == nil {
		session.Set("countnum", 1)
	} else {
		session.Set("countnum", (ct.(int) + 1))
	}
	t, _ := template.ParseFiles("count.html")
	w.Header().Set("Content-type", "text/html")
	t.Execute(w, session.Get("countnum"))
}

// 文件上传
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		token := generateToken()

		t, _ := template.ParseFiles("upload.html")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func main() {
	http.HandleFunc("/say", sayHelloName)
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
