export const getStartedUrl : string = "/getting-started";
export const contactUsUrl : string = "/contact-us";
export const homePageUrl : string = "/";
export const loginPageUrl : string = "/login";

export function createAssetUrl(asset : string) : string {
    return "/static/assets/" + asset;
}

export const learnMoreUrl : string = "/learn";

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
