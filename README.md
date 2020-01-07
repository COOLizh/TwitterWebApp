# TwitterWebApp

## Final project for courses at epam

### Hard level

- [x] POST /register - create new account with specified nick(unique in app), email, and password

- [x] POST /login - accept email and password and return token, you can use JWT or save session id in database

- [x] POST /subscribe - add account with login to your subscription list, you start seeing his tweets in your feeds

- [x] POST /tweets - create a tweet, account id should be found from JWT or fetched from database using session id

- [x] GET /tweets - return all tweets from your subscriptions
