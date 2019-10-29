export function replaceWithMark(text : string, mark : string) : string {
    var re = new RegExp(`(${mark})`, 'gi')
    return text.replace(re, '<mark>$1</mark>')
}

export function sanitizeTextForHTML(text: string) : string {
    return text
        .replace('&', '&amp;')
        .replace('<', '&lt;')
        .replace('>', '&gt;')
        .replace('"', '&#34;')
        .replace('\'', '&#39;')
}
