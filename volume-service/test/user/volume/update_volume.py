import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path']) # path for utils.py

from utils import device_login, handle_response

def update_volume(Serv_url, volumeID, volumeSize, deviceID, PASSWD):
    Serv_url=Serv_url+volumeID
    data = {
        'volume_size': volumeSize
    }
    response = requests.patch(Serv_url, json=data, auth=(deviceID, PASSWD)) # post request
    return response

if __name__ == "__main__":
    Serverurl = config['volumes'][1]

    print("볼륨 업데이트를 위한 정보를 받습니다.")
    dev_id, dev_pw = device_login()
    vol_id=input("변경할 볼륨 id를 입력하세요: ")
    volumesize=input("변경할 볼륨 크기를 입력하세요: ")
    response = update_volume(Serverurl, vol_id, volumesize, dev_id, dev_pw)
    
    if response.status_code != 401:
        handle_response(response)
    else:
        print("인증 오류: 다른 디바이스의 볼륨 정보를 수정할 수 없습니다.")