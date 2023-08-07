import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import device_login, handle_response, server_print, client_input

def register_volume(Serv_url, DeviceID, VolumeSize, PASSWD):
    data = {
        'device_id': DeviceID,
        'volume_size': VolumeSize
    }
    response = requests.post(Serv_url, json=data, auth=(DeviceID, PASSWD))
    return response

if __name__ == "__main__":
    dev_id, dev_pw = device_login()    
    Serverurl = config['volumes'][0]
    
    server_print("볼륨 생성을 위한 정보를 받습니다.")
    volumesize=client_input("요청할 볼륨 크기를 입력하세요: ")
    response = register_volume(Serverurl, dev_id, volumesize, dev_pw)

    if response.status_code == 200:     # worked properly
        print("볼륨이 정상적으로 생성되었습니다.")
    else:                               # error occurred
        handle_response(response)
