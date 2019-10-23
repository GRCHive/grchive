interface TGetAllOrgUsersInput {
    org : string
    csrf : string
}

interface TGetAllOrgUsersOutput {
    data: User[]
}
