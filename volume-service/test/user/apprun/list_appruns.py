import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path']) 

from utils import device_login, handle_response

def list_appruns(Serv_url, id, passwd):
    response = requests.get(Serv_url, auth=(id,passwd))
    return response

if __name__ == "__main__":
    ID, PASSWORD=device_login()
    serverurl = config['appruns'][0]
    response = list_appruns(serverurl, ID, PASSWORD)
    handle_response(response)
    
    if response.status_code == 200:         # properly worked
        print("앱 실행정보 목록을 반환합니다.\n")
        appruns_list=json.loads(response.text)['app_runs']
        for apprun_dict in appruns_list:
            for key in apprun_dict:
                print(f"{key}: {apprun_dict[key]}")
            print()
