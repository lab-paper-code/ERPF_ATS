import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path']) # path for utils.py

from utils import admin_login, handle_response

def list_volumes(Serv_url, ID, PASSWD):
    response = requests.get(Serv_url, auth=(ID,PASSWD)) # post request
    return response

if __name__ == "__main__":
    id, password=admin_login() 
    serverurl = config['volumes'][0]
    response = list_volumes(serverurl, id, password)
    
    if response.status_code != 401:
        handle_response(response)
        print("볼륨 정보를 반환합니다.")
        volumes_list=json.loads(response.text)['volumes']
        for vol_dict in volumes_list:
            for key in vol_dict:
                print(f"{key}: {vol_dict[key]}")
            print()
    else:
        print("인증 오류: 디바이스 전체 정보는 admin 계정으로만 가져올 수 있습니다.")