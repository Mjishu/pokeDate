import {expect, test} from "vitest"
import { ValidatePassword } from "../validation"


const badRequirements = [
      {
            name: "password hits all requirements",
            password: "I@mB3tter",
            equal:   {symbol: true, uppercase: true, lowercase: true, number:true, length: "I@mB3tter".length, isLength: true,}
      },
      {
            name: "password hits lowercase requirements",
            password: "me",
            equal: {symbol: false, uppercase: false, lowercase: true, number:false, length: "me".length, isLength: false}
      },
      {
            name:"password hits length requirements",
            password: "mememe",
            equal: {symbol: false, uppercase: false, lowercase: true, number:false, length: "mememe".length, isLength: true}
      },
      {
            name:"password hits special char requirements",
            password: "@@@@@@",
            equal: {symbol: true, uppercase: false, lowercase: false, number:false, length: "@@@@@@".length, isLength: true}
      }, 
      {
            name: "password hits upper case requirement",
            password: "HI",
            equal: {symbol: false, uppercase: true, lowercase: false, number:false, length: "HI".length, isLength: false}
      }
]

for(const req of badRequirements) {
      test(req.name, () => {
            expect(ValidatePassword(req.password)).toEqual(req.equal)
      })
}
