export function createGetStartedUrl() : string {
    return "/getting-started";
}

export function createContactUsUrl() : string {
    return "/contact-us";
}

export function createHomePageUrl() : string {
    return "/";
}

export function createLoginUrl() : string {
    return "/login";
}

export function createAssetUrl(asset : string) : string {
    return "/static/assets/" + asset;
}

export function createLearnMoreUrl() : string {
    return "/learn";
}

export function createMailtoUrl(user : string, domain : string) : Object {
    const email = createEmailAddress(user, domain)
    return Object.freeze({
        mailto: "mailto:" + email,
        email: email
    });
}

export function createEmailAddress(user : string, domain : string) : string {
    return user + "@" + domain;
}
