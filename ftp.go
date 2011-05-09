// FTP Client for Google Go language.
// Author: smallfish <smallfish.xy@gmail.com>
// Date  : 2011-05-09

package main

import "fmt"
import "os"
import "net"
import "strconv"
import "strings"

type FTP struct {
	host    string
	port    int
	user    string
	passwd  string
	pasv    int
	cmd     string
	code    int
	message string
	stream  []byte
	conn    net.Conn
	error   os.Error
}

func (ftp *FTP) Connect(host string, port int) {
	ftp.conn, ftp.error = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	// welcome
	ftp.Response()
	ftp.host = host
	ftp.port = port
}

func (ftp *FTP) Login(user, passwd string) {
	ftp.Request("USER " + user)
	ftp.Request("PASS " + passwd)
	ftp.user = user
	ftp.passwd = passwd
}

func (ftp *FTP) Response() string {
	ret := make([]byte, 1024)
	n, _ := ftp.conn.Read(ret)
	return string(ret[:n])
}

func (ftp *FTP) Request(cmd string) {
	ftp.conn.Write([]byte(cmd + "\r\n"))
	msg := ftp.Response()
	ftp.code, _ = strconv.Atoi(msg[:3])
	ftp.message = msg[4 : len(msg)-2]
	if cmd == "PASV" {
		start, end := strings.Index(ftp.message, "("), strings.Index(ftp.message, ")")
		s := strings.Split(ftp.message[start:end], ",", -1)
		l1, _ := strconv.Atoi(s[len(s)-2])
		l2, _ := strconv.Atoi(s[len(s)-1])
		ftp.pasv = l1*256 + l2
	}
	if (cmd != "PASV") && (ftp.pasv > 0) {
		newRequest(ftp.host, ftp.pasv, ftp.stream)
		ftp.pasv = 0
		ftp.stream = nil
		ftp.Response()
	}
}

func (ftp *FTP) Pasv() {
	ftp.Request("PASV")
}

func (ftp *FTP) Pwd() {
	ftp.Request("PWD")
}

func (ftp *FTP) List() {
	ftp.Pasv()
	ftp.Request("LIST")
}

func (ftp *FTP) Stor(file string, data []byte) {
	ftp.Request("CWD /")
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

// new connect to FTP pasv port, return data
func newRequest(host string, port int, b []byte) string {
	conn, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	defer conn.Close()
	if b != nil {
		conn.Write(b)
		return "OK"
	}
	ret := make([]byte, 4096)
	n, _ := conn.Read(ret)
	return string(ret[:n])
}
