package main

import (
    "fmt"
    "net/http"
    "io"
    "os"
    "bytes"
    "encoding/base64"
     "image"
     "image/png"
     "io/ioutil"
)


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

func uploadHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")

      file, _, err := r.FormFile("file")

      if err != nil {
        fmt.Fprintln(w, err)
        return
      }

      defer file.Close()

      out, err := os.Create("static/tmp/uploadedfile.png")
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

        html := `<html><body>
         <script type="text/javascript" src="https://code.jquery.com/jquery-1.11.3.js"></script>
         <script type="text/javascript" src="https://socketloop.com/public/tutorial/html2canvas.js"></script>

         <style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  background-color: #fff;
}

ul.ChatLog {
  list-style: none;
}

.ChatLog {
  max-width: 20em;
  margin: 0 auto;
}
.ChatLog .ChatLog__entry {
  margin: .5em;
}

.ChatLog__entry {
  display: flex;
  flex-direction: row;
  align-items: flex-end;
  max-width: 100%;
}

.ChatLog__entry.ChatLog__entry_mine {
  flex-direction: row-reverse;  
}

.ChatLog__avatar {
  flex-shrink: 0;
  flex-grow: 0;
  z-index: 1;
  height: 50px;
  width: 50px;
  border-radius: 25px;
  
}

.ChatLog__entry.ChatLog__entry_mine 
.ChatLog__avatar {
  display: none;
}

.ChatLog__entry .ChatLog__message {
  position: relative;
  margin: 0 12px;
}

.ChatLog__entry .ChatLog__message::before {
  position: absolute;
  right: auto;
  bottom: .6em;
  left: -12px;
  height: 0;
  content: '';
  border: 6px solid transparent;
  border-right-color: #ddd;
  z-index: 2;
}

.ChatLog__entry.ChatLog__entry_mine .ChatLog__message::before {
  right: -12px;
  bottom: .6em;
  left: auto;
  border: 6px solid transparent;
  border-left-color: #08f;
}

.ChatLog__message {
  background-color: #ddd;
  padding: .5em;
  border-radius: 4px;
  font-weight: lighter;
  max-width: 70%;
  font-family: Helvetica, Arial, sans-serif;
}

.ChatLog__entry.ChatLog__entry_mine .ChatLog__message {
  border-top: 1px solid #07f;
  border-bottom: 1px solid #07f;
  background-color: #08f;
  color: #fff;
}

.ChatLog__message .ChatLog__timestamp {
  display: none;
}
</style>

  <div id="target-area" style="background-color:white;display:table;">
  <div class="arrow"></div>
  <ul class="ChatLog">
    <li class="ChatLog__entry">
      <img class="ChatLog__avatar" src="tmp/uploadedfile.png" />
      <p class="ChatLog__message">
        Hello!
        <time class="ChatLog__timestamp">6 minutes ago</time>
      </p>
    </li>
    <li class="ChatLog__entry">
      <img class="ChatLog__avatar" src="tmp/uploadedfile.png" />
      <p class="ChatLog__message">
        What is going on here?
        <time class="ChatLog__timestamp">5 minutes ago</time>
      </p>
    </li>
    <li class="ChatLog__entry ChatLog__entry_mine">
      <img class="ChatLog__avatar" src="tmp/uploadedfile.png" />
      <p class="ChatLog__message">
        I have no idea.
        <time class="ChatLog__timestamp">4 minutes ago</time>
      </p>
    </li>
    <li class="ChatLog__entry">
      <img class="ChatLog__avatar" src="tmp/uploadedfile.png" />
      <p class="ChatLog__message">
        I have a neat idea. Maybe I should explain it to you in detail?
        <time class="ChatLog__timestamp">3 minutes ago</time>
      </p>
    </li>
    <li class="ChatLog__entry ChatLog__entry_mine">
      <img class="ChatLog__avatar" src="tmp/uploadedfile.png" />
      <p class="ChatLog__message">
        Sure thing. The more detail the better. In fact, if you could provide definitions for every single term you use, that would be terrific!
        <time class="ChatLog__timestamp">2 minutes ago</time>
      </p>
    </li>
  </ul>
  </div>
  <button type="button" onclick="captureDiv()">Save as Image</button>
         </body></html>
         <script type="text/javascript">
           function captureDiv() {
             html2canvas([document.getElementById('target-area')], {
                  onrendered: function(canvas)
                  {
                     var imgBase64 = canvas.toDataURL() // already base64
                     var img2html = '<img src=' + imgBase64 + '>'
                     $.post("/save", {data: imgBase64}, function () {
                        window.location.href = "save"});
                  }
             });
           }
         </script>`

         w.Write([]byte(html))

 }

func main() {  

    http.HandleFunc("/upload", uploadHandler)
    http.HandleFunc("/save", SaveFile)
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/", fs)

    http.ListenAndServe(":3000", nil)

}
