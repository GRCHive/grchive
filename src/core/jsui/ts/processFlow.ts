export function isProcessFullDataEmpty(data : FullProcessFlowData) : boolean {
    if (!data) {
        return true
    }

    if (!data.Nodes || data.NumNodes == 0) {
        return true
    }

    return false
}
