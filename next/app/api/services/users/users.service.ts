import apiClient from "../../api-client";
import { PaginationRequest } from "../request.type";
import { UserRequest, UserResponse, UsersListResponse } from "./users.service.type";

export async function getUsers(params: PaginationRequest): Promise<UsersListResponse> {
   try {
      const response = await apiClient.get("/users", { params });
      return response.data;
   } catch (error) {
      throw error;
   }
}

export async function createUser(user: UserRequest): Promise<UserResponse> {
   try {
      const formData = new FormData();
      Object.entries(user).forEach(([key, value]) => {
         if (value instanceof Date) {
            formData.append(key, value.toISOString());
         } else if (value instanceof File) {
            formData.append(key, value);
         } else if (value !== null && value !== undefined) {
            formData.append(key, String(value));
         }
      });

      const response = await apiClient.post("/users", formData, {
         headers: {
            "Content-Type": "multipart/form-data",
         },
      });
      return response.data;
   } catch (error) {
      throw error;
   }
}

export async function updateUser(userId: string, user: UserRequest): Promise<UserResponse> {
   try {
      const formData = new FormData();
      Object.entries(user).forEach(([key, value]) => {
         if (value instanceof Date) {
            formData.append(key, value.toISOString());
         } else if (value instanceof File) {
            formData.append(key, value);
         } else if (value !== null && value !== undefined) {
             // Avoid sending empty string or string 'null' for photo if it's not a file
            if (key === 'photo' && typeof value === 'string') return;
            formData.append(key, String(value));
         }
      });

      const response = await apiClient.put(`/users/${userId}`, formData, {
         headers: {
            "Content-Type": "multipart/form-data",
         },
      });
      return response.data;
   } catch (error) {
      throw error;
   }
}

export async function deleteUser(userId: string): Promise<UserResponse> {
   try {
      const response = await apiClient.delete(`/users/${userId}`);
      return response.data;
   } catch (error) {
      throw error;
   }
}