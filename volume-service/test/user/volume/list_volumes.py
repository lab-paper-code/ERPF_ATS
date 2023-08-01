import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path']) 

from utils import device_login, handle_response

def list_volumes(Serv_url, ID, PASSWD):
    response = requests.get(Serv_url, auth=(ID,PASSWD))
    return response

if __name__ == "__main__":
    id, password=device_login() 
    serverurl = config['volumes'][0]
    response = list_volumes(serverurl, id, password)
    handle_response(response)

    if response.status_code == 200:
        print("볼륨 정보를 반환합니다.")
        volumes_list=json.loads(response.text)['volumes']
        for vol_dict in volumes_list:
            for key in vol_dict:
                print(f"{key}: {vol_dict[key]}")
            print()