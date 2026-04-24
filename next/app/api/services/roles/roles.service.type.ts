import { BaseResponse, PaginationResponse } from "../response.type"

export type Role = {
    service_role_id: number;
    service_id: number;
    role_name: string;
    service_name: string;
    created_at: string;
    updated_at: string;
}

export type CreateRoleRequest = {
    service_id: number;
    role_name: string;
}

export type CreateRoleResponse = BaseResponse<Role>
export type GetRolesResponse = PaginationResponse<Role>