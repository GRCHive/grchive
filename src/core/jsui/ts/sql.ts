import { standardFormatTime } from './time'

export interface DbRefresh {
    Id                  : number
    DbId                : number
    OrgId               : number
    RefreshTime         : Date | null
    RefreshFinishTime   : Date | null
    RefreshSuccess      : boolean
    RefreshErrors       : string
}

export function cleanDbRefreshFromJson(db : DbRefresh) {
    if (!!db.RefreshTime) {
        db.RefreshTime = new Date(db.RefreshTime)
    }

    if (!!db.RefreshFinishTime) {
        db.RefreshFinishTime = new Date(db.RefreshFinishTime)
    }
}

export function dbRefreshIdentifier(refresh : DbRefresh) : string {
    let successStr : string = ""
    if (!refresh.RefreshFinishTime) {
        successStr = "PENDING"
    } else if (!refresh.RefreshSuccess) {
        successStr = "FAILED"
    } else {
        successStr = "COMPLETE"
    }

    return `${standardFormatTime(refresh.RefreshTime!)} [${successStr}]`
}

export interface DbSchema {
    Id          : number
    OrgId       : number
    RefreshId   : number
    SchemaName  : string
}

export interface DbTable {
    Id          : number
    OrgId       : number
    SchemaId    : number
    TableName   : string
}

export interface DbColumn {
    Id          : number
    OrgId       : number
    TableId     : number
    ColumnName  : string
    ColumnType  : string
}
