// FTP Client for Google Go language.
// Author: smallfish <smallfish.xy@gmail.com>

/*                                                                              
 * Copyright (C) 2013 Deepin, Inc.                                                 
 *               2013 Leslie Zhai <zhaixiang@linuxdeepin.com>                   
 *                                                                              
 * This program is free software: you can redistribute it and/or modify         
 * it under the terms of the GNU General Public License as published by         
 * the Free Software Foundation, either version 3 of the License, or            
 * any later version.                                                           
 *                                                                              
 * This program is distributed in the hope that it will be useful,              
 * but WITHOUT ANY WARRANTY; without even the implied warranty of               
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the                
 * GNU General Public License for more details.                                 
 *                                                                              
 * You should have received a copy of the GNU General Public License            
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.        
 */

package ftp

import (
	"fmt"
	"net"
	"strconv"
	"strings"
    "os"
    "io"
)

type FTP struct {
	host        string
	port        int
	user        string
	passwd      string
	pasv        int
	cmd         string
	Code        int
	Message     string
	Debug       bool
	stream      []byte
	conn        net.Conn
	Error       error
    offset      int
    Callback    func(int)
}

func (ftp *FTP) debugInfo(s string) {
	if ftp.Debug {
		fmt.Println(s)
	}
}

func (ftp *FTP) NewConnect() net.Conn {
    conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ftp.host, ftp.pasv))
    if err != nil {
        ftp.Error = err
        fmt.Println("ERROR:", ftp.Error.Error())
    }
    ftp.cmd = fmt.Sprintf("NewConnect:%d", ftp.pasv)
    return conn
}

func (ftp *FTP) Connect(host string, port int) bool {
	addr := fmt.Sprintf("%s:%d", host, port)
	ftp.conn, ftp.Error = net.Dial("tcp", addr)
	if ftp.Error != nil {
        fmt.Println("ERROR:", ftp.Error.Error())
        return false
    }
    ftp.Response()
	ftp.host = host
	ftp.port = port
    return true
}

func (ftp *FTP) Login(user, passwd string) {
	ftp.Request("USER " + user)
	ftp.Request("PASS " + passwd)
	ftp.user = user
	ftp.passwd = passwd
}

func (ftp *FTP) WriteToFile(conn net.Conn, f *os.File, offset int) {
    ftp.offset = offset
    defer conn.Close()
    data := make([]byte, 102400)
	for {
        n, err := conn.Read(data)
        if err != nil {
            if err != io.EOF { panic(err) }
        }
        f.WriteAt(data[:n], int64(ftp.offset))
        if ftp.Callback != nil {
            ftp.Callback(ftp.offset)
        }
        ftp.offset += n
        if err == io.EOF { return }
    }
    return
}

func (ftp *FTP) Response() (code int, message string) {
	ret := make([]byte, 1024)
	n, _ := ftp.conn.Read(ret)
	msg := string(ret[:n])
	code, _ = strconv.Atoi(msg[:3])
	message = msg[4 : len(msg)-2]
	ftp.debugInfo("<*cmd*> " + ftp.cmd)
	ftp.debugInfo(fmt.Sprintf("<*code*> %d", code))
	ftp.debugInfo("<*message*> " + message)
	return
}

func (ftp *FTP) Request(cmd string) {
	if ftp.conn == nil { return }
    ftp.conn.Write([]byte(cmd + "\r\n"))
	ftp.cmd = cmd
	ftp.Code, ftp.Message = ftp.Response()
	if cmd == "PASV" {
		start, end := strings.Index(ftp.Message, "("), strings.Index(ftp.Message, ")")
		s := strings.Split(ftp.Message[start:end], ",")
		l1, _ := strconv.Atoi(s[len(s)-2])
		l2, _ := strconv.Atoi(s[len(s)-1])
		ftp.pasv = l1*256 + l2
	}
}

func (ftp *FTP) Pasv() {
	ftp.Request("PASV")
}

func (ftp *FTP) Pwd() {
	ftp.Request("PWD")
}

func (ftp *FTP) Cwd(path string) {
	ftp.Request("CWD " + path)
}

func (ftp *FTP) Mkd(path string) {
	ftp.Request("MKD " + path)
}

func (ftp *FTP) Size(path string) (size int) {
	ftp.Request("SIZE " + path)
	size, _ = strconv.Atoi(ftp.Message)
	return
}

func (ftp *FTP) List() {
	ftp.Pasv()
	ftp.Request("LIST")
}

func (ftp *FTP) Stor(file string, data []byte) {
	ftp.Pasv()
	if data != nil {
		ftp.stream = data
	}
	ftp.Request("STOR " + file)
}

func (ftp *FTP) Quit() {
	ftp.Request("QUIT")
	ftp.conn.Close()
}
