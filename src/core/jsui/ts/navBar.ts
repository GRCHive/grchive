export interface NavBarItem {
    title: string
    icon: string
    url?: string
    disabled? : boolean
    hidden? : boolean
    children?: NavBarItem[]
}
