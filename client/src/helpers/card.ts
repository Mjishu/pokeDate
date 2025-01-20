import { GetTokens } from "./users"

export async function cardResponse(liked: boolean, id: string) {
      await GetTokens()
      const fetchParams = {
            method: "POST",
            headers: {
                  "Content-Type": "application/json",
                  "Authorization" : `Bearer ${localStorage.getItem("token")}`
            },
            body: JSON.stringify({ Animal_id: id, Liked: liked })
      }
      try {
            const response = await fetch(`/api/cards`, fetchParams)
            const data = await response.json()
            console.log(data)
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
            if (!response.ok) {
                  return null
            }
            const data = await response.json()
            return { data, statusCode: response.status }
      } catch (error) {
            console.error(`error trying to get card: ${error}`)
            return  null
      }
}