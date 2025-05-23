export interface ILoginResponse {
    message: string,
    user_id: string,
    full_name: string
}

export interface ILoginRequest {
    phone: string
}