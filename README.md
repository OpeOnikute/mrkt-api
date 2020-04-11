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

## Location Ranking
Locations will be ranked using a 5-day average of the number of level-3 and above incidents reported within a 5km radius. This is possible by taking advantage of Mongo's location GeoJSON and 2dsphere indexes.
| Rank        |    Average   |   Color    | 
| ---         |      ---     |    ---     |
| Safe        |    0 - 0.4   |   Green    | 
| Warning     |    0.5 - 0.9 |   Orange   |
| Unsafe      |     >= 1     |    Red     |

**N.B.** Further down the line we can use more incident levels and any other new features like incident upvotes to rank locations. For now, we can just stick to number of level-3 incidents and above reported.

## Building Docker Image
Regular Docker
- `docker build . -t opeo/mrkt-api`
- `docker login`
- `docker push opeo/mrkt-api`
Kubernetes (Local)
- eval $(minikube docker-env)
- `docker build . -t opeo/mrkt-api`
- Build mongo: `docker build -t opeo/mongo-auth -f mongo.Dockerfile .`

## TODO (Up next)
- [x] Sign up, login, post entries as user (anonymous or actual user id)
- [x] Alert types (other, potential harm, emergency, accident, fire, robbery) and their priority levels (1, 2, 3, 4, 5)
- [x] Meerkat ranking (alpha, beta, pup)
- [x] Ranking locations (Safety score)
- [x] Local Docker setup
- [ ] Add anonymous option when a user creates.
- [x] Kubernetes Setup (Local)
- [ ] Kubernetes Setup (Digital Ocean)
- [ ] Kubernetes Job (Calculate Alpha Ranking at 12am daily)
- [ ] Custom error message for all validation fields. The default one sucks.
- [ ] Tests
- [ ] Mongo driver: before find/find all, add { status: "enabled" }
- [ ] Forgot Password
- [ ] Better response than "mongo: no documents in result"
- [ ] Create location geoJSON from API not client
- [ ] Pass error instance to error handler and log stack trace properly
- [ ] Config package

## Kubernetes
Shortcut
- `kube-mrkt` ===> `kubectl --kubeconfig='mrkt-api-kubeconfig.yaml'`

Cluster-level resources
- Service Account (cicd)
- Role `~/kube-general/cicd-role.yml`
- Role binding `~/kube-general/cicd-role-binding.yml`
- Command `kube-mrkt apply -f ~/kube-general/ --kubeconfig="mrkt-api-kubeconfig.yaml"`

Port forwarding
- `kube-mrkt port-forward $(kube-mrkt get pod --selector="app=mrkt-api" --output jsonpath='{.items[0].metadata.name}') 8080:12345`

Secrets
- Create `kube-mrkt create secret generic mrkt-api-secrets --from-literal=MONGO_URL="mongodb+srv://<username>:<pwd>@<insert-url-here>" --from-literal=MONGO_DATABASE=mrkt --from-literal=PORT=12345 --from-literal=JWT_KEY="ssdsdsdsd" --from-literal=GOOGLE_MAPS_KEY=sdsdsd`

Get CI CD token
- `kube-mrkt get secret $(kube-mrkt get secret | grep cicd-token | awk '{print $1}') -o jsonpath='{.data.token}' | base64 --decode`

## Ideas
- Upvotes on incidents.
- Families/clans. Communities will be clans and you can invite people to join your clan.
- Government agencies. Each incident type would have a particular agency that it gets routed to. Admins would approve before it is passed on, or we would set some criteria for automatic routing. 

# Links
- JWT 
    - https://www.sohamkamani.com/golang/2019-01-01-jwt-authentication/
- Password hash 
    - https://medium.com/@jcox250/password-hash-salt-using-golang-b041dc94cb72
- Env variables 
    - https://dev.to/craicoverflow/a-no-nonsense-guide-to-environment-variables-in-go-a2f
- Request validation 
    - https://medium.com/@apzuk3/input-validation-in-golang-bc24cdec1835
- Kubernetes 
    - [Setup/Pipeline](https://www.digitalocean.com/community/tutorials/how-to-automate-deployments-to-digitalocean-kubernetes-with-circleci)
    - [Running local images on k8s](https://dzone.com/articles/running-local-docker-images-in-kubernetes-1)