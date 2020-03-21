db.createUser({
    user: "opeo",
    pwd: "root",
    roles: [
        {
            role: "readWrite",
            db: "mrkt"
        }
    ]
});