import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path']) # path for utils.py

from utils import handle_response, device_login

def get_device(server_url, deviceID, id, passwd):
    server_url=server_url+deviceID
    response = requests.get(server_url, auth=(id, passwd))
    return response

if __name__ == "__main__":
    serverurl = config['devices'][1]
    dev_id, dev_pw = device_login() # TODO: Notify user created deviceID
    device_id = input("반환할 디바이스 id: ")
    response = get_device(serverurl, device_id, dev_id, dev_pw)

    if response.status_code != 401: # print only if authorized
        handle_response(response)
        print(f"디바이스 {dev_id}의 정보를 반환합니다.")
        dev_dict=json.loads(response.text)['devices'][0] # make str to dict
        for key in dev_dict:
            print(f"{key}: {dev_dict[key]}")
    else:
        print("인증 오류: 다른 디바이스의 정보를 가져올 수 없습니다.") # TODO: check how to get admin access
