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
                  return true;
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

export async function GetTokens(): Promise<void> {
      console.log("get tokenscalled")
      try {
            const refreshToken = localStorage.getItem('refresh_token');
            const bearerToken = 'Bearer ' + refreshToken;
            if (!refreshToken) {
                  console.log('i dont have a refresh token. log in');
                  return;
            }
            const fetchParams = {
                  method: 'POST',
                  headers: {
                        'Content-Type': 'application/json',
                        "Authorization": bearerToken
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

export async function GetCurrentUser(): Promise<incomingUser | null> {
      await GetTokens()
      try {
            const token = localStorage.getItem('token');
            const bearerToken = 'Bearer ' + token;
            const fetchParams = {
                  method: 'POST',
                  headers: {
                        'Content-Type': 'application/json',
                        "Authorization": bearerToken
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
            let token = localStorage.getItem('token');
            let bearerToken = 'Bearer ' + token;
            const fetchParams = {
                  method: 'PUT',
                  headers: {
                        'Content-Type': 'application/json',
                        "Authorization": bearerToken
                  },
                  body: JSON.stringify(userBody)
            };
            const response = await fetch('/api/users/current', fetchParams);
            if (!response.ok) {
                  const data = await response.json();
                  alert(`error: ${data.message}`)
            }
            return response.status
      } catch (error: any) {
            console.error(`error trying to update user ${error}`);
            return error.status;
      }
}
