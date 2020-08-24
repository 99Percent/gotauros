# Go Tauros API

Golang wrapper to access Tauros REST API endpoints




# Usage #

## Obtain api key and api secret ##

1. Create an account at https://tauros.io (or https://staging.tauros.io for sandbox environment), verify email and phone 
2. Go to Perfil de Usuario>General
3. Activate "Modo desarollador" with the toggle. It will ask for your account password
4. The API option will appear on the left side bar, click on it.
5. Click "Crear API key"
6. Enter any name for "Nombre del Token", click all the checkboxes except "Verificar IPs" 
7. Click "Crear API key", a modal will appear showing the API Key and API secret. This is what you will copy and paste in the json file as explained below.

## Create json tokens file ###

To run create a json file ```tokens.json``` like this

```json
{
  "email" : "[email of account]",
  "api_key" : "[api key from tauros]",
  "api_secret" : "[api secret from tauros]",
  "URL": "[https://staging.tauros.io] (sandbox) OR [https://tauros.io] (live)"
}
```
## Example:

```golang
import (
  "encoding/json"
  "io/ioutil"
  "log"
  "github.com/99Percent/gotauros"
)
// declare tauros object
  var tauros TauApi
// get API credentials from token json file
	in, err := ioutil.ReadFile("tokens.json")
	if err != nil {
		log.Fatalf("Unable to load tokens file tokens.json: %v", err)
	}
	if err := json.Unmarshal(in, &tauros); err != nil {
		log.Fatalf("Unable to unmarshall tokens file: %v", err)
  }
// use tauros object
  coins, _ := tauros.Getcoins()
  log.Printf("Available coins: %v",coins)

```