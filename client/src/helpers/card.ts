
export async function cardResponse(liked: boolean, id: string) {
      console.log("Card was " + (liked === true ? "liked" : "disliked"))
      const fetchParams = {
            method: "POST",
            headers: {
                  "Content-Type": "application/json",
            },
            body: JSON.stringify({ id, liked })
      }
      // const response = await fetch(`/api/cards/${id}`,fetchParams)
      // const data = await response.json()
      // console.log(data)
}