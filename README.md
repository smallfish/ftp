FTP client for Google Go language 
==================================

install 
========
go get github.com/xiangzhai/goftp

example 
========

sudo apt-get install vsftpd 
cp ~/test.png /home/ftp/

```go
package main

import (                                                                        
    "fmt"                                                                       
    "os"                                                                        
    "github.com/xiangzhai/goftp"                                                
)

const (
    fileName string = "test.png"
)

var (
    received int = 0
)

func callBack(int n) {
    received += n
    fmt.Println("received:", received)
}

func main() {                                                                   
    f, _ := os.Create(fileName)
    defer f.Close()
    // new ftp                                                                  
    ftp := new(ftp.FTP)                                                         
    // set debug, default false                                                 
    ftp.Debug = true
    // set callback function pointer
    ftp.Callback = callBack                                                            
    // connect                                                                  
    ftp.Connect("localhost", 21)                                                
    // login                                                                    
    ftp.Login("anonymous", "")                                                  
    // login failure                                                            
    if ftp.Code == 530 {                                                        
        fmt.Println("error: login failure")                                     
        os.Exit(-1)                                                             
    }                                                                           
    // Switching to Binary mode
    ftp.Request("TYPE I")
    // Directory changed to "/"
    ftp.Cwd("/")
    // Entering Passive Mode
    ftp.Pasv()
    // connect to ftp host with new pasv port
    conn := ftp.NewConnect()
    // download file
    // Restart position accepted (0)
    ftp.Request("REST 0")
    // Opening BINARY mode data connection for test.png (128997 bytes)
    ftp.Request("RETR " + fileName)
    // for looping write to file
    ftp.WriteToFile(conn, f, 0)
    // quit                                                                     
    ftp.Quit()                                                                  
}
```
