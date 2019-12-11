export interface Comment {
	Id      :  number
	UserId  :  number
	PostTime:  Date
	Content :  string
}

export enum CommentResource {
    DocumentRequest = 1 
}
