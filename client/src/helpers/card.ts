
export async function cardResponse(liked: boolean, id: string) {
      const fetchParams = {
            method: "POST",
            headers: {
                  "Content-Type": "application/json",
            },
            body: JSON.stringify({ id: id, liked: liked })
      }
      try {
            const response = await fetch(`http://localhost:8080/cards`, fetchParams)
            const data = await response.json()
            console.log(data)
            return data
      } catch (error) {
            console.log("Error trying to get card: " + id + error)
            return
      }
}

export async function getRandomCard() {
      try {
            const response = await fetch("http://localhost:8080/cards")
            const data = await response.json()
            console.log(data)
            return data
      } catch (error) {
            console.error(`error trying to get card: ${error}`)
            return
      }
}