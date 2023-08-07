import getpass

server_prefix = "[Server]"
client_prefix = "[Client]"

def client_input(prompt):                               # client_input is input function with prefix([Client])
    user_input = input(client_prefix + " " + prompt)
    return user_input

def server_print(*args, sep=' ', end='\n'):             # server_print is print function with prefix([Client])
    message = server_prefix + sep.join(map(str, args)) + end
    print(message, end='')

def admin_login():
    id=getpass.getpass("[Client] ID: ")
    passwd=getpass.getpass("[Client] Password: ")
    print()
    return id, passwd

def device_login(): 
    id=input("[Client] Device ID: ")
    passwd=getpass.getpass("[Client] Password: ")
    print()
    return id, passwd

def handle_response(response): # handle error message
    if response.status_code == 400:
        server_print("Bad Request: 잘못된 요청입니다. 로그인 정보 또는 입력이 정확한지 확인하세요.")
    elif response.status_code == 401:
        server_print("Unauthorized: 인증되지 않은 요청입니다.")
    elif response.status_code == 403:
        server_print("Forbidden: 접근이 금지되었습니다.")
    elif response.status_code == 404:
        server_print("Not Found: 요청한 리소스를 찾을 수 없습니다.")
    elif response.status_code == 500:
        server_print("Internal Server Error: 서버 내부 에러가 발생했습니다.")
    else:
        server_print("Unknown Error: 알 수 없는 에러가 발생했습니다.")