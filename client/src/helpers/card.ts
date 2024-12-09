
export async function cardResponse(liked: boolean, id: string) {
      const fetchParams = {
            method: "POST",
            headers: {
                  "Content-Type": "application/json",
            },
            body: JSON.stringify({ id: id, liked: liked })
      }
      const response = await fetch(`http://localhost:8080/cards`, fetchParams)
      const data = await response.json()
      console.log(data)
      return data
}

export async function getRandomCard() {
      const response = await fetch("http://localhost:8080/cards")
      const data = await response.json()
      console.log(data)
      return data
}