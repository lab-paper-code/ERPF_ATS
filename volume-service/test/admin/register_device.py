import requests, json, sys, getpass

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import admin_login, handle_response

def register_device(Serv_url, IP, PASSWORD, Description, ID, PASSWD):
    data = {
        'ip': IP,
        'password': PASSWORD,
        'description': Description
    }
    response = requests.post(Serv_url, json=data, auth=(ID,PASSWD))
    return response

if __name__ == "__main__":
    ID, PASSWORD=admin_login()
    Serverurl = config['devices'][0]

    print("디바이스 등록을 위한 정보를 받습니다.")
    ip=input("IP를 입력하세요. 생략 시 요청을 보낸 IP를 입력합니다: ")
    password=getpass.getpass("Password를 입력하세요: ")
    description=input("디바이스에 대한 설명을 입력하세요. 생략 시 공백으로 처리됩니다: ")
    response = register_device(Serverurl, ip, password, description, ID, PASSWORD)
    
    if response.status_code == 401:
        print("디바이스 등록은 admin 계정에서만 가능합니다.")
    else: 
        handle_response(response)
        

    
