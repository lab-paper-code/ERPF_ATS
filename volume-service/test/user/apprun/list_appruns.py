import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path']) # path for utils.py

from utils import admin_login, handle_response

def list_appruns(Serv_url, id, passwd):
    response = requests.get(Serv_url, auth=(id,passwd)) # post request
    return response

if __name__ == "__main__":
    ID, PASSWORD=admin_login() # admin_login
    serverurl = config['appruns'][0]
    response = list_appruns(serverurl, ID, PASSWORD)
    handle_response(response)
    
    print("앱 실행정보 목록을 반환합니다.\n")
    appruns_list=json.loads(response.text)['app_runs']
    for apprun_dict in appruns_list:
        for key in apprun_dict:
            print(f"{key}: {apprun_dict[key]}")
        print()
