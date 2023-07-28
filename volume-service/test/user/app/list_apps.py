import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path']) # path for utils.py

from utils import admin_login, handle_response

def list_apps(Serv_url, id, passwd):
    response = requests.get(Serv_url, auth=(id,passwd)) # post request
    return response

if __name__ == "__main__":
    ID, PASSWORD=admin_login() # admin_login
    serverurl = config['apps'][0]
    response = list_apps(serverurl, ID, PASSWORD)
    handle_response(response)
    
    print("앱 목록을 반환합니다.")
    apps_list=json.loads(response.text)['apps']
    for app_dict in apps_list:
        for key in app_dict:
            print(f"{key}: {app_dict[key]}")
        print()
    
