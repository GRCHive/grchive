import MetadataStore from './metadata'

export interface DatabaseType {
    Id: number
    Name: string
}

export interface Database {
    Id: number
    Name: string
    OrgId: number
    TypeId: number
    OtherType: string
    Version: string
}

export interface DatabaseConnection {
    Id: number
    DbId: number
    OrgId: number
    ConnString: string
    Username: string
    Password: string
}

export const otherTypeId = 2
const otherTypeName = "Other"

export function isDatabaseSupported(db: Database) : boolean {
    return db.TypeId != otherTypeId
}

export function getDbTypeAsString(db : Database) : string {
    let typ = MetadataStore.state.idToDbType[db.TypeId]

    if (typ.Name == otherTypeName) {
        return `${typ.Name} (${db.OtherType})`
    } else {
        return typ.Name
    }
}
