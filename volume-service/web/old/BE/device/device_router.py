from fastapi import APIRouter, Header
from starlette import status

import domain.device.device_crud as device_crud

from domain.device.device_schema import AuthInfo, DeviceRegister, DeviceUpdate, Device

router = APIRouter(
    prefix = '/devices'
)

@router.post("")
def register_device(device_register: DeviceRegister, admin_id: str = Header(..., alias='admin_id'), admin_password: str = Header(..., alias='admin_password')):
    auth_info = AuthInfo(id = admin_id, password = admin_password)
    return device_crud.register_device(auth_info, device_register)

@router.get("/{device_id}", response_model=Device)
def get_device(device_id: str, dev_id: str = Header(..., alias='dev_id'), dev_password: str = Header(..., alias='dev_password')):
    auth_info = AuthInfo(id=dev_id, password=dev_password)
    return device_crud.get_device(auth_info, device_id)

@router.get("", response_model=list[Device])
def list_devices(dev_id: str = Header(..., alias='dev_id'), dev_password: str = Header(..., alias='dev_password')):
    auth_info = AuthInfo(id=dev_id, password=dev_password)
    return device_crud.list_devices(auth_info)

@router.patch("/{device_id}", status_code=status.HTTP_200_OK)
def update_device(device_id: str, device_update: DeviceUpdate, dev_id: str = Header(..., alias='dev_id'), dev_password: str = Header(..., alias='dev_password')):
    auth_info = AuthInfo(id=dev_id, password=dev_password)
    return device_crud.update_device(auth_info, device_update)