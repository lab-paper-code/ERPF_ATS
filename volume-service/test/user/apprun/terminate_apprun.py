import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import handle_response, device_login, server_print, client_input

def terminate_apprun(server_url, apprunID, id, passwd):
    server_url=server_url+apprunID
    response = requests.delete(server_url, auth=(id, passwd))
    return response

if __name__ == "__main__":
    serverurl = config['appruns'][1]
    dev_id, dev_pw = device_login()
    
    apprun_id=client_input("종료할 Apprun ID: ")
    response = terminate_apprun(serverurl, apprun_id, dev_id, dev_pw)
    
    if response.status_code == 200:         # worked properly
        server_print("앱 실행을 종료하였습니다.")
    else:
        handle_response(response)