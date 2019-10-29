export function nodeTypeToClass(typeId : number) : string {
    // TODO: How do we keep this in sync with the server?
    //       Maybe it should get queried along with the types.
    switch(typeId){
        case 1:
            return "activity-manual"
        case 2:
            return "activity-automated"
        case 3:
            return "decision"
        case 4:
            return "start"
        case 5:
            return "general-ledger-entry"
        case 6:
            return "system"
        default:
            break;
    }
    return ""
}
