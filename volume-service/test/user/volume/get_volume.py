import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import handle_response, device_login, server_print, client_input

def get_volume(server_url, volumeID, deviceID, passwd):
    server_url=server_url+volumeID
    response = requests.get(server_url, auth=(deviceID, passwd))
    return response

if __name__=="__main__":
    dev_id, dev_pw = device_login()
    serverurl = config['volumes'][1] 
    vol_id=client_input("반환할 볼륨 id를 입력하세요: ")
    response = get_volume(serverurl, vol_id, dev_id, dev_pw)

    if response.status_code == 200:         # worked properly
        server_print(f"볼륨 {vol_id}의 정보를 반환합니다.\n")
        vol_dict=json.loads(response.text)
        for key in vol_dict:
            if key == "id":
                server_print(f"{'volume-id'}: {vol_dict[key]}")
            else: 
                server_print(f"{key}: {vol_dict[key]}")

    else:                                   # error occured 
        handle_response(response)