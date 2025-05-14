<script>
    import request from "../../lib/request";
    import {setCookie , getCookie} from "../../lib/auth.js";
    import {clearOutput, formatDate} from "../../lib/common.js"

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

    let device = null;
    let devices = [];
    let deviceID = "";
    let isMulti = false;
    
    const getDevices = async () => {
        try {
            clearOutput(device, devices);

            let url = isMulti ? "/devices" : `/devices/${deviceID}`;
            const response = await request("GET", url, {}, {}, authInfo.id, authInfo.password);

            setCookie("ID", authInfo.id,30);
            setCookie("PW", authInfo.password,30);

            const parsedResponse = typeof response === 'string' ? JSON.parse(response) : response;
            console.log("parsed server response:", parsedResponse);

            if (isMulti) {
                if (response && response.devices) {
                    devices = [...response.devices].map(dev => ({
                        ...dev,
                        created_at: formatDate(dev.created_at),
                        updated_at: formatDate(dev.updated_at)
                    }));
                }
            } else {
                if (response) {
                    device = {
                        ...response,
                        created_at: formatDate(response.created_at),
                        updated_at: formatDate(response.updated_at)
                    };
                }
            }
        } catch (error) {
            alert(`Error: ${error.message}`);
        }
    };
</script>
<div class="input container mx-auto" style="margin-top: 60px">
    <label for="id">ID:</label>
    <input type="text" id="id" style="margin-right: 10pt; margin-bottom: 5pt;" bind:value={authInfo.id} placeholder="Enter ID" />

    <label for="password">Password:</label>
    <input type="password" id="password" bind:value={authInfo.password} placeholder="Enter Password" />

<form on:submit|preventDefault>
    {#if !isMulti}
        <label for="deviceID">DeviceID:</label>
        <input type="text" id="deviceID" style="margin-right: 10pt;" bind:value={deviceID} placeholder="Enter DeviceID" />
    {/if}

    <label>
        <input type="checkbox" bind:checked={isMulti} on:change={clearOutput(device, devices)} /> Multi-device
    </label>

    <button style="margin-left: 10pt;" on:click={getDevices}>Search Device</button>
</form>

{#if isMulti}
    {#if devices.length > 0}
        <div class="device-container" style="margin-top: 20px">
            <strong style="margin-left: 10px;">Device Info:</strong>
            {#each devices as device}
                <div class="device">
                    <p>ID: {device.id}</p>
                    <p>IP: {device.ip}</p>
                    <p>Description: {device.description || "N/A"}</p>
                    <p>Password: {device.password}</p>
                    <p>Created At: {device.created_at}</p>
                    <p>Updated At: {device.updated_at}</p>
                </div>
            {/each}
        </div>
    {/if}
{:else if device}
    <div class="device-container" style="margin-top: 20px">
        <strong style="margin-left: 10px;">Device Info:</strong>
        <div class="device">
            <p>ID: {device.id}</p>
            <p>IP: {device.ip}</p>
            <p>Description: {device.description || "N/A"}</p>
            <p>Password: {device.password}</p>
            <p>Created At: {device.created_at}</p>
            <p>Updated At: {device.updated_at}</p>
        </div>
    </div>
{/if}
</div>

<style>
    .device-container {
        border: 2px solid #323131;
    }

    .device {
        border-bottom: 1px solid #323131;
        padding: 10px;
    }

    .device:last-child {
        border-bottom: none;
    }
</style>
