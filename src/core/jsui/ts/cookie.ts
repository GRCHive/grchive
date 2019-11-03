// Source: https://developer.mozilla.org/en-US/docs/Web/API/Document/cookie/Simple_document.cookie_framework

export function getCookie(sKey : string) : string | null {
    return decodeURIComponent(document.cookie.replace(new RegExp("(?:(?:^|.*;)\\s*" + 
        encodeURIComponent(sKey).replace(/[\-\.\+\*]/g, "\\$&") + "\\s*\\=\\s*([^;]*).*$)|^.*$"), "$1")) || null;
}
