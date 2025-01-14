export type userData = {
      Username: string;
      Password: string;
      Email?: string;
      C_password?: string;
}

export type incomingUser = {
      Id: string;
      Username: string;
      Email: string;
      Profile_picture: string;
      Date_of_birth: string;
}

export async function CreateUser(formData: userData): Promise<boolean> {
      if (formData.Password != formData.C_password) {
            alert("passwords do not match")
            return false
      }
      const fetchParams = {
            method: "POST",
            headers: {
                  'Content-Type': 'application/json'
            },
            body: JSON.stringify({ Username: formData.Username, Email: formData.Email, Password: formData.Password })
      };

      try {
            const response = await fetch('/api/users/create', fetchParams);
            const data = await response.json();
            if (data.status != 200) {
                  return false
            }
            if (data.token && data.refresh_token) {
                  localStorage.setItem("token", data.token)
                  localStorage.setItem("refresh_token", data.refresh_token)
                  return true
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

export async function GetTokens(): Promise<{statusCode: number}> {
      const refresh_token = localStorage.getItem('token')
      if (refresh_token ==null) return {statusCode: 400} 
      console.log("get tokenscalled")
      try {
            const fetchParams = {
                  method: 'POST',
                  headers: {
                        'Content-Type': 'application/json',
                        "Authorization": `Bearer ${refresh_token}`
                  }
            };

            const response = await fetch('/api/refresh', fetchParams);
            const data = await response.json();
            if (data.token) {
                  localStorage.setItem('token', data.token);
            }
            return {statusCode: 200}
      } catch (error) {
            console.error(`error fetching tokens ${error}`);
            return {statusCode : 400};
      }
}

export async function GetCurrentUser(): Promise<incomingUser | null> {
      const refreshStatus = await GetTokens()
      if (refreshStatus.statusCode == 400) return null
      try {
            const fetchParams = {
                  method: 'POST',
                  headers: {
                        'Content-Type': 'application/json',
                        "Authorization": `Bearer ${localStorage.getItem('token')}`
                  }
            };
            const response = await fetch('/api/users/current', fetchParams);
            const data = await response.json();
            if (response.status == 200) {
                  return data;
            } else {
                  return null;
            }
      } catch (err) {
            console.error(`error fetching curernt user data ${err}`);
            return null;
      }
}

type updatedUser = {
      Username: string;
      Email: string;
      Date_of_birth: string;
      Id: string;
}

export async function UpdateUser(userBody: updatedUser): Promise<number> {
      try {
            const fetchParams = {
                  method: 'PUT',
                  headers: {
                        'Content-Type': 'application/json',
                        "Authorization": `Bearer ${localStorage.getItem("token")}`
                  },
                  body: JSON.stringify(userBody)
            };
            const response = await fetch('/api/users/current', fetchParams);
            if (!response.ok) {
                  const data = await response.json();
                  alert(`error: ${data.message}`)
            }
            return response.status
      } catch (error: unknown) {
            console.error(`error trying to update user ${error}`);
            return 400;
      }
}
