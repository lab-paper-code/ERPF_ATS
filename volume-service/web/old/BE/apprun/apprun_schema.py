from pydantic import BaseModel, validator

class AuthInfo(BaseModel):          # info used for request authentication
    id : str
    password : str

class AppRunExecute(BaseModel):    # used for apprun execute
    device_id : str
    volume_id : str
    
    @validator('device_id', 'volume_id')
    def not_empty(cls, v):
        if not v or not v.strip():
            raise ValueError('빈 값은 허용되지 않습니다.')
        return v
    

class AppRun(BaseModel):            # return type of get/list devices
    id : str 
    device_id : str
    volume_id : str
    app_id : str
    terminated : bool | None = None
    created_at : str
    updated_at : str