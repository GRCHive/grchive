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

export const otherTypeId = 2
const otherTypeName = "Other"

export function getDbTypeAsString(db : Database) : string {
    let typ = MetadataStore.state.idToDbType[db.TypeId]

    if (typ.Name == otherTypeName) {
        return `${typ.Name} (${db.OtherType})`
    } else {
        return typ.Name
    }
}
