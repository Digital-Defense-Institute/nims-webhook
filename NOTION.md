 # Notion Configuration
 
 ## Database IDs

You need both the `Asset Database` ID and the `Alert Database` ID to use this tool.

1. Navigate to the database within NIMS under `Databases`
2. Right click on the database and click `Copy link`
    ![link](./screenshots/link.png)
3. Locate the database ID in the URL
    * The database ID is the long string of letters and numbers in the URL after the last `/` and before the `?` or `#` if present
    * Example: 
        * Link: `https://www.notion.so/184cdc5a1ef3710badc2d2b1271aeb81?v=174cdc3a1ef181719981000cab12bf54&pvs=4`
        * ID: `184cdc5a1ef3710badc2d2b1271aeb81`
4. Copy the ID 
5. Repeat the above for the other database

 ## Auth Token and Access

This will walk you through creating a Notion integration, getting the auth token, and adding the integration to the proper NIMS databases.

1. Go to `Manage connections` in Notion
    ![connection](./screenshots/connection.png)

2. Click `Develop or manage integrations`
    ![connection](./screenshots/manage.png)

3. Click `New integration`
    ![connection](./screenshots/new.png)

4. Configure the new integration
    * Give it a name, ex: `nims_template`
    * Choose the workspace
    * Type: `Internal`
    * Click `Save`
    ![connection](./screenshots/integration.png)

5. Click `Configure integration settings`
    ![connection](./screenshots/configure.png)

6. Copy the `Internal Integration Secret` -- this is your auth token for your `.env` file
    * Click `Save`
     ![connection](./screenshots/token.png)

7. Navigate to your `Alert Database`
    * Click the 3-dot menu and find `Connections`
    * Click on your newly created integration
    ![connection](./screenshots/alerts.png)

8. Click `Confirm`  
    ![connection](./screenshots/confirm.png)

9. Repeat steps 7 and 8 for the `Asset Database`