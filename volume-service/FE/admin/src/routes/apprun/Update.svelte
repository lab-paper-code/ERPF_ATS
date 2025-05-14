<script>
    import request from "../../lib/request";
    import {setCookie , getCookie} from "../../lib/auth.js";

    let authInfo = {
        id: "",
        password: "",
    };

    let authCookieID = getCookie("ID");
    let authCookiePW = getCookie("PW");
    if (authCookieID && authCookiePW){ 
        authInfo.id = authCookieID;
        authInfo.password = authCookiePW;
    }

    let appRunUpdate = {
        app_id: "",
        device_id: "",
        volume_id: "",
    };

    let apprun_id = "";

    const updateAppRun = async () => {
        try {
            const url = `/appruns/${apprun_id}`;
            const response = await request("PATCH", url, appRunUpdate, {}, authInfo.id, authInfo.password);
            
            setCookie("ID", authInfo.id,30);
            setCookie("PW", authInfo.password,30);
            
            console.log("Received from server:", response); // debugging line
            alert(`Successfully updated AppRun with ID: ${apprun_id}`);
        } catch (error){
        alert(`Error: ${error.message}`);
        }
    };
</script>


<div class="input container" style="margin-top: 60px;">
    <form>
        <label for="id">ID:</label>
        <input
            type="password"
            id="id"
            style=" margin-right: 10px;"
            bind:value={authInfo.id}
            placeholder="Enter ID"
        />
    
        <label for="password">Password:</label>
        <input
            type="password"
            id="password"
            style=" margin-right: 10px;"
            bind:value={authInfo.password}
            placeholder="Enter Password"
        />

        <label for="apprun_id">ID:</label>
        <input
            type="text"
            id="apprun_id"
            style=" margin-right: 10px;"
            bind:value={apprun_id}
            placeholder="Enter AppRunID"
        />

        <label for="app_id">AppID:</label>
        <input
            type="text"
            id="app_id"
            style=" margin-right: 10px;"
            bind:value={appRunUpdate.app_id}
            placeholder="Enter AppID"
        />
        </form>

        <form>
            <label for="device_id">DeviceID:</label>
        <input
            type="text"
            id="device_id"
            style="  margin-top: 5px;margin-right: 10px;"
            bind:value={appRunUpdate.device_id}
            placeholder="Enter DeviceID"
        />

        <label for="volume_id">VolumeID:</label>
        <input
            type="text"
            id="volume_id"
            style="margin-right: 10px;"
            bind:value={appRunUpdate.volume_id}
            placeholder="Enter VolumeID"
        />
        <button on:click={updateAppRun}>Update AppRun</button>
    </form>
</div>


