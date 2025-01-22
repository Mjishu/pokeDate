
export type StrongPassword = {
      symbol: boolean;
      uppercase: boolean;
      lowercase: boolean; 
      number: boolean;
      isLength: boolean
      length: number;
}

const specialCharacters = ["!", "@", "#", "$", "%", "^", "&","*","(",")", "-", "+", "=", "{", "}", "\\", "|", ";", ":", '"', "<", ">", "/", ",", "'"]

export function ValidatePassword(password: string): StrongPassword {
      const strength: StrongPassword = {
            symbol: false,
            uppercase: false,
            lowercase: false,
            number:false,
            isLength: password.length >= 6,
            length: password.length
      }
      for (let i = 0; i <password.length; i ++ ){
            if ("a".charCodeAt(0) <= password.charCodeAt(i) && password.charCodeAt(i) <= "z".charCodeAt(0)) {
                  strength.lowercase = true
            } 
            if ("A".charCodeAt(0) <= password.charCodeAt(i) && password.charCodeAt(i) <= "Z".charCodeAt(0)) {
                  strength.uppercase = true
            }
            if ("0".charCodeAt(0) <= password.charCodeAt(i) && password.charCodeAt(i) <= "9".charCodeAt(0)) {
                  strength.number = true
            }
            for(let special = 0; special < specialCharacters.length ; special++) {
                  const ord = specialCharacters[special].charCodeAt(0)
                  if (password.charCodeAt(i) == ord) {
                        strength.symbol = true
                  }
            }
      }
      console.log(strength)
      return strength
}