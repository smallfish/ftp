FTP client for Google Go language 
==================================

install 
========
go get github.com/xiangzhai/goftp

example 
========
```go
package main

import (                                                                        
    "fmt"                                                                       
    "os"                                                                        
    "github.com/xiangzhai/goftp"                                                
)

func main() {                                                                   
    // new ftp                                                                  
    ftp := new(ftp.FTP)                                                         
    // set debug, default false                                                 
    ftp.Debug = true                                                            
    // connect                                                                  
    ftp.Connect("localhost", 21)                                                
    // login                                                                    
    ftp.Login("anonymous", "")                                                  
    // login failure                                                            
    if ftp.Code == 530 {                                                        
        fmt.Println("error: login failure")                                     
        os.Exit(-1)                                                             
    }                                                                           
    // pwd                                                                      
    ftp.Pwd()                                                                   
    fmt.Println("code:", ftp.Code, ", message:", ftp.Message)                   
    // quit                                                                     
    ftp.Quit()                                                                  
}
```
