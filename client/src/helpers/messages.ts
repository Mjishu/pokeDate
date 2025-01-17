import { GetTokens } from "./users";
import type { incomingUser } from "./users";

export type Conversation = {
	Id:                string;
	Conversation_name: string;
	Members:           Conversation_member[];
	Messages:          Messages[]
}

export type Messages = {
	Id:              string;
	From_id:         string; //References users
	Conversation_id: string;
	message:    string;
	Sent_datetime:   Date;
      From_user: incomingUser;
};

export type Conversation_member = {
	Member_id:       string;
	Conversation_id: string; // References Conversation
	Joined_datetime: Date;
	Left_datetime:   Date;
      User: incomingUser
}

export async function CurrentUserMessages(): Promise<Conversation[] | null> {
      await GetTokens()
      try {
            const fetchParams = {
                  method: "POST", 
                  headers : {
                        "Content-Type" :"application/json",
                        "Authorization": `Bearer ${localStorage.getItem("token")}`
                  }
            }
            const response = await fetch("/api/messages", fetchParams)
            const data = await response.json()
            console.log(data)
            if (!response.ok) {
                  return null
            }
            return data
      } catch (error) {
            console.error(`error trying to get current user messages ${error}`)
            return null
      }
}

export async function GetMessage(id: string): Promise<Conversation | null> {
      await GetTokens()
      try {
            const fetchParams = {
                  method: "POST", // might need to be post
                  headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${localStorage.getItem("token")}`
                  }
            }
            console.log(`reaching out to /api/messages/${id}`)
            const response = await fetch(`/api/messages/${id}`, fetchParams)
            const data = await response.json()
            if (!response.ok) {
                  console.error("error fetching message")
                  return null
            }
            console.log(data)
            return data
      }
      catch (error) {
            console.error(`error getting message ${error}`)
            return null
      }
}

export async function CreateMessage(): Promise<string | null> {
      await GetTokens()
      try{
            const fetchParams = {
                  method:"POST",
                  headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${localStorage.getItem("token")}`
                  }
            }
            const response = await fetch("/api/messages/create", fetchParams)
            const data = await response.json()
            console.log(data)
            if(!response.ok) {
                  console.error("response not ok")
                  return null
            }
            return data.id
      }catch(error) {
            console.error("error trying to create new message", error)
            return null
      }

}

export async function SendMessage(message: string, conversation_id :string) {
      await GetTokens()
      try {
            const fetchParams = {
                  method:"POST",
                  headers: {
                        "Content-Type":"application/json",
                        "Authorization": `Bearer ${localStorage.getItem("token")}`
                  },
                  body: JSON.stringify({message})
            }
            const response = await fetch(`/api/messages/${conversation_id}/send`, fetchParams)
            const data = await response.json()
            if(!response.ok) {
                  console.log("sending message response NOT OK")
                  return
            }
            console.log(data)
      }catch(error) {
            console.error(`could not send message ${error}`)
      }
}
