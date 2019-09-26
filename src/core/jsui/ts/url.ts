export function createGetStartedUrl() : string {
    return createContactUsUrl();
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
