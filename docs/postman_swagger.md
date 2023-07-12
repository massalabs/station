# Use postman with the swagger spec

**1. Locate the last version of swagger.yml file under `api/swagger/server/restapi/resource/swagger.yml`**

![](images/postman_swagger/1.png)

**2. Copy the whole content of `swagger.yml`**

![](images/postman_swagger/2.png)

**3. You can also copy the content of `swagger.yml` directly from your local repository.**

![](images/postman_swagger/3.png)

**4. Now, install and Open your Postman**

![](images/postman_swagger/4.png)

**5. Click on APIs**

![](images/postman_swagger/5.png)

**6. Click on Import**

![](images/postman_swagger/6.png)

**7. Select Raw text**

![](images/postman_swagger/7.png)

**8. Paste the content of the `swagger.yml`**

![](images/postman_swagger/8.png)

**9. Click on Continue**

![](images/postman_swagger/9.png)

**10. Now Postman will detect MassaStation API and its type.**

![](images/postman_swagger/10.png)

**11. Click Import and wait**

![](images/postman_swagger/11.png)

**12. When Import is Done, you will be able to see MassaStation API imported, and you can browse all its endpoints**

![](images/postman_swagger/12.png)

**13. Let's try to call an endpoint, but before you should have your massa station running in the background. You can run it directly from the terminal, inside station repo, type `task build-run`**

![](images/postman_swagger/13.png)

**14. Expand MassaStation folder, and let's try to call an endpoint**

![](images/postman_swagger/14.png)

15. **Replace the `{{baseUrl}}` by `station.massa` or define it as an environment variable in Postman:**
   
   For more flexibility, define `baseUrl` as your own environment variable in Postman. This allows you to easily switch between different environments or server instances without modifying the requests individually.
   
   Refer to the [Postman documentation](https://learning.postman.com/docs/sending-requests/variables/) for details on defining and managing environment variables in Postman.



![](images/postman_swagger/15.png)

**16. Click on Send, and wait a few seconds...**

![](images/postman_swagger/16.png)

**17. You will receive in response all websites uploaded so far to the current DNS.**

![](images/postman_swagger/17.png)
