// Rules meant to be used with Vuetify form elements.

export function createMaxLength(len : number) : (_: string) => boolean | string {
    return (v : string) => !!v && v.length <= len || `Invalid input length, must have less than ${len} characters.`;
}

export function required(v : any) : boolean | string {
    return (!!v && v != Object()) || (Array.isArray(v) && v.length == 0)  || "Input required.";
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
