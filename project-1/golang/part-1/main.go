package main

import (
	"fmt"
	"io"
	"net"
	"log"
	"os"
	"strings"
)


func get_addr(req_addr string) (string, string, string) {

  strip_addr := req_addr[7:]

  req_path := "/"

  port := "80"

  path_arr := strings.Split(strip_addr, "/")

  prefix_addr := path_arr[0]

  strip_arr := strings.Split(prefix_addr,":")

  strip_addr = strip_arr[0]

  if len(strip_arr) > 1 {
    port = strip_arr[1]
  }


  
  if len(path_arr) > 1 {
    var sb strings.Builder
    for index, value := range path_arr {
      if index > 0 {
        sb.WriteString("/")
        sb.WriteString(value)
      }
    }

    req_path = sb.String()
  }

  ip_arr, err := net.LookupIP(strip_addr)

  if err != nil {
    log.Fatalf(err.Error())
  }

  ip_addr := ""

  if len(ip_arr) > 1 {
    ip_addr = fmt.Sprintf("%s:%s",ip_arr[1],port)
  } else {
    ip_addr = fmt.Sprintf("%s:%s",ip_arr[0], port)
  }

  return strip_addr, ip_addr, req_path
}

func get_msg(ip_addr string, req_path string, dns_addr string) string {
  
  req_string := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", req_path, dns_addr)

  socket, err := net.Dial("tcp", ip_addr)

  defer socket.Close()

  if err != nil {
    log.Fatalf(err.Error())
  }

  writ, err  := socket.Write([]byte(req_string))

  if err != nil {
    log.Fatalf(err.Error())
  } else if writ != len(req_string) {
    log.Fatalf("Couldn't write entire request string")
  }

  var msg strings.Builder

  part_msg  := make([]byte, 1024) 

  for {
    _, err := socket.Read(part_msg)

    if err != nil {
      if err != io.EOF { 
        log.Fatalf(err.Error())
      } 
      break
    }

    msg.WriteString(string(part_msg))
    part_msg = make([]byte, 1024)
  }

  return  msg.String()
}


func get_response(req_addr string, redirects int) {
  
  if strings.HasPrefix(req_addr, "https") {
    log.Fatalf("This client can only reach http")
  } else if redirects == 10 {
    log.Fatalf("There have been too many redirects")
  }

  if !strings.HasPrefix(req_addr, "http://") {
    log.Fatalf("This domain is missing http://")
  }

  dns_addr, ip_addr, req_path := get_addr(req_addr)

  msg := get_msg(ip_addr, req_path, dns_addr)

  headers := strings.Split(msg, "\r\n")

  segments := strings.Split(msg,"\r\n\r\n")

  req_body := ""

  for _, line := range segments {

    if strings.Contains(line, "<html") || strings.Contains(line, "<HTML") {
      req_body = line
      break;
    }
  }

  is_redirect  := false
  is_html := false

  for _, line := range headers {

    if strings.HasPrefix(line, "HTTP") {
      status_code := line[9:12]

      switch status_code[0:2] {
        case "30":
          is_redirect = true
        case "40":
          log.Fatalf("Status Code: %s", status_code)
      }
    } else if is_redirect && strings.HasPrefix(line, "Location") { 
      location_arr := strings.Split(line, " ")
      
      if len(location_arr) != 1 {
        fmt.Println("Redirected to %s",location_arr[1])
        redirects += 1
        get_response(location_arr[1], redirects) 
      } else {
        log.Fatalf("Redirect address not provided")
      }
    } else if strings.HasPrefix(line, "Content-Type:") && strings.Contains(line, "text/html"){
        is_html = true
    }   
  }

  if is_html && !is_redirect {
    fmt.Println(req_body)
  }
}


func main() {
  
  args := os.Args[1]

  get_response(args,0)
}
