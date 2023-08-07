import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import device_login, handle_response, server_print, client_input

def update_volume(Serv_url, deviceIP, devicePasswordMod, deviceID, PASSWD):
    Serv_url=Serv_url+deviceID
    data = {
        'ip': deviceIP,
        'password': devicePasswordMod
    }
    response = requests.patch(Serv_url, json=data, auth=(deviceID, PASSWD))
    return response

if __name__ == "__main__":
    dev_id, dev_pw = device_login()
    Serverurl = config['devices'][1]

    server_print("디바이스 업데이트를 위한 정보를 받습니다.")
    dev_ip=client_input("변경할 디바이스 ip를 입력하세요: ")
    dev_pw_mod=client_input("변경할 디바이스 패스워드를 입력하세요: ")
    response = update_volume(Serverurl, dev_ip, dev_pw_mod, dev_id, dev_pw)
    
    if response.status_code == 200:                     # worked properly
        server_print("디바이스가 정상적으로 업데이트 되었습니다.")
    elif response.status_code == 401:                   # unauthorized error, print error message from server
        pass              
    else:                                               # other error occurred
        handle_response(response) 
        