import apiClient from "../../api-client";
import { PaginationRequest } from "../request.type";
import { CreateRoleRequest, CreateRoleResponse, GetRolesResponse } from "./roles.service.type";

export async function getRoles(params: PaginationRequest): Promise<GetRolesResponse>{
    try {
        const response = await apiClient.get("/roles", {params});
        return response.data;
    } catch(error) {
        throw(error)
    }
}

export async function createRole(data: CreateRoleRequest): Promise<CreateRoleResponse>{
    try {
        const response = await apiClient.post("/roles", data);
        return response.data;
    } catch(error) {
        throw(error)
    }
}

export async function updateRole(roleId: number, data: CreateRoleRequest): Promise<CreateRoleResponse>{
    try {
        const response = await apiClient.put(`/roles/${roleId}`, data);
        return response.data;
    } catch(error) {
        throw(error)
    }
}

export async function deleteRole(roleId: number): Promise<CreateRoleResponse>{
    try {
        const response = await apiClient.delete(`/roles/${roleId}`);
        return response.data;
    } catch(error) {
        throw(error)
    }
}