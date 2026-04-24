import apiClient from "../../api-client";
import { OAUTH_AUTHORIZE_URL } from "../../config";
import { BaseResponse } from "../response.type";
import { UserAccess } from "../user-access/user-access.service.type";
import { LoginRequest, LoginResponse, TokenExchangeRequest, OauthAuthorizationParams, TokenExchangeResponse, UpdateAccountRequest, JWTClaimResponse, ForgotUsernameorPasswordRequest, ResetPasswordrequest } from "./oauth.type";


export async function oauthLogin(data: LoginRequest, params: OauthAuthorizationParams|null): Promise<LoginResponse> {
    try {
        const response = await apiClient.post("/oauth/login", data, {params});
        return response.data;
    } catch (error) {
        throw error;
    }
}

export async function oauthTokenExchange(data: TokenExchangeRequest): Promise<TokenExchangeResponse>{
    try {
        const response = await apiClient.post("/oauth/token", data);
        return response.data;
    } catch(error) {
        throw error;
    }   
}

export async function oauthGetMe(): Promise<JWTClaimResponse>{
    try {
        const response = await apiClient.get("/oauth/me");
        return response.data;
    } catch (error) {
        throw error;
    }
}

export async function oauthLogout() {
    try {
        const response = await apiClient.post("/oauth/logout", {});
        return response.data;
    } catch (error) {
        throw error;
    }
}

export async function updateAccountRequest(data: UpdateAccountRequest): Promise<BaseResponse<TokenExchangeResponse>> {
    try {
        const formData = new FormData();
        if (data.full_name) formData.append("full_name", data.full_name);
        if (data.email) formData.append("email", data.email);
        if (data.username) formData.append("username", data.username);
        if (data.phone) formData.append("phone", data.phone);
        if (data.password) formData.append("password", data.password);
        if (data.password_confirmation) formData.append("password_confirmation", data.password_confirmation);
        if (data.photo) formData.append("photo", data.photo);

        const response = await apiClient.put("/oauth/update-account", formData, {
            headers: { "Content-Type": "multipart/form-data" },
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}

export async function forgotUsernameorPasswordRequest(data: ForgotUsernameorPasswordRequest): Promise<BaseResponse<{token: string}>> {
    try {
        const response = await apiClient.post("/oauth/forgot-password", data);
        return response.data;
    } catch (error) {
        throw error;
    }
}

export async function resetPasswordRequest(data: ResetPasswordrequest): Promise<BaseResponse<{token: string}>> {
    try {
        const response = await apiClient.post("/oauth/reset-password", data);
        return response.data;
    } catch (error) {
        throw error;
    }
}


// Helper
export function generateAuthorizationURL(params: Partial<OauthAuthorizationParams>): string{
    const queryParams = new URLSearchParams();

    if (params.client_id) queryParams.append("client_id", params.client_id);
    if (params.redirect_uri) queryParams.append("redirect_uri", params.redirect_uri);
    if (params.response_type) queryParams.append("response_type", params.response_type);
    if (params.scope) queryParams.append("scope", params.scope);
    if (params.state) queryParams.append("state", params.state);
    if (params.code_challenge) queryParams.append("code_challenge", params.code_challenge);
    if (params.code_challenge_method) queryParams.append("code_challenge_method", params.code_challenge_method);

    return OAUTH_AUTHORIZE_URL + "?" + queryParams.toString()
}

export function generateAppAuthorizationURL(app: UserAccess): string {
    const queryParams = new URLSearchParams();

    if (app.client_id) queryParams.append("client_id", app.client_id);
    if (app.redirect_url) queryParams.append("redirect_uri", app.redirect_url);
    queryParams.append("response_type", "code");
    queryParams.append("scope", "openid profile email");
    queryParams.append("state", "state");

    return OAUTH_AUTHORIZE_URL + "?" + queryParams.toString()
}