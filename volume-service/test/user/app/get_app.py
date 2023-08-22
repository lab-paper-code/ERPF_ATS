import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import handle_response, device_login, server_print, client_input

def get_app(server_url, AppID, id, passwd):
    server_url=server_url+AppID
    response = requests.get(server_url, auth=(id, passwd))
    return response

if __name__ == "__main__":
    dev_id, dev_pw = device_login()    
    serverurl = config['apps'][1]
    
    app_id=client_input("반환할 App ID: ")
    response = get_app(serverurl, app_id, dev_id, dev_pw)
    
    if response.status_code == 200:         # worked properly
        server_print("앱 정보를 반환합니다.")         
        app_dict=json.loads(response.text)  # make app_dict from response
        for key in app_dict:                # print app_dict
            if key == "id":
                server_print(f"{'app-id'}: {app_dict[key]}") 
            else:
                server_print(f"{key}: {app_dict[key]}")

    else:
        handle_response(response)
