FTP client for Google Go language 
==================================

install 
========
go get github.com/smallfish/ftp.go

example 
========
```go
package main

import (                                                                        
    "fmt"                                                                       
    "os"                                                                        
    "github.com/smallfish/ftp.go"
)

func main() {                                                                   
    ftp := new(ftp.FTP)                                                         
    // debug default false
    ftp.Debug = true
    ftp.Connect("localhost", 21)                                                

    // login
    ftp.Login("anonymous", "")
    if ftp.Code == 530 {                                                         
        fmt.Println("error: login failure")                                     
        os.Exit(-1)                                                             
    }                                                                           
    
    // pwd
    ftp.Pwd()
    fmt.Println("code:", ftp.Code, ", message:", ftp.Message)                   

    // make dir
    ftp.Mkd("/path")
    ftp.Request("TYPE I")

    // stor file
    b, _ := ioutil.ReadFile("/path/a.txt")
    ftp.Stor("/path/a.txt", b)
    
    ftp.Quit()                                                                  
}
```
