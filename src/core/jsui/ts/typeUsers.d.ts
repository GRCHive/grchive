interface User {
    Id: number
    FirstName : string
    LastName : string
    Email : string
}

interface UserWithRole {
    User: User
    RoleId: number
    OrgId: number
}
