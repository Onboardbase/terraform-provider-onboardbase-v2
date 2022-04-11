const express = require('express')
const  CryptoJS  = require('crypto-js');


const decryptSecrets = async (
    secrets,
    passcode
  ) => {
    const encryptionPassphrase = passcode;
    try {
      const bytes = CryptoJS.AES.decrypt(secrets, encryptionPassphrase);
      return bytes.toString(CryptoJS.enc.Utf8);
    } catch (error) {
      console.log(error);
    }
  };

const aesDecryptSecret = (secret, passcode) => {
    return decryptSecrets(secret, passcode);
  };

const app = express()
app.use(express.json());

app.post("/decode", async (req, res) => {
    let secrets = req.body.secrets
    let wanted = req.body.secret
    let passcode = req.body.passcode
    let found = false
    let decoded = {}
    for (index in secrets){
        try{
        let r = await aesDecryptSecret(secrets[index], passcode)
        r = JSON.parse(r)
        console.log(r);
        if(r.key == wanted){
            found = true
            console.log(found);
            return res.json(r)
        }
        } catch(err){
            return res.json({error: "Unable to decrypt secret. Passcode may be invalid"})
        }
    }
    if (found == false) return res.json({error: "Secret not found"})
    
    
})

app.listen("3000", () => console.log("Listening"))