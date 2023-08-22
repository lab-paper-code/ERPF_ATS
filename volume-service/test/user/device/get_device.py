import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path'])

from utils import handle_response, device_login, server_print, client_input

def get_device(server_url, deviceID, id, passwd):
    server_url=server_url+deviceID
    response = requests.get(server_url, auth=(id, passwd))
    return response

if __name__ == "__main__":
    dev_id, dev_pw = device_login()
    serverurl = config['devices'][1]
    
    device_id = client_input("반환할 디바이스 id: ")
    response = get_device(serverurl, device_id, dev_id, dev_pw)

    if response.status_code == 200:             # worked properly 
        server_print(f"디바이스 {dev_id}의 정보를 반환합니다.")
        dev_dict=json.loads(response.text)['devices'][0]
        for key in dev_dict:
            if key == "id":
                server_print(f"{'device-id'}: {dev_dict[key]}")
            else:
                server_print(f"{key}: {dev_dict[key]}")
    elif handle_response == 401:                # unauthorized error, print error message from server
        pass           
    else:                                       # other error occurred
        handle_response(response)
