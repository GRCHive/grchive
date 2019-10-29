export function replaceWithMark(text : string, mark : string) : string {
    var re = new RegExp(`(${mark})`, 'gi')
    return text.replace(re, '<mark>$1</mark>')
}
