# Mrkt
This is a rewrite of [this application](https://github.com/OpeOnikute/safety-alert-api) in Golang. 

## Authorization
This is done using JWTs. Authorized endpoints require the JWT token be passed as a Bearer Token in a `Authorization` header.
- Call the `/login` endpoint and then store the token.
- When calling an authorised endpoint, pass in a header called `Authorization` with the value `Bearer <token>`.

## Alert Types
These are available for users to select when creating the entry. When they select one, the priority is automatically assigned. The types are managed from the admin so they can be dynamic. They are added to an entry by passing just the ID.
The priority levels are loosely based on [DEFCON](https://en.wikipedia.org/wiki/DEFCON). 

### Levels
| Level |       Description                     | 
| ---   |          ---                          | 
| 5     | Highest level. Very serious           | 
| 4     | Important. Should be resolved asap.   | 
| 3     | Getting serious.                      | 
| 2     | Potential problem.                    | 
| 1     | Minor.                                | 

### Types
| Type | Level  | Example  | 
| ---  | ---    |    ---       |
| potential harm | 2 |  Naked wire etc.  | 
| emergency | 4  | Fire, blast etc.  |  
| accident  | 4  | Car, Danfo etc  | 
| fire  | 5 |  --  | 
| robbery  | 5 | --  | 

## Meerkat Ranking
To make this more fun, we want to make users see their ranking. They are:
| Rank | Description  | Criteria  | 
| ---  | ---    |    ---       |
| Pup | You are just getting started. You have a lot to learn, and it's great that you have a community willing to help. |  Less than 5 incidents reported.  | 
| Beta | You are starting to find your feet. There is more to come from you.  | Less than 20 incidents reported.  |  
| Alpha  | You are amongst the elite. You care about the safety of the collective, and your contributions are making an impact.  | Up to 20 incidents reported.  | 

### Notes
- The criteria should be built in a dynamic way, meaning users can be demoted if we shift the goal posts. This is fair because their ranking should be in relation to the rankings of others in the clan.
- With time, the criteria of number of incidents will change e.g. 5 becomes 10, 20 becomes 40 etc. But I believe it's good for now.
- If/when we add the ability to upvote incidents, this would also factor. e.g. 5 incidents with at least 10 upvotes each etc.
- We can also determine this dynamically by getting the user with the most incidents reported and then work our way down from there. A MongoDB aggregation that splits users into buckets and then picks out the position of the user in relation to others can be done. But this would mean calculating this all the time, which can be resource intensive. To solve that, we can have a specific time ranks are updated. e.g. 12am every day. Hm.
- Final solution:
    1. Get the user with the highest number reported.
    2. When users need to see their ranks, calculate the distribution percentiles. 0-40%, 40-80%, 80-100%. i.e. four numbers including zero.
    3. Use their position in the distribution to determine their rank.

## TODO (Up next)
- [x] Sign up, login, post entries as user (anonymous or actual user id)
- [x] Alert types (other, potential harm, emergency, accident, fire, robbery) and their priority levels (1, 2, 3, 4, 5)
- [ ] Meerkat ranking (alpha, beta, pup)
- [ ] Ranking locations (Safety score)
- [ ] Custom error message for all validation fields. The default one sucks.
- [ ] Tests
- [ ] Mongo driver: before find/find all, add { status: "enabled" }
- [ ] Forgot Password
- [ ] Better response than "mongo: no documents in result"

## Ideas
- Upvotes on incidents.
- Families/clans. Communities will be clans and you can invite people to join your clan.

# Links
- JWT - https://www.sohamkamani.com/golang/2019-01-01-jwt-authentication/
- Password hash - https://medium.com/@jcox250/password-hash-salt-using-golang-b041dc94cb72
- Env variables - https://dev.to/craicoverflow/a-no-nonsense-guide-to-environment-variables-in-go-a2f
- Request validation - https://medium.com/@apzuk3/input-validation-in-golang-bc24cdec1835