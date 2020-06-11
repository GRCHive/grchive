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

export interface DbFunction {
    Id          : number
    OrgId       : number
    SchemaId    : number
    Name        : string
    Src         : string
    RetType     : string | null
}

export interface RawDbColumn {
    Name        : string
    Type        : string
}

export interface DbTable {
    Id          : number
    OrgId       : number
    SchemaId    : number
    TableName   : string
    Columns     : RawDbColumn[]
}

export interface DbSqlQueryMetadata {
    Id          : number
    DbId        : number
    OrgId       : number
    Name        : string
    Description : string
}

export interface DbSqlQuery {
    Id           : number
    MetadataId   : number
    Version      : number
    UploadTime   : Date
    UploadUserId : number
    OrgId        : number
    Query        : string
}

export interface SqlResult {
    Columns: string[]
    CsvText: string
}

export interface DbSqlQueryRequest {
    Id              : number
    QueryId         : number
    UploadTime      : Date
    UploadUserId    : number
    AssigneeUserId  : number | null
    DueDate         : Date | null
    OrgId           : number
}

export function cleanDbSqlQueryFromJson(q : DbSqlQuery) {
    q.UploadTime = new Date(q.UploadTime)
}

export function cleanDbSqlRequestFromJson(q : DbSqlQueryRequest) {
    q.UploadTime = new Date(q.UploadTime)

    if (!!q.DueDate) {
        q.DueDate = new Date(q.DueDate)
    }
}

export interface DbSqlQueryRequestApproval {
	RequestId        : number
	OrgId            : number
	ResponseTime     : Date
	ResponsderUserId : number
	Response         : boolean
	Reason           : string
}

export function cleanDbSqlRequestApprovalFromJson(q : DbSqlQueryRequestApproval) {
    q.ResponseTime = new Date(q.ResponseTime)
}
