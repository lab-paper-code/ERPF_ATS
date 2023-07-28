import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path']) # path for utils.py

from utils import device_login, handle_response

def update_volume(Serv_url, deviceIP, devicePasswordMod, deviceID, PASSWD):
    Serv_url=Serv_url+deviceID
    data = {
        'ip': deviceIP,
        'password': devicePasswordMod
    }
    response = requests.patch(Serv_url, json=data, auth=(deviceID, PASSWD)) # post request
    return response

if __name__ == "__main__":
    Serverurl = config['devices'][1]

    print("디바이스 업데이트를 위한 정보를 받습니다.")
    dev_id, dev_pw = device_login()
    dev_ip=input("변경할 디바이스 ip를 입력하세요: ")
    dev_pw_mod=input("변경할 디바이스 패스워드를 입력하세요: ")
    response = update_volume(Serverurl, dev_ip, dev_pw_mod, dev_id, dev_pw)
    
    if response.status_code != 401:
        handle_response(response)
    else:
        print("인증 오류: 다른 디바이스의 정보를 업데이트할 수 없습니다.")