#!/bin/bash
python3 "/home/palisade1/Desktop/ksv/volume-service/test/admin\
/register_device.py"
python3 "/home/palisade1/Desktop/ksv/volume-service/test/user\
/volume/create_volume.py"
python3 "/home/palisade1/Desktop/ksv/volume-service/test/user\
/volume/list_volumes.py"
python3 "/home/palisade1/Desktop/ksv/volume-service/test/user\
/app/list_apps.py"
python3 "/home/palisade1/Desktop/ksv/volume-service/test/user\
/apprun/execute_apprun.py"
python3 "/home/palisade1/Desktop/ksv/volume-service/test/user\
/volume/mount_volume.py"
