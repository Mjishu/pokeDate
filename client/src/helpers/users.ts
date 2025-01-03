export type userData = {
      Username: string;
      Password: string;
      Email?: string;
      C_password?: string;
}

export type incomingUser = {
      Username: string;

}

export async function userFormSubmit(url: string, method: string, formData: userData): Promise<boolean> {
      const fetchParams = {
            method: "POST",
            headers: {
                  'Content-Type': 'application/json'
            },
            body: JSON.stringify({ formData, exp_seconds: 3600 })
      };

      try {
            const response = await fetch('/api/users' + url, fetchParams);
            const data = await response.json();
            if (data.status == 200) {
                  console.log("status 200")
                  localStorage.setItem("token", data.token)
            }
      } catch (error) {
            console.error('There was an error creating user: ' + error);
            return false;
      }
      return false
}

export async function loginUser(formData: userData): Promise<boolean> {
      const fetchParams = {
            method: "POST",
            headers: {
                  'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username: formData.Username, password: formData.Password })
      };

      try {
            const response = await fetch('/api/users/login', fetchParams);
            const data = await response.json();
            if (data.status == 200) {
                  console.log("status 200")
                  localStorage.setItem("refresh_token", data.refresh_token)
                  localStorage.setItem("token", data.token)
            }
      } catch (error) {
            console.error('There was an error creating user: ' + error);
            return false;
      }
      return false
}

export async function GetCurrentUser() {
      try {
            const jwtToken = localStorage.getItem("token")
            const bearerToken = "Bearer " + jwtToken
            const fetchParams = {
                  method: "POST",
                  headers: {
                        "Authorization": bearerToken
                  }
            }
            const response = await fetch("/api/users/current", fetchParams)
            const data = await response.json()
            return data
      } catch (error) {
            console.error(`error getting current user ${error}`)
            return
      }
}


export async function LogoutUser() {
      try {
            const refreshToken = localStorage.getItem("refresh_token")
            const bearerToken = "Bearer " + refreshToken
            const fetchParams = {
                  method: "POST",
                  headers: {
                        "Authorization": bearerToken
                  }
            }
            const response = await fetch("/api/revoke", fetchParams)
            if (response.status != 204) {
                  alert("issue revoking token")
                  return
            }

            localStorage.removeItem("token")
            localStorage.removeItem("refresh_token")
      } catch (error) {
            console.error(`error trying to sign you out ${error}`)
            return
      }
}