import { BaseResponse, PaginationResponse } from "../response.type";

export type User = {
    user_id: number;
    full_name: string;
    email: string;
    username: string;
    photo: File|string;
    phone: string;
    password: string;
    role_name: string | null;
    created_at: Date;
    updated_at: Date;
}

export type UserRequest = User
export type UserResponse = BaseResponse<User>
export type UsersListResponse = PaginationResponse<User>