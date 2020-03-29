# Mrkt
This is a rewrite of [this application](https://github.com/OpeOnikute/safety-alert-api) in Golang. 

### Authorization
This is done using JWTs. Authorized endpoints require the JWT token be passed as a Bearer Token in a `Authorization` header.
- Call the `/login` endpoint and then store the token.
- When calling an authorised endpoint, pass in a header called `Authorization` with the value `Bearer <token>`.

### Alert Types
These are available for users to select when creating the entry. When they select one, the priority is automatically assigned. The types are managed from the admin so they can be dynamic.
The priority levels are loosely based on [DEFCON](https://en.wikipedia.org/wiki/DEFCON). 

#### Levels
| Level |       Description                     | 
| ---   |          ---                          | 
| 5     | Highest level. Very serious           | 
| 4     | Important. Should be resolved asap.   | 
| 3     | Getting serious.                      | 
| 2     | Potential problem.                    | 
| 1     | Minor.                                | 

#### Types
| Type | Level  | Example  | 
| ---  | ---    |    ---       |
| potential harm | 2 |  Naked wire etc.  | 
| emergency | 4  | Fire, blast etc.  |  
| accident  | 4  | Car, Danfo etc  | 
| fire  | 5 |  --  | 
| robbery  | 5 | --  | 

### TODO (Up next)
- [x] Sign up, login, post entries as user (anonymous or actual user id)
- [ ] Alert types (other, potential harm, emergency, accident, fire, robbery) and their priority levels (1, 2, 3, 4, 5)
- [ ] Meerkat ranking (alpha, beta, pup)
- [ ] Ranking locations (Safety score)
- [ ] Custom error message for all validation fields. The default one sucks.
- [ ] Tests
- [ ] Mongo driver: before find/find all, add { status: "enabled" }
- [ ] Forgot Password
- [ ] Better response than "mongo: no documents in result"

# Links
- JWT - https://www.sohamkamani.com/golang/2019-01-01-jwt-authentication/
- Password hash - https://medium.com/@jcox250/password-hash-salt-using-golang-b041dc94cb72
- Env variables - https://dev.to/craicoverflow/a-no-nonsense-guide-to-environment-variables-in-go-a2f
- Request validation - https://medium.com/@apzuk3/input-validation-in-golang-bc24cdec1835