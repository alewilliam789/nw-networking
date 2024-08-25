import socket
import sys


def get_addr(req_address):
    split_arr = []
    split_addr = ""
    req_path = "/"

    
    if("http://" in req_address):
        split_arr = req_address.split("http://")
        split_addr = split_arr[-1]
    elif("www." in req_address):
        split_arr = req_address.split("www.")
        split_addr = split_arr[-1]
    else:
        split_addr =req_address

    if(len(split_arr) > 1):
        split_addr = split_arr[1]

    if("/" in split_addr):
        if(split_addr[-1] == "/"):
            leng = len(split_addr)
            split_addr = split_addr[:leng-1]
        else:
            split_arr = split_addr.split("/")
            split_addr = split_arr[0]
            req_path = "/" + split_arr[1]

    ip_address = ""

    try:
        ip_address = socket.gethostbyname(split_addr)
    except socket.gaierror as e:
        if (not(e.errno == -2)):
            print(f"{e.strerror}")
        else:
            return "0.0.0.0", split_addr, req_path
    
    return ip_address, split_addr, req_path



def get_mess(sock, ip_addr, req_string):

    sock.connect((ip_addr, 80))

    sock.send(req_string.encode('UTF-8'))

    msg = ""

    while True:
        try:
            part_msg = sock.recv(1024).decode()

            if(len(part_msg)  != 0):
                msg += part_msg
            else:
                return msg

        except Exception as e:
            print(f"{e}")




def make_request(req_address, redirects = 0):

    is_redirect = False
    is_text = False

    sock = socket.socket()

    if("https" in req_address):
        print("This client can only reach http")
        sys.exit(1)
    elif(redirects == 10):
        print("Too many redirects on current path")
        sys.exit(1)
    elif(not("http://" in req_address)):
        print("This url doesn't contain http://")
        sys.exit(1)


    ip_addr, strip_addr, req_path = get_addr(req_address)


    req_string = f"GET {req_path} HTTP/1.1\r\nHost: {strip_addr}\r\nConnection: close\r\n\r\n"
    
    msg = get_mess(sock, ip_addr, req_string)

    sock.close()

    http_header = msg.split("\r\n")

    resp_body = ""

    for line in http_header:
        if("HTTP/1.1" in line):
            if(not(line.find("30") == -1)):
                is_redirect = True
            elif(not(line.find("40") == -1)):
                status_code_arr = line.split(" ")
                print(f"Status Code: {status_code_arr[1]}\nResponse Body: {resp_body}\n")
                sys.exit(1)
        elif("Location:" in line and is_redirect):
            redirects += 1
            location_array = line.split(" ")
            print(f"Redirected to: {location_array[1]}")
            make_request(location_array[1],redirects)
        elif("Content-Type:" in line and not(is_redirect)):
            if("text/html" in line):
                is_text = True
        elif("html" in line or "HTML" in line):
            resp_body = line
    

    if(not(is_redirect) and is_text):
        print(resp_body)


web_address = sys.argv[1]
    
make_request(web_address)

sys.exit(0)
