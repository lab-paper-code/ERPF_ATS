import requests, json, sys, getpass

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import admin_login, handle_response, client_input, server_print, client_prefix

def register_device(Serv_url, IP, PASSWORD, Description, ID, PASSWD):
    data = {
        'ip': IP,
        'password': PASSWORD,
        'description': Description
    }
    response = requests.post(Serv_url, json=data, auth=(ID,PASSWD))
    return response

def list_devices(Serv_url, ID, PASSWD):                           
    response = requests.get(Serv_url, auth=(ID,PASSWD))
    return response

if __name__ == "__main__":
    ID, PASSWORD=admin_login()
    serverurl = config['devices'][0]

    server_print("디바이스 등록을 위한 정보를 받습니다.")
    ip=client_input("IP를 입력하세요. 생략 시 요청을 보낸 IP를 입력합니다: ")
    password=getpass.getpass(client_prefix, "Password를 입력하세요: ")
    description=client_input("디바이스에 대한 설명을 입력하세요. 생략 시 공백으로 처리됩니다: ")
    response = register_device(serverurl, ip, password, description, ID, PASSWORD)
    
    if response.status_code == 200:
        server_print("디바이스가 정상적으로 등록되었습니다.")
                                 # non-admin user request device register
    elif response.status_code == 401:                       # work properly
        server_print("디바이스 등록은 admin 계정에서만 가능합니다.")       
    else:
        handle_response(response)                           # other error occured
        
    response2 = list_devices(serverurl, ID, PASSWORD)       # return registered deviceID

    if  response2.status_code == 200:                       # worked properly
        print()
        server_print("생성된 디바이스 정보를 반환합니다.")
        dev_dict=json.loads(response2.text)['devices'][-1]  # get most recent device
        for key in ['id', 'ip', 'description']:
            if key == "id":
                server_print(f"{'device-id'}: {dev_dict[key]}") # change to normal print if needed
            else:
                server_print(f"{key}: {dev_dict[key]}")
    else:                                                   # already admin -> other error occured
        handle_response(response2)
        exit()