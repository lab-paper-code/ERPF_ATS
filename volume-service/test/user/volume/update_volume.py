import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import device_login, handle_response, server_print, client_input

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

    server_print("볼륨 업데이트를 위한 정보를 받습니다.")
    vol_id=client_input("변경할 볼륨 id를 입력하세요: ")
    volumesize=client_input("변경할 볼륨 크기를 입력하세요: ")
    response = update_volume(Serverurl, vol_id, volumesize, dev_id, dev_pw)
    
    if response.status_code == 200:         # worked properly
        server_print("볼륨 정보가 정상적으로 변경되었습니다.")
    elif response.status_code == 401:       # unauthorized error, print error message from server
        pass                                
    else:                                   # error occurred       
        handle_response(response) 