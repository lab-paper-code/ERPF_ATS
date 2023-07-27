import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path']) # path for utils.py

from utils import admin_login, handle_response # 유저는 못 하는지 확인 

def list_devices(Serv_url, ID, PASSWD):
    response = requests.get(Serv_url, auth=(ID,PASSWD)) # post request
    return response

if __name__ == "__main__":
    id, password=admin_login() # admin_login
    serverurl = config['devices'][0]
    response = list_devices(serverurl, id, password)
    handle_response(response)
    if response.status_code == 401:
        print("디바이스 전체 정보는 admin 계정으로만 가져올 수 있습니다.")
    else:
        print("디바이스 정보를 반환합니다.")
        print()
        devices_list=json.loads(response.text)['devices']
        for dev_dict in devices_list:
            for key in dev_dict:
                print(f"{key}: {dev_dict[key]}")
            print()
        
