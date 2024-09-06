// api接口统一风格响应
export type ApiEcho<T> = {
    bizCode: string;
    message: string;
    data: T,
    version: string;
    requestId: string;
}


export type PageInfo = {
    pageInt: number;
    pageSize: number;
}


export type LoginResponse = {
    id: number;
    email: string;
    username: string;
    avatar: string;
    role: number;
    accessToken: string;
    refreshToken: string;
}

export type NewPostRequest = {
    "title": string;
    "content": string;
    "attrId"?: number;
    "tags"?: number[];
}

