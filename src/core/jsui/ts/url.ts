export const getStartedUrl : string = "/getting-started";
export const contactUsUrl : string = "/contact-us";
export const homePageUrl : string = "/";
export const loginPageUrl : string = "/login";

export function createAssetUrl(asset : string) : string {
    return "/static/assets/" + asset;
}

export const learnMoreUrl : string = "/learn";
export const dashboardUrl : string = "/dashboard";

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

export const baseLogoutUrl : string = "/logout";
export function createLogoutUrl(csrf : string) : string {
    return baseLogoutUrl + "?csrf=" + csrf;
}

export const myAccountBaseUrl : string = "/dashboard/user";
export function createMyAccountUrl(email: string) : string {
    return myAccountBaseUrl + "/" + email;
}
