export type userData = {
      Username: string;
      Password: string;
      Email?: string;
      C_password?: string;
}

export async function userFormSubmit(method: string, formData: userData): Promise<boolean> {
      const fetchParams = {
            method: 'POST',
            headers: {
                  'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
      };

      try {
            const response = await fetch('/api/users', fetchParams);
            const data = await response.json();
            if (data.success) {
                  return true
            }
      } catch (error) {
            console.error('There was an error creating user: ' + error);
            return false;
      }
      return false
}