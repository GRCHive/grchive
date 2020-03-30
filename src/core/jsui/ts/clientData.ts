export interface ClientData {
	Id          : number
	OrgId       : number
	Name        : string
	Description : string
}

export interface ClientDataVersion {
	Id      : number
	OrgId   : number
	DataId  : number
	Version : number
	Kotlin  : string
}

export interface DataSourceOption {
	Id          : number
	Name        : string
	KotlinClass : string
}

export interface DataSourceLink {
	OrgId        : number
	DataId       : number
	SourceId     : number
	SourceTarget : Record<string, any>
}

export interface FullClientDataWithLink {
    Data : ClientData
    Link: DataSourceLink
}

export class DataSourceOptionNode {
    _tag : string
    _fullTag : string
    _val : DataSourceOption | null

    _children : DataSourceOptionNode[]
    _tagToChild: Map<string, DataSourceOptionNode>

    constructor(tag : string, fullTag : string, val : DataSourceOption | null) {
        this._tag = tag
        this._fullTag = fullTag
        this._children = []
        this._tagToChild = new Map<string, DataSourceOptionNode>()
        this._val = val
    }

    addDirectChild(tag : string, val : DataSourceOption | null) : DataSourceOptionNode {
        if (this._tagToChild.has(tag)) {
            return this._tagToChild.get(tag)!
        }

        let node = new DataSourceOptionNode(tag, this._fullTag + '.' + tag, val)
        this._children.push(node)
        this._tagToChild.set(tag, node)
        return node
    }

    addNestedChild(o : DataSourceOption) {
        let tagsToParse = o.Name.replace(this._fullTag + '.', '').split('.')
        if (tagsToParse.length == 1) {
            this.addDirectChild(tagsToParse[0], o)
        } else {
            let next = this.addDirectChild(tagsToParse[0], null)
            next.addNestedChild(o)
        }
    }

    printNode(prefix : string = "") {
        console.log(`${prefix}${this._tag} [${this._fullTag}]`)
        for (let c of this._children) {
            c.printNode(prefix + '\t')
        }
    }

    get numChildren() : number {
        return this._children.length
    }

    get isLeaf() : boolean {
        return (this.numChildren == 0 && !!this._val)
    }

    getChildTags() : string[] {
        return Array.from(this._tagToChild.keys())
    }

    getChild(t : string) : DataSourceOptionNode | null {
        if (!this._tagToChild.has(t)) {
            return null
        }
        return this._tagToChild.get(t)!
    }
}

export class DataSourceOptionTree {
    _options : DataSourceOption[]

    _rootNode : DataSourceOptionNode

    constructor(options : DataSourceOption[]) {
        this._options = options

        this._rootNode = new DataSourceOptionNode("Root", "Root", null)
        for (let o of options) {
            this._rootNode.addNestedChild(o)
        }
    }

    traverse(tags : string[]) : DataSourceOptionNode | null {
        let node : DataSourceOptionNode | null = this._rootNode
        let tagCopy : string[] = tags.slice()
        while (tagCopy.length > 0) {
            let t = tagCopy.shift()
            node = node.getChild(t!)
            if (!node) {
                return null
            }
        }
        return node
    }

    printTree() {
        this._rootNode.printNode()
    }
}
