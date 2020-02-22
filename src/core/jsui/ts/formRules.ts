// Rules meant to be used with Vuetify form elements.

export function createMaxLength(len : number) : (_: string) => boolean | string {
    return (v : string) => !!v && v.length <= len || `Invalid input length, must have less than ${len} characters.`;
}

export function createMinLength(len : number) : (_: string) => boolean | string {
    return (v : string) => !!v && v.length >= len || `Invalid input length, must have more than ${len} characters.`;
}

export function nonZero(v : any) : boolean | string {
    return (v != 0) || "Invalid choice.";
}

export function required(v : any) : boolean | string {
    return (!!v && v != Object()) || (Array.isArray(v) && v.length == 0) || (v === 0) || "Input required.";
}

export function email(v : string) : boolean | string {
    if (!v) {
        return true;
    }

    const regex : RegExp = /(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])/;
    const match = v.match(regex);
    return !!match || "A valid email is required.";
}

export function numeric(v : string) : boolean | string {
    return !isNaN(Number(v)) || "Must be numeric."
}

export function hasLowerCase(v : string) : boolean | string {
    return (/[a-z]/.test(v)) || "Must have lower case characters."
}

export function hasUpperCase(v : string) : boolean | string {
    return (/[A-Z]/.test(v)) || "Must have upper case characters."
}

export function hasNumeric(v : string) : boolean | string {
    return (/[0-9]/.test(v)) || "Must have numeric characters."
}

export function password(v: string) : boolean | string {
    // This needs to be kept in sync w/ whatever rules we have on the server.
    //  - Min Length: 8
    //  - Has lower case
    //  - Has upper case
    //  - Has number
    let minLength = createMinLength(8)

    let rMinLength = minLength(v)
    if (!rMinLength || typeof rMinLength == 'string') {
        return rMinLength
    }

    let rLowerCase = hasLowerCase(v)
    if (!rLowerCase || typeof rLowerCase == 'string') {
        return rLowerCase
    }

    let rUpperCase = hasUpperCase(v)
    if (!rUpperCase || typeof rUpperCase == 'string') {
        return rUpperCase
    }

    let rNumeric = hasNumeric(v)
    if (!rNumeric || typeof rNumeric == 'string') {
        return rNumeric
    }

    return true
}

export function createPerElement(fn : (arg0 : string) => boolean | string) : (arg0 : Array<string>) => boolean | string {
    return function(v: Array<string>) : boolean | string  {
        for (const ele of v) {
            const res = fn(ele)
            if (!res || typeof res == 'string') {
                return res + ` [${ele}]`
            }
        }
        return true
    }
}
