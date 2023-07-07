''' 
목표: 앱 등록 # AppRun은 별도일 것
구현할 함수: 앱 등록(app_register_request), 상태 출력(handle_response)
input: 
output: 
'''
import requests

def app_register_request(Serv_url, Name,RequireGPU, DockerImage, Description, Arguments, OpenPorts): 
    data = {
        'Name': Name,
        'RequireGPU': RequireGPU,
        'Description': Description,
        'DockerImage': DockerImage,
        'Arguments': Arguments,
        'OpenPorts': OpenPorts
    }
    response = requests.post(Serv_url, data=data) # post request
    return response

Serverurl = '155.230.36.25:31200/apps' # set serverURL, change this to var.
print("앱 등록을 위한 정보를 받습니다.")

RequireGPU=False
Description=""
Arguments=""
OpenPorts="12300"
Name=input("앱 이름을 입력하세요: ")
RequireGPU=input("GPU 요청여부를 입력하세요(사용:True, 사용하지 않음:False). 입력하지 않으면 False으로 처리됩니다.:  ")
Description=input("Description을 입력하세요 입력하지 않으면 공백으로 처리됩니다.: ")
DockerImage=input("사용할 도커 이미지를 입력하세요: ")
Arguments=input("실행에 필요한 인자(Arguments)를 입력하세요. 입력하지 않으면 인자를 받지 않는 것으로 처리됩니다.: ")
OpenPorts=input("포트를 입력하세요. 입력하지 않으면 기본포트(12300)으로 처리됩니다. : ")
response = app_register_request(Serverurl, Name,RequireGPU, DockerImage, Description, Arguments, OpenPorts)


def handle_response(response): # print status
    if response.status_code == 200:
        print("앱이 정상적으로 등록되었습니다.")
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