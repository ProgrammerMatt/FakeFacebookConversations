package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io"
    "os"
    "bytes"
    "encoding/base64"
     "image"
     "image/png"
     "io/ioutil"
     "log"
    "html/template"
)

var testTemplate *template.Template


type ViewData struct {
  Filename string
}

func SaveFile(w http.ResponseWriter, r *http.Request) {

         if r.Method == "POST" {

                 imgBase64 := r.FormValue("data")

                 // remove "data:image/png;base64,"
                 imgBase64cleaned := imgBase64[len("data:image/png;base64,"):len(imgBase64)]

                 // decode base64 to buffer bytes
                 imgBytes, _ := base64.StdEncoding.DecodeString(imgBase64cleaned)

                 // convert []byte to image for saving to file
                 img, _, _ := image.Decode(bytes.NewReader(imgBytes))

                 imgFile, err := os.Create("./screen-capture.png")
                 if err != nil {
                         panic(err)
                 }

                 // save to file on your webserver
                 png.Encode(imgFile, img)
                 fmt.Println("screen-capture.png file saved")

         }

         // NOTE : For some odd reason this part has to be outside the r.Method == "POST"
         //        in order for the streaming/download straight to browser to work

         streamBytes, _ := ioutil.ReadFile("./screen-capture.png")

         //output to browser
         w.Header().Set("Content-Disposition", "attachment; filename='screen-capture.png'")
         w.Header().Set("Content-Type", "image/png")
         w.Header().Set("Content-Transfer-Encoding", "binary")
         w.Header().Set("Content-Description", "File Transfer")

         //w.Header().Set("Content-Length", string(len(streamBytes)))
         // will produce warning message : http: invalid Content-Length of "뚴"

         // ok, pump it out to the browser!
         b := bytes.NewBuffer(streamBytes)
         if _, err := b.WriteTo(w); err != nil {
                 fmt.Fprintf(w, "%s", err)
         }

 }

    type Inventory struct {
        Material string
        Count    uint
      }

func JSONWriter(w http.ResponseWriter, val interface{}) {
    w.Header().Set("Content-Type", "application/json")
    b, _ := json.Marshal(val)
    w.Write(b)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")

      file, _, err := r.FormFile("file")

      if err != nil {
        fmt.Fprintln(w, err)
        return
      }

      defer file.Close()

      filename := "tmp/uploadedfile.png"

      out, err := os.Create("static/"+filename)
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

   
      // vd := ViewData{filename}
      // tmpl, err := template.New("master").Delims("<<", ">>").ParseFiles("index.gohtml")
      // if err != nil { panic(err) }
      // err = tmpl.ExecuteTemplate(w, "index.gohtml", vd)
      // if err != nil { panic(err) }

      JSONWriter(w, filename)

       }


 var results []string  

type Data struct {
  Messages [] struct{
      User string
      Msg string
  }
}


// PostHandler converts post request body to string
func generateConversationFromJSON(rw http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
      var cd Data   
      err := decoder.Decode(&cd)
      if err != nil {
          panic(err)
      }
      defer req.Body.Close()
      log.Println(cd.Messages)
}

func main() {  

    http.HandleFunc("/upload", uploadHandler)
    http.HandleFunc("/save", SaveFile)
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/", fs)
    http.HandleFunc("/generate", generateConversationFromJSON)


    http.ListenAndServe(":3000", nil)

}
