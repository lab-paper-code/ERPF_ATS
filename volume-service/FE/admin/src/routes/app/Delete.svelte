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

  let app_id = ""; // This will hold the ID of the App to delete

  const deleteApp = async () => {
    const url = `/apps/${app_id}`;

    try {
      const response = await request("DELETE", url, {}, {}, authInfo.id, authInfo.password);
      
      setCookie("ID", authInfo.id,30);
      setCookie("PW", authInfo.password,30);
      
      // console log to check server response
      const parsedResponse = typeof response === 'string' ? JSON.parse(response) : response;
      console.log("parsed server response:", parsedResponse);

      alert(`Successfully deleted App with ID: ${app_id}`);
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

    <label for="app_id">App ID:</label>
    <input
      type="text"
      id="app_id"
      style=" margin-right: 10px;"
      bind:value={app_id}
      placeholder="Enter App ID"
    />
    <button on:click={deleteApp}>Delete App</button>
  </form>
</div>
