import { BaseResponse } from "../response.type";

export type Client = {
    service_id: number;
    service_name: string;
    description: string;
    logo: string;
    client_id: string;
    client_secret: string;
    redirect_url: string;
    is_active:  boolean;
    created_at: Date;
    updated_at: Date;
}
export type CreateClientRequest = {
    name: string;
    description: string;
    logo: string | File;
    redirect_url: string;
}

export type CreateClientResponse = BaseResponse<Client>;