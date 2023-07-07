''' 
목표: REST API -> device input 받기
구현할 함수: device_register_request(POST 요청), handle_response(결과 출력)
input: IP, Password, (Description)
output: 정상 등록되면 메시지 출력 / 오류 발생해도 출력
'''

# ref. /volume-service/rest/device_handlers.go: handleRegisterDevice 함수

import requests

def device_register_request(Serv_url, IP, Password, Description):
    data = {
        'IP': IP,
        'Password': Password,
        'Description': Description # default: blank, optional
    }
    response = requests.post(Serv_url, data=data) # post request
    return response
    
Serverurl = '155.230.36.25:31200/devices' # set serverURL, change this to var.
print("디바이스 등록을 위한 정보를 받습니다.")

IP=input("IP를 입력하세요: ")
Password=input("Password를 입력하세요: ")
Description=input("디바이스에 대한 설명을 입력하세요. 입력하지 않으면 공백으로 처리됩니다.\n: ")

response = device_register_request(Serverurl, IP, Password, Description)

def handle_response(response): # print status
    if response.status_code == 200:
        print("디바이스가 정상적으로 등록되었습니다.")
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

handle_response(response)

    