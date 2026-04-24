import { BaseResponse } from "../response.type";

export type LoginRequest = {
    email_or_username: string;
    password: string;
}

export type LoginResponse = BaseResponse<{
    session_id: string;
    full_name: string;
    email: string;
    username: string;
    photo: string;
    phone: string;
    is_admin: boolean;
    redirect_url?: string;
}>

export type OauthUserResponse = BaseResponse<{
    user_id: string;
    full_name: string;
    email: string;
    username: string;
    photo: string;
    phone: string;
    is_admin: boolean;
}>;

export type OauthAuthorizationParams = {
    client_id: string,
    redirect_uri: string,
    response_type: string,
    scope: string,
    state: string,
    code_challenge?: string | null,
    code_challenge_method?: string | null,
}

export type TokenExchangeRequest = {
    grant_type: "authorization_code",
    code: string,
    client_id: string,
    redirect_uri: string,
    code_verifier: string,
}

export type TokenExchangeResponse = BaseResponse<{
    access_token: string;
    refresh_token: string;
    expires_in: number;
    token_type: string;
}>


export type UpdateAccountRequest = {
    full_name?: string;
    email?: string;
    username?: string;
    photo?: string | File;
    phone?: string;
    password?: string;
    password_confirmation?: string;
};

export type ValidateAppRequest = {
    client_id: string;
    redirect_uri: string;
}

export type ValidateAppResponse = BaseResponse<{
    client_id: string;
    client_name: string;
    logo_url: string;
    redirect_uri: string;
    scope: string;
}>

export type JWTPayload = {
    user_id: string;
    full_name: string;
    email: string;
    username: string;
    phone: string;
    photo: string;
    service_name: string;
    role_name: string;
}

export type JWTClaimResponse = BaseResponse<JWTPayload>;


// Forgot & Reset Password
export type ForgotUsernameorPasswordRequest = {
    email: string;
    forgot_type: "username" | "password";
}

export type ResetPasswordrequest = {
    token: string;
    password: string;
    password_confirmation: string;
}