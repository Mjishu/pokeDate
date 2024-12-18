export type Animal = {
      Id: string;
      Name: string;
      Species: string;
      Date_of_birth: string;
      Sex: string;
      Available: boolean;
      Price: number;
      Breed: string;
      Image_src: string;
      Shots: AnimalShot[]
};

export type Shot = {
      Id: number;
      Name: string;
      Description: string;
}

export type AnimalShot = {
      id: string;
      name: string;
      description: string;
      date_given: string;
      next_due: string
}

export type NewAnimal = {
      name: string;
      species: string;
      date_of_birth: string | undefined;
      price: number;
      available: boolean;
      sex: string;
      breed: string;
      shots: NewShot[];
}

type NewShot = {
      id: string;
      date_given: string;
      date_due: string;
}

export type UpdatedAnimal = {
      name: string;
      date_of_birth: string | undefined;
      price: number;
      available: boolean;
      shots: NewShot[];
}


export async function createAnimal(animal: NewAnimal) {
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
            const response = await fetch("http://localhost:8080/animals", fetchParams)
            const data = await response.json()
            return data
      } catch (error) {
            console.error(`Error trying to create animal: ${error}`)
            return
      }
}

export async function getOrganizationAnimals() {
      try {
            const response = await fetch("http://localhost:8080/organizations/animals")
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
            const response = await fetch("http://localhost:8080/organizations/animals", fetchParams)
            const data = await response.json()
            console.log(data)
            return data
      } catch (error) {
            console.error(`error trying to get animal: ${error}`)
            return
      }
}

export async function updateAnimalById(e: Event, id: string, updatedAnimal: UpdatedAnimal): Promise<void> {
      e.preventDefault()
      const fetchParams = {
            method: "PUT",
            headers: {
                  "Content-Type": "application/json"
            },
            body: JSON.stringify({
                  id: id, name: updatedAnimal.name, date_of_birth: updatedAnimal.date_of_birth + "T00:00:00Z",
                  price: updatedAnimal.price, available: updatedAnimal.available
            })
      }
      try {
            const response = await fetch("http://localhost:8080/animals", fetchParams)
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
            const response = await fetch("http://localhost:8080/animals", fetchParams)
            const data = await response.json()
            return data
      } catch (error) {
            console.error(`error trying to delete animal: ${error}`)
            return
      }
}

export async function GetAllShots(): Promise<Shot[]> {
      try {
            const response = await fetch("http://localhost:8080/shots")
            const data = await response.json()
            return data
      } catch (error) {
            console.error(`there was an error trying to fetch all shots: ${error}`)
            return []
      }
}