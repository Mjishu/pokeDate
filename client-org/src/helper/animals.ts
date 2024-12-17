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
};

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
      name: string;
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
            const response = await fetch("http://localhost:8080/animals", fetchParams)
            const data = await response.json()
            console.log(data)
            return data
      } catch (error) {
            console.error(`error trying to get animal: ${error}`)
            return
      }
}

export async function updateAnimalById(id: string, updatedAnimal: UpdatedAnimal) {
      const fetchParams = {
            method: "PUT",
            headers: {
                  "Content-Type": "application/json"
            },
            body: JSON.stringify({ id })
      }
      try {
            console.log(updatedAnimal)
            // const response = await fetch("http://localhost:8080/animals", fetchParams)
            // const data = await response.json()
            // return data
      } catch (error) {
            console.error(`error trying to update animal: ${error}`)
            return
      }
}