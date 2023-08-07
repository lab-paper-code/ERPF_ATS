import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import device_login, handle_response

def update_volume(Serv_url, VolumeID, VolumeSize, DeviceID, Passwd):
    Serv_url=Serv_url+VolumeID
    data = {
        'volume_size': VolumeSize,
    }
    response = requests.patch(Serv_url, json=data, auth=(DeviceID, Passwd))
    return response

if __name__ == "__main__":
    dev_id, dev_pw = device_login()
    Serverurl = config['volumes'][1]

    print("볼륨 업데이트를 위한 정보를 받습니다.")
    vol_id=input("변경할 볼륨 id를 입력하세요: ")
    volumesize=input("변경할 볼륨 크기를 입력하세요: ")
    response = update_volume(Serverurl, vol_id, volumesize, dev_id, dev_pw)
    
    if response.status_code == 401: pass           # unauthorized error, print error message from server
    else:                                          
        handle_response(response) 