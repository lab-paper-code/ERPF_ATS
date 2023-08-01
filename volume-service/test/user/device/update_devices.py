import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import device_login, handle_response

def update_volume(Serv_url, deviceIP, devicePasswordMod, deviceID, PASSWD):
    Serv_url=Serv_url+deviceID
    data = {
        'ip': deviceIP,
        'password': devicePasswordMod
    }
    response = requests.patch(Serv_url, json=data, auth=(deviceID, PASSWD))
    return response

if __name__ == "__main__":
    Serverurl = config['devices'][1]

    print("디바이스 업데이트를 위한 정보를 받습니다.")
    dev_id, dev_pw = device_login()
    dev_ip=input("변경할 디바이스 ip를 입력하세요: ")
    dev_pw_mod=input("변경할 디바이스 패스워드를 입력하세요: ")
    response = update_volume(Serverurl, dev_ip, dev_pw_mod, dev_id, dev_pw)
    
    if response.status_code == 401: pass           # unauthorized error, print error message from server
    else:                                          
        handle_response(response) 
        