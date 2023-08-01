import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import device_login, handle_response

def mount_volume(Serv_url, volumeID, deviceID, passwd):
    Serv_url=Serv_url+volumeID
    response = requests.post(Serv_url, auth=(deviceID, passwd)) 
    return response

if __name__ == "__main__":
    Serverurl = config['mounts'][0]

    dev_id, dev_pw = device_login()    
    print("볼륨 마운트를 위한 정보를 받습니다.")
    volumeid=input("볼륨 id를 입력하세요: ")
    response = mount_volume(Serverurl, volumeid, dev_id, dev_pw)

    handle_response(response)
