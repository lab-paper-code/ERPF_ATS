import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import device_login, handle_response, server_print, client_input

def execute_apprun(Serv_url, appID, volumeID, deviceID, password):
    Serv_url=Serv_url+appID
    data = {
        'device_id': deviceID,
        'volume_id': volumeID
    }
    response = requests.post(Serv_url, json=data, auth=(deviceID, password))
    return response

if __name__ == "__main__":
    dev_id, dev_pw = device_login()
    Serverurl = config['appruns'][1]
    
    server_print("앱 실행을 위한 정보를 받습니다.")
    app_id = client_input("사용할 앱 id를 입력하세요: ")
    vol_id=client_input("사용할 볼륨 id를 입력하세요: ")
    response = execute_apprun(Serverurl, app_id, vol_id, dev_id, dev_pw)

    if response.status_code == 200:             # worked properly
        server_print("앱이 정상적으로 실행되었습니다.")
    else:                                       # error occured
        handle_response(response)
