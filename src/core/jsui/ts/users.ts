export function createUserString(user : User | null) : string {
    if (!user) {
        return "No User"
    }
        
    return `${user.FirstName} ${user.LastName} (${user.Email})`
}
