import { BaseResponse, PaginationResponse } from "../response.type";

export type Assignation = {
    user_id: number;
    service_role_id: number;
    service_id: number;
    full_name: string;
    role_name: string;
    service_name: string;
    created_at: Date;
    updated_at: Date;
}

export type AssignationResponse = PaginationResponse<Assignation>

export type AssignationCreateRequest = {
    user_ids: number[];
    role_id: number;
    service_id: number;
};

export type AssignationDeleteRequest = {
    user_id: number;
    role_id: number;
    service_id: number;
};

export type AssignationGetRequest = {
    service_id: number;
    role_id?: number;
    show: number
};