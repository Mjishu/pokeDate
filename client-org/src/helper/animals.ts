import { formatISO } from "date-fns";

export type Animal = {
      Id: string;
      Name: string;
      Species: string;
      Date_of_birth: string;
      Sex: string;
      Price: number;
      Available: boolean;
      Breed: string;
      Image_src: string;
      Shots: AnimalShot[] | NewShot[]
};

export type Shot = {
      Id: number;
      Name: string;
      Description: string;
}

export type AnimalShot = {
      Id: number;
      Name: string;
      Description: string;
      Date_given: string;
      Next_due: string
}

export type NewAnimalImage = {
      animal_id: string;
}

export type NewAnimal = {
      Name: string;
      Species: string;
      Date_of_birth: string | undefined;
      Price: number;
      Available: boolean;
      Sex: string;
      Breed: string;
      Shots: NewShot[];
}

type NewShot = {
      Id: number;
      Date_given: string;
      Next_due: string;
}

export type UpdatedAnimal = {
      Name: string;
      Date_of_birth: string | undefined;
      Price: number;
      Available: boolean;
      Shots: NewShot[];
}

export async function createAnimal(animal: NewAnimal, image: File) {
      //! Switch date_birth to be in the 2024-10-09T00:00:00Z format
      console.log(animal)
      if (animal == undefined) {
            console.error("cannot make an empty animal")
            return
      }
      const fetchParams = {
            method: "POST",
            headers: {
                  "Content-Type": "application/json",
            },
            body: JSON.stringify({ ...animal })
      }
      console.log(fetchParams.body)
      try {
            const response = await fetch("/api/animals", fetchParams)
            if (!response.ok) {
                  throw new Error(`issue uploading animal: ${response.statusText}`)
            }
            const data = await response.json()
            // createAnimalImage(image, data.animal.Id, true)
      } catch (error) {
            console.error(`Error trying to create animal: ${error}`)
            return
      }
}

export async function getOrganizationAnimals() {
      try {
            const response = await fetch("/api/organizations/animals")
            const data = await response.json()
            return data
      } catch (error) {
            console.error(`Error trying to get all animals: ${error}`)
            return
      }
}

export async function getAnimalById(id: string) {
      const fetchParams = { //causes error on backend bc it thinks its getting a new animal to create
            method: "POST",
            headers: {
                  "Content-Type": "application/json"
            },
            body: JSON.stringify({ id: id })
      }
      try {
            const response = await fetch("/api/organizations/animals", fetchParams)
            const data = await response.json()
            return data
      } catch (error) {
            console.error(`error trying to get animal: ${error}`)
            return
      }
}

export async function updateAnimalById(id: string, updatedAnimal: UpdatedAnimal): Promise<void> {
      const fetchParams = {
            method: "PUT",
            headers: {
                  "Content-Type": "application/json"
            },
            body: JSON.stringify({
                  ...updatedAnimal,
                  Id: id,
            })
      }
      try {
            console.log(fetchParams.body)
            const response = await fetch("/api/animals", fetchParams)
            const data = await response.json()
            return data
      } catch (error) {
            console.error(`error trying to update animal: ${error}`)
            return
      }
}

export async function DeleteAnimalById(id: string): Promise<void> {
      const fetchParams = {
            method: "DELETE",
            headers: {
                  "Content-Type": "application/json"
            },
            body: JSON.stringify({ id })
      }
      try {
            const response = await fetch("/api/animals", fetchParams)
            const data = await response.json()
            return data
      } catch (error) {
            console.error(`error trying to delete animal: ${error}`)
            return
      }
}

export async function GetAllShots(): Promise<Shot[]> {
      try {
            const response = await fetch("/api/shots")
            const data = await response.json()
            return data
      } catch (error) {
            console.error(`there was an error trying to fetch all shots: ${error}`)
            return []
      }
}