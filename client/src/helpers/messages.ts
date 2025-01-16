import { GetTokens } from "./users";

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
	Message_text:    string;
	Sent_datetime:   Date;
};

export type Conversation_member = {
	Member_id:       string;
	Conversation_id: string; // References Conversation
	Joined_datetime: Date;
	Left_datetime:   Date;
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
            const response = await fetch(`/api/messages/${id}`, fetchParams)
            const data = await response.json()
            if (!response.ok) {
                  console.error("error fetching message")
                  return null
            }
            return data
      }
      catch (error) {
            console.error(`error getting message ${error}`)
            return null
      }
}