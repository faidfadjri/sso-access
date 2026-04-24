import apiClient from "../../api-client";
import { PaginationRequest } from "../request.type";
import { BaseResponse } from "../response.type";
import { CreateUserAccessRequest, CreateUserAccessResponse, GetUserAccessResponse } from "./user-access.service.type";

export async function getUserAccess(params: PaginationRequest): Promise<GetUserAccessResponse> {
    try {
        const response = await apiClient.get("/users/access", {
            params
        });
        return response.data;
    } catch(error) {
        throw(error)
    }
}

export async function createUserAccess(data: CreateUserAccessRequest): Promise<CreateUserAccessResponse> {
    try {
        const response = await apiClient.post("/users/access", data);
        return response.data;
    } catch(error) {
        throw(error)
    }
}

export async function deleteUserAccess(accessId: number): Promise<BaseResponse<null>> {
    try {
        const response = await apiClient.delete(`/users/access/${accessId}`);
        return response.data;
    } catch(error) {
        throw(error)
    }
}