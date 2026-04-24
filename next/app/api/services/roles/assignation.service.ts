import apiClient from "../../api-client";
import { Assignation, AssignationCreateRequest, AssignationDeleteRequest, AssignationGetRequest, AssignationResponse } from "./assignation.service.type";

export async function getAssignations(params: AssignationGetRequest): Promise<AssignationResponse> {
    try {
        const response = await apiClient.get("roles/assign", {
            params
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}

export async function createAssignation(assignation: AssignationCreateRequest): Promise<Assignation> {
    try {
        // Correcting the likely typo in user request to match the standard resource path pattern
        const response = await apiClient.post(`/roles/assign`, assignation);
        return response.data;
    } catch (error) {
        throw error;
    }
}

export async function deleteAssignation(assignation: AssignationDeleteRequest): Promise<Assignation> {
    try {
        const response = await apiClient.delete(`/roles/assign`, {
            data: assignation
        });
        return response.data;
    } catch (error) {
        throw error;
    }
}
