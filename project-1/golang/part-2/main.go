package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
  "net/http"
)


type Path struct {
  RequestPath string
  FilePath string
  ContentType string
}


func create_server(port string) net.Listener {

  if port != "" {
    port_number, err := strconv.Atoi(port)

    if err != nil {
      log.Fatalf("The passed in socket couldn't be converted to a number")
    }

    if !(port_number >= 1024) {
      log.Fatalf("This port isn't over")
    }
  } else {
    port = "1024"
  }

  fmt.Printf("Listening on port %s \n", port)

  full_port := fmt.Sprintf(":%s",port)
  
  ln, err := net.Listen("tcp",full_port)


  if err != nil {
    log.Fatalf("A server couldn't listen on this port")
  }

  return ln
}

func get_request(conn net.Conn) string {
  
  var msg strings.Builder

  for {
    buf := make([]byte, 1024)

    len_, err := conn.Read(buf)

    if err != nil {
      log.Fatalf("Couldn't read from the connection")
    }

    msg.WriteString(string(buf))

    if len_ == 0  || len_ < 1024 {
      break
    }

    buf = make([]byte, 1024)
  }

  return msg.String()
}

func parse_request(request string) map[string] string {
  request_arr := strings.Split(request,"\r\n")

  request_map := make(map[string] string)
  
  for index, line := range request_arr {
    if index == 0 {
      first_arr := strings.Split(line, " ")

      request_map["Method"] = first_arr[0]
      request_map["Path"] = first_arr[1]
    } else {
      line_arr := strings.Split(line, ": ")


      if len(line_arr) > 1 {
        request_map[line_arr[0]] = line_arr[1]
      }
    }
  }

  return request_map
}

func get_path(request_map map[string] string, paths []Path) Path {
  
  request_path := request_map["Path"]

  for _, path := range paths {
    if request_path == path.RequestPath {
      return path
    }
  }

  return Path{ContentType: "text/html; charset=utf-8", }
}

func send_file(path Path, header bytes.Buffer, conn net.Conn) {
  
  conn.Write(header.Bytes())
  
  if path.FilePath != "" {
    file, err := os.Open(path.FilePath)

    if err != nil {
      log.Fatalf("Either an incorrect file path was provided or this file doesn't exist")
    }

    _, write_err := io.Copy(conn, file)

    if write_err != nil {
      if write_err != io.EOF {
        log.Fatalf("There was an issue reading the file to the TCP connection")
      }
    }
  }
}

func send_response(path Path, conn net.Conn) { 
  var hdr bytes.Buffer


  if path.RequestPath == "" {
    hdr.Write([]byte("HTTP/1.1 404 Not Found\r\n"))
  } else {
    hdr.Write([]byte("HTTP/1.1 200 OK\r\n"))
    hdr.Write([]byte("Connection: Keep-Alive\r\n"))
  }
  
  hdr.Write([]byte(fmt.Sprintf("Content-Type: %s\r\n", path.ContentType)))
  
  
  now := time.Now()

  s := now.UTC().Format(http.TimeFormat)

  hdr.Write([]byte(fmt.Sprintf("Date: %s\r\n\r\n", s)))

  
  send_file(path, hdr, conn)

  conn.Close()
}


func main() {
  paths := []Path {{RequestPath: "/rfc2616", FilePath: "./rfc2616.html", ContentType: "text/html"}}
  
  port := ""

  if len(os.Args) > 1 {
    port = os.Args[1]
  }

  ln := create_server(port)

  defer ln.Close()

  conn_chan := make(chan net.Conn)

  go func() {
    for {
      conn, err := ln.Accept()
  
      if err != nil {
        log.Fatalf("A connection couldn't be made to the tcp socket")
      }

      conn_chan <- conn
    }
  }()

  for {
    select {
      case conn := <- conn_chan:
        fmt.Printf("Accepting New Connection\n")
        go func(conn net.Conn){
          request := get_request(conn)

          request_map := parse_request(request) 

          file_path := get_path(request_map, paths)

          send_response(file_path, conn)
        }(conn)
    }
  }
}






