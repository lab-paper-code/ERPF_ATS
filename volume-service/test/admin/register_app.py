import requests, json, sys, array

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path']) # path for utils.py

from utils import admin_login, handle_response

def register_app(Serv_url, Name, RequireGPU, Description, DockerImage, Arguments, OpenPorts, ID, PASSWD): #function sending post request
    data = {
        'name': Name,
        'require_gpu': RequireGPU,
        'description': Description,
        'docker_image': DockerImage,
        'arguments': Arguments,
        'open_ports': OpenPorts # change it later to open_ports
    }
    response = requests.post(Serv_url, json=data, auth=(ID, PASSWD)) # post request
    return response

if __name__ == "__main__": 
    id, password=admin_login() # get ID,PW
    serverurl = config['apps'][0] 
    requiregpu=False
    print("앱 등록을 위한 정보를 받습니다.")
    name=input("앱 이름을 입력하세요: ")
    requiregpu=bool(input("GPU 요청여부를 입력하세요(사용:True, 사용하지 않음:False).\n입력하지 않으면 False으로 처리됩니다: "))
    if requiregpu!=True:
        requiregpu=False
    description=input("Description을 입력하세요. 입력하지 않으면 공백으로 처리됩니다: ")
    dockerimage=input("사용할 도커 이미지를 입력하세요: ")
    arguments=input("실행에 필요한 인자(Arguments)를 입력하세요. 입력하지 않으면 인자를 받지 않는 것으로 처리됩니다: ")
    openports_str=input("앱 실행 시 사용할 포트를 입력하세요(comma로 구분): ")
    if len(openports_str) >= 2:
        openports_list=openports_str.split(',')
        openports=[int(port) for port in openports_list]
        
    response = register_app(serverurl, name, requiregpu, description, dockerimage, arguments, openports, id, password)
    
    handle_response(response)
