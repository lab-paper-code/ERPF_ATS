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

    let device_id = "";

    let deviceUpdate = {
        ip: "",
        password: "",
        description: "",
    };

    

    const updateDevice = async () => {
        try {
            const url = `/devices/${device_id}`;
            
            const response = await request("PATCH", url, deviceUpdate, {}, authInfo.id, authInfo.password); // capture response
            
            setCookie("ID", authInfo.id,30);
            setCookie("PW", authInfo.password,30);
            
            console.log("sent to server", updateDevice);
            console.log("Received from server:", response); // debugging line
            alert(`Successfully updated Device with ID: ${device_id}`);
        } catch (error) {
            alert(`Error: ${error.message}`);
        }
    };
</script>

<div class="input container mx-auto" style="margin-top: 60px">
    <form>
        <label for="id">ID:</label>
        <input
            type="text"
            id="id"
            style="margin-bottom: 5px; margin-right: 10px;"
            bind:value={authInfo.id}
            placeholder="Enter ID"
        />

        <label for="password">Password:</label>
        <input
            type="password"
            id="password"
            bind:value={authInfo.password}
            placeholder="Enter Password"
        />

        <label for="id">DeviceID:</label>
        <input
            type="text"
            id="device_id"
            bind:value={device_id}
            placeholder="Enter DeviceID"
        />

        <label for="ip">IP:</label>
        <input
            type="text"
            id="ip"
            style="margin-right: 10px;"
            bind:value={deviceUpdate.ip}
            placeholder="Enter IP to Change"
        />

        <label for="password">Device Password:</label>
        <input
            type="password"
            id="password"
            style="margin-right: 10px;"
            bind:value={deviceUpdate.password}
            placeholder="Enter Device Password to Change"
        />

        <label for="description">Device Description:</label>
        <input
            type="text"
            id="description"
            style="margin-right: 10px;"
            bind:value={deviceUpdate.description}
            placeholder="Enter Device Description to Change"
        />
        <button on:click={updateDevice}>Update Device</button>
    </form>
</div>
