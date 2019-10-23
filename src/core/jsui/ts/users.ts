export function createUserString(user : User) : string {
    return `${user.firstName} ${user.lastName} (${user.email})`
}
