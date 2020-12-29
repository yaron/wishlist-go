# Wishlist in GO
This wishlist is a simple go based wishlist backend that can be used to create a webapplication to host your wishlist. This only includes a backend which can be contacted through it's API.

## Usage
To run in production just get the binary from the last release and run it. You can use these environment variables to change the way it works a bit.
| Variable                | Description       |
| ----------------------- | -----------------| 
| WISH_PATH | Change the path in which the database file and hash file are stored. This can be usefull when running multiple instances. Default is "".|
| WISH_DEBUG | Set this to "1" to create a test user if the database does not yet exist. The user will have the username "test2" and the password "test". |
| WISH_MAILGUN_DOMAIN | The mailgun domain to use for sending 'unclaim' mails. |
| WISH_MAILGUN_KEY | The mailgun key to use. |
| WISH_URL | The domain to use in the 'unclaim' mail. Do not unclude https:// (eg. example.com). "/unclaim/[secret-key]" will be appended, you need to make sure something exists there. |
| WISH_NAME | The name of the site, to be used in the 'unclaim' mail.| 

## TODO
 - Document API
 - Make mail text and subject overwritable
 - Build frontend