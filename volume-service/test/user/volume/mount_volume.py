import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import device_login, handle_response, server_print, client_input

def mount_volume(Serv_url, VolumeID, MountPath, DeviceID, Passwd):
    Serv_url=Serv_url+VolumeID
    data ={
        'mount_path': MountPath
    }
    response = requests.post(Serv_url, json=data, auth=(DeviceID, Passwd)) 
    return response

if __name__ == "__main__":
    dev_id, dev_pw = device_login()
    Serverurl = config['mounts'][0]
        
    server_print("볼륨 마운트를 위한 정보를 받습니다.")
    volumeid=client_input("볼륨 id를 입력하세요: ")
    mountpath=client_input("마운트 위치를 입력하세요(ex: /uploads): ")
    response = mount_volume(Serverurl, volumeid, mountpath, dev_id, dev_pw)

    if response.status_code == 200: # worked properly
        server_print("볼륨이 정상적으로 마운트되었습니다.")
    else:                           # error occurred
        handle_response(response)
