def get_login_info(): # function getting ID, PASSWORD for sending request, combine with AppRegisterRequest? -> Q. Login everytime we post?
    id=input("ID: ")
    passwd=input("PASSWORD: ")
    print()
    return id, passwd

def handle_response(response): # print status
    if response.status_code == 200:
        print("정상적으로 처리되었습니다.")
    elif response.status_code == 400:
        print("Bad Request: 잘못된 요청입니다.")
    elif response.status_code == 401:
        print("Unauthorized: 인증되지 않은 요청입니다.")
    elif response.status_code == 403:
        print("Forbidden: 접근이 금지되었습니다.")
    elif response.status_code == 404:
        print("Not Found: 요청한 리소스를 찾을 수 없습니다.")
    else:
        print("Unknown Error: 알 수 없는 에러가 발생했습니다.")