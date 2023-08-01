import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import device_login, handle_response

def list_apps(Serv_url, id, passwd):
    response = requests.get(Serv_url, auth=(id,passwd))
    return response

if __name__ == "__main__":
    ID, PASSWORD=device_login()
    serverurl = config['apps'][0]
    response = list_apps(serverurl, ID, PASSWORD)
    handle_response(response)
    
    if response.status_code == 200:
        print("앱 목록을 반환합니다.")
        apps_list=json.loads(response.text)['apps']
        for app_dict in apps_list:
            for key in app_dict:
                print(f"{key}: {app_dict[key]}")
            print()