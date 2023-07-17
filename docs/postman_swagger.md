# Using Postman with Swagger

This guide provides a high-level overview of using Postman with Swagger to interact with an API. Follow these 5 steps to get started:

1. **Locate the Swagger File**: Find the `swagger.yml` [file](https://github.com/massalabs/thyra/blob/e25eef54d5901ea16dddd7258ee09481a3e794a0/api/swagger/server/restapi/resource/swagger.yml) for the API you want to work with and copy its entire content.

2. **Importing the Swagger**: Launch Postman and navigate to the "APIs" section. From there, click on the "Import" button and select the "Raw text" option. Next, paste the contents of the copied `swagger.yml` file into the designated text field.

3. **Ensure MassaStation is Running**: Start your MassaStation by using the command `task build-run`.

4. **Interact with an Endpoint**: Expand the imported API in Postman and select the desired endpoint (e.g., "all Domains Getter"). Replace `{{baseUrl}}` with `station.massa` by clicking on the "..." menu, then go to the Variables tab, and add the variable. From there, Click on the "Send" button in Postman to make the API request

Congratulations! You've successfully used Postman with Swagger to interact with an API.
