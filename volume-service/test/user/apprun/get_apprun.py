import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import handle_response, device_login, server_print, client_input

def get_apprun(server_url, apprunID, id, passwd):
    server_url=server_url+apprunID
    response = requests.get(server_url, auth=(id, passwd))
    return response

if __name__ == "__main__":
    dev_id, dev_pw = device_login()
    serverurl = config['appruns'][1]
    
    apprun_id=client_input("반환할 Apprun ID: ")
    response = get_apprun(serverurl, apprun_id, dev_id, dev_pw)
    

    if response.status_code == 200:         # worked properly
        server_print("앱 배포 정보를 반환합니다.")
        apprun_dict=json.loads(response.text)
        for key in apprun_dict:
            if key == "id":
                server_print(f"{'apprun_id'}: {apprun_dict[key]}")
            else:
                server_print(f"{key}: {apprun_dict[key]}")
    else:                                   # error occured
        handle_response(response)
