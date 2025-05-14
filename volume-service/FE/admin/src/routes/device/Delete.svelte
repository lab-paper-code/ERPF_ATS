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

  let device_id = ""; // This will hold the ID of the Device to delete

  const deleteDevice = async () => {
    const url = `/devices/${device_id}`;

    try {
      const response = await request("DELETE", url, {}, {}, authInfo.id, authInfo.password);
      
      setCookie("ID", authInfo.id,30);
      setCookie("PW", authInfo.password,30);

      // console log to check server response
      const parsedResponse = typeof response === 'string' ? JSON.parse(response) : response;
      console.log("parsed server response:", parsedResponse);

      alert(`Successfully deleted Device with ID: ${device_id}`);
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

    <label for="device_id">Device ID:</label>
    <input
      type="text"
      id="device_id"
      style=" margin-right: 10px;"
      bind:value={device_id}
      placeholder="Enter Device ID"
    />
    <button on:click={deleteDevice}>Delete Device</button>
  </form>
</div>
