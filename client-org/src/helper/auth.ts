
export type orgData = {
      Name: string;
      Password: string;
      Email?: string;
      C_password?: string;
}

export type Organization = {
      Id?: string;
      Name: string;
      Email: string;
      Profile_picture?: string;
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
            // if (response.status != 200) {
            //       throw new Error("could not create new user")
            // }
            // TODO log user in and store token in local storage? (maybe this is done on backend -> respond with token so localstorage.setitem(data.token))
            if (data.token && data.refresh_token) {
                  localStorage.setItem('token', data.token)
                  localStorage.setItem("refresh_token", data.refresh_token)
                  return true
            } else {
                  throw new Error("could not find token in response")
            }
      } catch (error) {
            throw new Error("could not create new user, " + error)
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

export async function LogoutOrganization(): Promise<number> {
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
                  return 400
            }

            localStorage.removeItem("token")
            localStorage.removeItem("refresh_token")
            location.reload()
            return 204
      } catch (error) {
            console.error(`error trying to sign you out ${error}`)
            return 400
      }
}

export async function GetTokens(): Promise<{ statusCode: number }> {
      try {
            const fetchParams = {
                  method: 'POST',
                  headers: {
                        'Content-Type': 'application/json',
                        Authorization: `Bearer ${localStorage.getItem("refresh_token")}`
                  }
            };

            const response = await fetch('/api/refresh', fetchParams);
            const data = await response.json();
            if (data.token) {
                  localStorage.setItem('token', data.token);
                  return { statusCode: 200 }
            }
            return { statusCode: 400 }
      } catch (error) {
            console.error(`error fetching tokens ${error}`);
            return { statusCode: 400 };
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

export async function UpdateOrganization(orgData: Organization): Promise<number> {
      try {
            const fetchParams = {
                  method: "PUT",
                  headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${localStorage.getItem("token")}`
                  },
                  body: JSON.stringify(orgData)
            }

            const response = await fetch("/api/organizations/update", fetchParams)
            const data = await response.json()
            if (!response.ok) {
                  console.error("response was not ok")
                  return 400
            }
            console.log(data)
            return 200
      } catch (error) {
            console.error(`error trying to update organization ${error}`)
            return 400
      }
}