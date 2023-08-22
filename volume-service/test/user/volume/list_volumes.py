import requests, json, sys

with open('./test/conf.json','r') as conf:
    config = json.load(conf)
sys.path.append(config['sys_path']) 

from utils import device_login, handle_response, server_print, client_input

def list_volumes(Serv_url, ID, PASSWD):
    response = requests.get(Serv_url, auth=(ID,PASSWD))
    return response

if __name__ == "__main__":
    id, password=device_login() 
    serverurl = config['volumes'][0]
    response = list_volumes(serverurl, id, password)

    if response.status_code == 200:                     # worked properly 
        server_print("볼륨 정보를 반환합니다.")
        volumes_list=json.loads(response.text)['volumes']
        for vol_dict in volumes_list:
            for key in vol_dict:
                if key == "id":
                    server_print(f"{'volume-id'}: {vol_dict[key]}")
                else:
                    server_print(f"{key}: {vol_dict[key]}")
            print()
    else:                                               # error occurrred
        handle_response(response)
