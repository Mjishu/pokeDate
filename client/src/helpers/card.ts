
export async function cardResponse(liked: boolean, id: string) {
      const fetchParams = {
            method: "POST",
            headers: {
                  "Content-Type": "application/json",
            },
            body: JSON.stringify({ id: id, liked: liked })
      }
      try {
            const response = await fetch(`/api/cards`, fetchParams)
            const data = await response.json()
            return response.status
      } catch (error) {
            console.log("Error trying to get card: " + id + error)
            return
      }
}

export async function getRandomCard() {
      try {
            const fetchParams = {
                  method: "GET",
                  headers: {
                        "Authorization": `Bearer ${localStorage.getItem("token")}`
                  }
            }
            const response = await fetch("/api/cards", fetchParams)
            const data = await response.json()
            return { data, statusCode: response.status }
      } catch (error) {
            console.error(`error trying to get card: ${error}`)
            return { statusCode: 400 }
      }
}