import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import handle_response, device_login

def get_device(server_url, deviceID, id, passwd):
    server_url=server_url+deviceID
    response = requests.get(server_url, auth=(id, passwd))
    return response

if __name__ == "__main__":
    serverurl = config['devices'][1]
    dev_id, dev_pw = device_login()
    device_id = input("반환할 디바이스 id: ")
    response = get_device(serverurl, device_id, dev_id, dev_pw)
    if response.status_code == 200:             # properly worked
        print(f"디바이스 {dev_id}의 정보를 반환합니다.")
        dev_dict=json.loads(response.text)['devices'][0]
        for key in dev_dict:
            print(f"{key}: {dev_dict[key]}")
    elif handle_response == 401: 
        pass           # unauthorized error, print error message from server
    else:                                       # other error occurred
        handle_response(response)
