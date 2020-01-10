export function sanitizeUrl(url : string) : string {
    if (url.startsWith('http')) {
        return url
    }

    return '//' + url
}
