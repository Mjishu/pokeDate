
export type Message = {
      From: string;
      Content: string;
      Date_sent: Date;
}

export type Messages = {
      message: Message[];
}

export async function CurrentUserMessages(): Promise<Messages[] | null> {
      try {
            const fetchParams = {
                  method: "GET" // Might need to be POST
            }
            const response = await fetch("/api/messages", fetchParams)
            const data = await response.json()
            if (!response.ok) {
                  return null
            }
            return data
      } catch (error) {
            console.error(`error trying to get current user messages ${error}`)
            return null
      }
}

export async function GetMessage(id: string): Promise<Messages | null> {
      try {
            const fetchParams = {
                  method: "GET", // might need to be post
                  headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${localStorage.getItem("token")}`
                  }
            }
            const response = await fetch(`/api/messages/${id}`, fetchParams)
            const data = await response.json()
            if (!response.ok) {
                  console.error("error fetching message")
                  return null
            }
            return data
      }
      catch (error) {

            return null
      }
}