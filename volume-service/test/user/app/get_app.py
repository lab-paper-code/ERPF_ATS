import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path']) # path for utils.py

from utils import handle_response, device_login

def get_app(server_url, AppID, id, passwd):
    server_url=server_url+AppID
    response = requests.get(server_url, auth=(id, passwd))
    return response

if __name__ == "__main__":
    serverurl = config['apps'][1]
    dev_id, dev_pw = device_login() # TODO: Notify user created deviceID
    app_id=input("반환할 App ID: ")
    response = get_app(serverurl, app_id, dev_id, dev_pw)
    handle_response(response)

    print("앱 정보를 반환합니다.")
    # print(response.text)
    app_dict=json.loads(response.text) # make str to dict
    for key in app_dict:
        print(f"{key}: {app_dict[key]}")
