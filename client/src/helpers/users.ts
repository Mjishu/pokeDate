export type userData = {
      Username: string;
      Password: string;
      Email?: string;
      C_password?: string;
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
                  localStorage.setItem("token", data.token)
            }
      } catch (error) {
            console.error('There was an error creating user: ' + error);
            return false;
      }
      return false
}