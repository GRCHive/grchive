export function createUserString(user : User) : string {
    return `${user.FirstName} ${user.LastName} (${user.Email})`
}
