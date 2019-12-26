# TwitterWebApp

## Final project for courses at epam
### Low level
 - [ ] POST /tweets - create a tweet
 - [ ] GET /tweets?account_id=“your account id” - return all tweets expect yours order by timestamp newest first
### Middle level
 - [ ] POST /register - create new account with specified nick(unique in app), email, and password
 - [ ] POST /login - accept email and password and return token, you can use JWT or save session id in database
 - [ ] POST /tweets - create a tweet, account id should be found from JWT or fetched from database using session id
 - [ ] GET /tweets - return all tweets expect yours order by timestamp newest first find your account id from Token]
 ### Hard level
 
 - [ ] POST /register - create new account with specified nick(unique in app), email, and password
 - [ ] POST /login - accept email and password and return token, you can use JWT or save session id in database
 - [ ] POST /subscribe - add account with login to your subscription list, you start seeing his tweets in your feeds
 - [ ] POST /tweets - create a tweet, account id should be found from JWT or fetched from database using session id
 - [ ] GET /tweets - return all tweets from your subscriptions
 