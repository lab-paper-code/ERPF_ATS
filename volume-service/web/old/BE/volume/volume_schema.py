from pydantic import BaseModel
from typing import Optional

class AuthInfo(BaseModel):          # info used for authenticating requests
    id : str
    password : str

class VolumeCreate(BaseModel):      # used for creating volumes
    device_id : Optional[str | None] = None                 
    volume_size : Optional[int | str | None] = None

class VolumeUpdate(BaseModel):
    volume_id : Optional[str | None] = None
    mod_volume_size : Optional[int | str | None] = None

class Volume(BaseModel):            # return type of get/list volumes
    id : str
    device_id : str
    volume_size : int
    Mounted : bool | None = None
    created_at : str
    updated_at : str