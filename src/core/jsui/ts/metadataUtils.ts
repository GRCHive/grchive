import Metadata from './metadata'

export function lazyGetUserFromId(userId : number | null) : Promise<User | null> {
    return new Promise<User | null>((resolve) => {
        let getUser = () => {
            if (!!userId && userId in Metadata.state.idToUsers) {
                resolve(Metadata.state.idToUsers[userId])
            } else {
                resolve(null)
            }
        }

        if (Metadata.state.usersInitialized) {
            getUser()
        } else {
            // TODO: Work under the assumption that the initialize function
            // was already call. We should break that assumption and do the
            // init here probably.
            Metadata.watch((state) => {
                return state.usersInitialized
            }, () => {
                getUser()
            })
        }
    })
}

export function lazyGetControlTypeFromId(controlId : number) : Promise<ProcessFlowControlType> {
    return new Promise<ProcessFlowControlType>((resolve) => {
        let getControlType = () => {
            if (controlId in Metadata.state.idToControlTypes) {
                resolve(Metadata.state.idToControlTypes[controlId])
            } else {
                resolve(Object() as ProcessFlowControlType)
            }
        }

        if (Metadata.state.controlTypeInitialized) {
            getControlType()
        } else {
            // TODO: Work under the assumption that the initialize function
            // was already call. We should break that assumption and do the
            // init here probably.
            Metadata.watch((state) => {
                return state.controlTypeInitialized
            }, () => {
                getControlType()
            })
        }
    })

}
