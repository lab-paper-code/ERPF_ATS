from fastapi import APIRouter, Header
from starlette import status

import domain.volume.volume_crud as volume_crud

from domain.volume.volume_schema import AuthInfo, VolumeCreate, VolumeUpdate, Volume

router = APIRouter()

@router.post("/volumes", status_code=status.HTTP_200_OK)
def create_volume(volume_create: VolumeCreate, dev_id: str = Header(..., alias='dev_id'), dev_password: str = Header(..., alias='dev_password')):
    auth_info = AuthInfo(id = dev_id, password = dev_password)
    return volume_crud.create_volume(auth_info, volume_create)

@router.get("/volumes/{volume_id}")
def get_volume(volume_id: str, dev_id: str = Header(..., alias='dev_id'), dev_password: str = Header(..., alias='dev_password')):
    auth_info = AuthInfo(id=dev_id, password=dev_password)
    return volume_crud.get_volume(auth_info, volume_id)

@router.get("/volumes", response_model=list[Volume])
def list_volumes(dev_id: str = Header(..., alias='dev_id'), dev_password: str = Header(..., alias='dev_password')):
    auth_info = AuthInfo(id=dev_id, password=dev_password)
    return volume_crud.list_volumes(auth_info)

@router.patch("/volumes/{volume_id}", status_code=status.HTTP_200_OK)
def update_volumes(volume_id: str, volume_update: VolumeUpdate, dev_id: str = Header(..., alias='dev_id'), dev_password: str = Header(..., alias='dev_password')):
    auth_info = AuthInfo(id=dev_id, password=dev_password)
    return volume_crud.update_volumes(auth_info, volume_update, volume_id)

@router.post("/mounts/{volume_id}", status_code=status.HTTP_204_NO_CONTENT)
def mount_volume(volume_id: str, dev_id: str = Header(..., alias='dev_id'), dev_password: str = Header(..., alias='dev_password')):
    auth_info = AuthInfo(id = dev_id, password = dev_password)
    return volume_crud.mount_volume(auth_info, volume_id)

@router.delete("/mounts/{volume_id}", status_code=status.HTTP_204_NO_CONTENT)
def unmount_volume(volume_id: str, dev_id: str = Header(..., alias='dev_id'), dev_password: str = Header(..., alias='dev_password')):
    auth_info = AuthInfo(id = dev_id, password = dev_password)
    return volume_crud.unmount_volume(auth_info, volume_id)