import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import admin_login, handle_response, client_input, server_print

def register_app(Serv_url, Name, RequireGPU, Description, DockerImage, Arguments, OpenPorts, ID, Passwd):
    data = {
        'name': Name,
        'require_gpu': RequireGPU,
        'description': Description,
        'docker_image': DockerImage,
        'arguments': Arguments,
        'open_ports': OpenPorts
    }
    response = requests.post(Serv_url, json=data, auth=(ID, Passwd))
    return response

if __name__ == "__main__": 
    id, password=admin_login()
    serverurl = config['apps'][0] 
    require_gpu=False

    server_print("앱 등록을 위한 정보를 받습니다.")
    name=client_input("앱 이름을 입력하세요: ")
    require_gpu=client_input("GPU 요청여부를 입력하세요(사용:True 사용 X: False).\n입력하지 않으면 False으로 처리됩니다: ")
    require_gpu=True if (require_gpu=="True" or require_gpu=="true") else False 
    description=client_input("Description을 입력하세요. 입력하지 않으면 공백으로 처리됩니다: ")
    dockerimage=client_input("사용할 도커 이미지를 입력하세요: ")
    arguments=client_input("실행에 필요한 인자(Arguments)를 입력하세요. 입력하지 않으면 인자를 받지 않는 것으로 처리됩니다: ")
    openports_str=client_input("앱 실행 시 사용할 포트를 입력하세요(comma로 구분): ")
    if len(openports_str) >= 2:
        openports_list=openports_str.split(',')
        openports=[int(port) for port in openports_list]
    else: 
        openports=list(openports_str)

    response = register_app(serverurl, name, require_gpu, description, dockerimage, arguments, openports, id, password)

    if response.status_code == 200:                 # worked properly
        server_print("앱이 정상적으로 등록되었습니다. ")
    elif response.status_code == 401:               # unauthorized error occured
        server_print("앱 등록은 admin 계정으로만 가능합니다.")
    else:                                           # todo: need to change response format for each obj    
        handle_response(response)               
