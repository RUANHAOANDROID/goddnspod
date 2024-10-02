
interface Config {
    userAgent: string;
    tokenId: string;
    loginToken: string;
    domain: string;
    subDomain: string;
    timer: string;
    support: string;
}

export function getCurrentUrl() :string{
    const currentUrl = new URL(window.location.href);
    currentUrl.port = "6565"; 
    return currentUrl.href; 
}

export const defConfig: Config = {
    userAgent: "Hao88 DDNS/0.1Alpha(52927295@qq.com)",
    tokenId: "xxxxxx",
    loginToken: "xxxxxxxxxxxxxx",
    domain: "hao88.cloud",
    subDomain: "xxxx",
    timer: "5m30s",
    support: "v4"
}
