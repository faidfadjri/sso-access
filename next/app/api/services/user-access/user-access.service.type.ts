import { BaseResponse, PaginationResponse } from "../response.type";

export type UserAccess = {
    access_id: number;
    user_id: number;
    service_id: number;
    status: "active" | "revoke";
    username: string;
    full_name:string;
    email: string;
    phone: string;
    service_name:string;
    logo?: string;
    redirect_url: string;
    client_id: string;
    created_at: Date;
    updated_at: Date;
}

export type GetUserAccessResponse = PaginationResponse<UserAccess>

export type CreateUserAccessRequest = {
    user_id: number;
    service_ids: number[];
    status: "active" | "revoke";
}

export type CreateUserAccessResponse = BaseResponse<UserAccess>