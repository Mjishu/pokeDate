
export type orgData = {
      Name: string;
      Password: string;
      Email?: string;
      C_password?: string;
}

type Organization = {
      Id?: string;
      Name: string;
      Email: string;
}

export async function CreateOrganization(formdata: orgData): Promise<boolean> {
      const fetchParams = {
            method: "POST",
            headers: {
                  "Content-Type": "application/json"
            },
            body: JSON.stringify(formdata)
      }
      try {
            const response = await fetch("/api/organizations/create", fetchParams) // would like data to have new org ID so i can do the json verify
            const data = await response.json()
            if (response.status != 200) {
                  throw new Error("could not create new user")
            }
            // TODO log user in and store token in local storage? (maybe this is done on backend -> respond with token so localstorage.setitem(data.token))
            if (data.token) {
                  localStorage.setItem('token', data.token)
            } else {
                  throw new Error("could not find token in response")
            }
            return true
      } catch (error) {
            throw new Error("could not create new user")
      }
      return false
}

export async function loginOrganization(formData: orgData): Promise<boolean> {
      const fetchParams = {
            method: "POST",
            headers: {
                  'Content-Type': "application/json"
            },
            body: JSON.stringify({ Name: formData.Name, Password: formData.Password })
      }
      try {
            const response = await fetch("/api/organizations/login", fetchParams)
            const data = await response.json();
            if (response.status == 200) {
                  console.log("status 200")
                  localStorage.setItem("refresh_token", data.refresh_token)
                  localStorage.setItem('token', data.token)
                  return true
            }
      } catch (error) {
            console.error(`error creating user ${error}`)
            return false;
      }
      return false;
}

export async function LogoutOrganization() {
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

export async function GetTokens(): Promise<void> {
      try {
            const fetchParams = {
                  method: 'POST',
                  headers: {
                        'Content-Type': 'application/json',
                        Authorization: `Bearer ${localStorage.getItem("token")}`
                  }
            };

            const response = await fetch('/api/refresh', fetchParams);
            const data = await response.json();
            if (data.token) {
                  localStorage.setItem('token', data.token);
            }
      } catch (error) {
            console.error(`error fetching tokens ${error}`);
            return;
      }
}

export async function GetCurrentOrganization(): Promise<Organization | null> {
      await GetTokens();
      try {
            const fetchParams = {
                  method: "POST",
                  headers: {
                        "Conte-Type": "application/json",
                        Authorization: `Bearer ${localStorage.getItem("token")}`
                  },
            }
            const response = await fetch("/api/organizations/current", fetchParams)
            const data = await response.json()
            if (response.status == 200) {
                  return data
            } else {
                  return null;
            }
      } catch (error) {
            throw new Error(`error fetching current Organization Data ${error}`)
      }
}