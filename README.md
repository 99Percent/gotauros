# How to Use #

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
  "api_key" : "[api_key from tauros]",
  "api_secret" : "[api secret from tauros]",
  "URL": "[https://staging.tauros.io] (sandbox) [https://tauros.io] (live)"
}
```