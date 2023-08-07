import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import device_login, handle_response, server_print, client_input

def unmount_volume(Serv_url, VolumeID, DeviceID, Passwd):
    Serv_url=Serv_url+VolumeID
    response = requests.delete(Serv_url, auth=(DeviceID, Passwd))
    return response

if __name__ == "__main__":
    dev_id, dev_pw = device_login()    
    Serverurl = config['mounts'][0]

    server_print("볼륨 마운트 해제를 위한 정보를 받습니다.")
    volumeid=client_input("볼륨 id를 입력하세요: ")
    response = unmount_volume(Serverurl, volumeid, dev_id, dev_pw)

    if response.status_code == 200: # worked properly
        print("볼륨이 정상적으로 마운트 해제되었습니다.")
    else:                           # error occurred
        handle_response(response)
