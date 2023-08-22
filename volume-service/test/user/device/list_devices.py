import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import device_login, handle_response, server_print

def list_devices(Serv_url, ID, PASSWD):
    response = requests.get(Serv_url, auth=(ID,PASSWD))
    return response

if __name__ == "__main__":
    id, password=device_login()
    serverurl = config['devices'][0]
    
    response = list_devices(serverurl, id, password)

    if response.status_code == 200:   # properly worked
        server_print("디바이스 정보를 반환합니다.")
        print()
        devices_list=json.loads(response.text)['devices']
        for dev_dict in devices_list:
            for key in dev_dict:
                if key == "id":
                    server_print(f"{'device-id'}: {dev_dict[key]}")
                else:
                    server_print(f"{key}: {dev_dict[key]}")
            print()
    elif response.status_code == 401:     # unauthorized error occured
        server_print("디바이스 전체 정보는 admin 계정으로만 가져올 수 있습니다.")
    else:                               # other error occured
        handle_response(response)
