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

    let appUpdate = {
        name: "",
        require_gpu: false,
        description: "",
        docker_image: "",
        arguments: "",
        commands: "",
        stateful: false,
        open_ports: "",
    };

    let app_id = "";

    const updateApp = async () => {
        try {
            // Convert the comma-separated string to an array of numbers
            // appUpdate.open_ports = appUpdate.open_ports.split(",").map(Number); 
        
            const url = `/apps/${app_id}`;
            const response = await request("PATCH", url, appUpdate, {}, authInfo.id, authInfo.password);
            
            setCookie("ID", authInfo.id,30);
            setCookie("PW", authInfo.password,30);
            
            console.log("Received from server:", response); // debugging line
            alert(`Successfully updated App with ID: ${app_id}`);
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

        <label for="app_id">AppID:</label>
        <input
            type="text"
            id="id"
            style=" margin-right: 10px;"
            bind:value={app_id}
            placeholder="Enter AppID"
        />

        <label for="name">Name:</label>
        <input
            type="text"
            id="name"
            style=" margin-right: 10px;"
            bind:value={appUpdate.name}
            placeholder="Enter App Name"
        />
        </form>

        <form>
            <label for="description">Description:</label>
        <input
            type="text"
            id="description"
            style="  margin-top: 5px;margin-right: 10px;"
            bind:value={appUpdate.description}
            placeholder="Enter Description"
        />

        <label for="docker_image">Docker Image:</label>
        <input
            type="text"
            id="docker_image"
            style="margin-right: 10px;"
            bind:value={appUpdate.docker_image}
            placeholder="Enter Docker Image"
        />
        <label for="commands">Commands:</label>
        <input
            type="text"
            id="commands"
            style=" margin-right: 10px;"
            bind:value={appUpdate.commands}
            placeholder="Enter Commands"
        />
        <label for="arguments">Arguments:</label>
        <input
            type="text"
            id="arguments"
            style=" margin-right: 10px;"
            bind:value={appUpdate.arguments}
            placeholder="Enter Arguments"
        />
        </form>
        <form>
        <label for="open_ports">Open Ports:</label>
        <input
            type="text"
            id="open_ports"
            style=" margin-top: 5px; margin-right: 10px;"
            bind:value={appUpdate.open_ports}
            placeholder="Enter Open Ports"
        />
        <label for="require_gpu">Require GPU:</label>
        <input
            type="checkbox"
            id="require_gpu"
            style=" margin-right: 10px;"
            bind:checked={appUpdate.require_gpu}
        />
        <label for="stateful">Stateful:</label>
            <input
            type="checkbox"
            id="stateful"
            style=" margin-right: 10px;"
            bind:checked={appUpdate.stateful}
        />
        <button on:click={updateApp}>Update App</button>
    </form>
</div>


