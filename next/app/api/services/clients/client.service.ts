import apiClient from "../../api-client";
import { PaginationRequest } from "../request.type";
import { PaginationResponse } from "../response.type";
import { Client, CreateClientResponse } from "./client.type";

export async function getServiceClients(params: PaginationRequest): Promise<PaginationResponse<Client>> {
    const response = await apiClient.get('/service/clients', {
        params
    });
    return response.data;
}

export async function createServiceClient(client: FormData): Promise<CreateClientResponse> {
    const response = await apiClient.post('/service/clients', client, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
    return response.data;
}

export async function updateServiceClient(serviceId: number, client: FormData): Promise<void> {
    const response = await apiClient.put(`/service/clients/${serviceId}`, client, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
    return response.data;
}

export async function deleteServiceClient(serviceId: number): Promise<void> {
    await apiClient.delete(`/service/clients/${serviceId}`);
}