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
      date_of_birth: Date | undefined;
      price: number;
      available: boolean;
      sex: string;
      breed: string;
}

export async function createAnimal(e: Event, animal: NewAnimal) {
      e.preventDefault();
      if (animal == undefined) {
            console.error("cannot make an empty animal")
            return
      }
      const fetchParams = {
            method: "POST",
            headers: {
                  "Content-Type": "application/json",
            },
            body: JSON.stringify({ animal })
      }
      console.log("creating: " + animal)
      // try {
      //       const response = await fetch("http://localhost:8080/animals", fetchParams)
      //       const data = await response.json()
      //       return data
      // } catch (error) {
      //       console.error(`Error trying to create animal: ${error}`)
      //       return
      // }
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