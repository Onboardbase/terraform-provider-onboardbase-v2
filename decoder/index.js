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

exports.handler = async (event) => {
  // TODO implement
  if (event.body){
    let body = JSON.parse(event.body)
    let secrets = body.secrets
    let wanted = body.secret
    let passcode = body.passcode
    let found = false
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
            return JSON.stringify({error: "Unable to decrypt secret. Passcode may be invalid"})
        }
    }
    if (found == false) return JSON.stringify({error: "Secret not found"}) 
  }
  return JSON.stringify({"error": "No request body specified"});
};
