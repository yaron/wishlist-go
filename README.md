# Wishlist in GO
This wishlist is a simple go based wishlist backend that can be used to create a webapplication to host your wishlist. This only includes a backend which can be contacted through it's API.

Documentation of the API can be found on https://yaron.github.io/wishlist-go/ .
I'll be building a frontend here: https://github.com/yaron/wishlist-front .

## Usage
To run in production just get the binary from the last release and run it. You can use these environment variables to change the way it works a bit.
| Variable                | Description       |
| ----------------------- | -----------------| 
| WISH_PATH | Change the path in which the database file and hash file are stored. This can be usefull when running multiple instances. Default is "".|
| WISH_DEBUG | Set this to "1" to create a test user if the database does not yet exist. The user will have the username "test2" and the password "test". |
| WISH_MAILGUN_DOMAIN | The mailgun domain to use for sending 'unclaim' mails. |
| WISH_MAILGUN_KEY | The mailgun key to use. |
| WISH_MAIL_FROM | The from mail address (eg. "My wishlist \<my@wishlist.com>") |
| WISH_MAIL_SUBJECT | The subject of the mail that users can get after claiming an item. |
| WISH_MAIL_BODY | The body of the mail that users can get after claiming an item. |

For both the body and the subject a templating language can be used to make the content use some variables. See here for more info on the template language: https://golang.org/pkg/text/template/ .
The variables that are avalable are these:
- ItemName string
- ItemID   int
- Key      string
- Mail     string