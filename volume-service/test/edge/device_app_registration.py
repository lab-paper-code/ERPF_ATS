''' 
목표: 앱 등록 # AppRun은 별도일 것
구현할 함수: 앱 등록(app_register_request), 상태 출력(handle_response)
input: 
output: 
'''
import requests

def getLoginInfo(): # function getting ID, PASSWORD for sending request, combine with AppRegisterRequest? -> Q. Login everytime we post?
    print("앱 등록을 위해서는 ID, PASSWORD가 필요합니다.")
    id=input("ID: ")
    passwd=input("PASSWORD: ")
    return id, passwd

def AppRegisterRequest(Serv_url, Name, RequireGPU, Description, DockerImage, Arguments, OpenPorts, id, passwd): # function requesting registration
    data = {
        'name': Name,
        'require_gpu': RequireGPU,
        'description': Description,
        'docker_image': DockerImage,
        'arguments': Arguments,
        'OpenPorts': OpenPorts # change it later to open_ports
    }
    response = requests.post(Serv_url, json=data, auth=(id,passwd), verify=False) # post request, added auth
    return response

def handleResponse(response): # print status
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


if __name__ == "__main__": 
    ID, PASSWORD=getLoginInfo() # get LoginInfo
    
    ''' basic settings for registration'''
    print("앱 등록을 위한 정보를 받습니다.")
    Serverurl = "http://155.230.36.25:31200/apps" # set serverURL, change this to var.
    RequireGPU=False
    Description=""
    Arguments=""
    OpenPorts=[12300]
    Name=input("앱 이름을 입력하세요: ")
    RequireGPU=bool(input("GPU 요청여부를 입력하세요(사용:True, 사용하지 않음:False).\n입력하지 않으면 False으로 처리됩니다: "))
    Description=input("Description을 입력하세요.\n입력하지 않으면 공백으로 처리됩니다: ")
    DockerImage=input("사용할 도커 이미지를 입력하세요: ")
    Arguments=input("실행에 필요한 인자(Arguments)를 입력하세요.\n입력하지 않으면 인자를 받지 않는 것으로 처리됩니다: ")
    OpenPorts[0]=int(input("포트를 입력하세요.\n입력하지 않으면 기본포트(12300)으로 처리됩니다: "))
    response = AppRegisterRequest(Serverurl, Name, RequireGPU, Description, DockerImage, Arguments, OpenPorts, ID, PASSWORD)
    
    handleResponse(response)