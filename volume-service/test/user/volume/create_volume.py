import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path']) # path for utils.py

from utils import device_login, handle_response

def register_volume(Serv_url, DeviceID, VolumeSize, PASSWD):
    data = {
        'device_id': DeviceID,
        'volume_size': VolumeSize
    }
    response = requests.post(Serv_url, json=data, auth=(DeviceID, PASSWD)) # post request
    return response

if __name__ == "__main__":
    Serverurl = config['volumes'][0]

    print("볼륨 생성을 위한 정보를 받습니다.")
    dev_id, dev_pw = device_login()
    volumesize=input("요청할 볼륨 크기를 입력하세요: ")
    response = register_volume(Serverurl, dev_id, volumesize, dev_pw)

    handle_response(response)
