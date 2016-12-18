package main

import (
    "fmt"
    "net/http"
    "html/template"
    "io"
    "os"
    "strconv"
)


func uploadHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")

  for i := 1; i < 3; i++ {

      file, _, err := r.FormFile("file"+strconv.Itoa(i))

      if err != nil {
        fmt.Fprintln(w, err)
        return
      }

      defer file.Close()

      out, err := os.Create("static/tmp/uploadedfile"+strconv.Itoa(i)+".png")
      if err != nil {
        fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
        return
      }

      defer out.Close()

      // write the content from POST to the file
      _, err = io.Copy(out, file)
      if err != nil {
        fmt.Fprintln(w, err)
      }
  }

         testTemplate, err := template.ParseFiles("chat.gohtml")
        if err != nil {
          panic(err)
        }

        vd := ViewData{"tmp/uploadedfile1.png"}
        err = testTemplate.Execute(w, vd)
        if err != nil {
          http.Error(w, err.Error(), http.StatusInternalServerError)
        }
 }

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

var testTemplate *template.Template

type ViewData struct {  
  Name string
}

func main() {  

  fs := http.FileServer(http.Dir("static"))
  http.Handle("/", fs)

  http.HandleFunc("/upload", uploadHandler)
  //http.HandleFunc("/", handler)
  http.ListenAndServe(":3000", nil)
}


func handler(w http.ResponseWriter, r *http.Request) { 
  w.Header().Set("Content-Type", "text/html")
  var err error
  testTemplate, err = template.ParseFiles("index.gohtml")
  if err != nil {
    panic(err)
  }
  vd := ViewData{"John Smith"}
  err = testTemplate.Execute(w, vd)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}