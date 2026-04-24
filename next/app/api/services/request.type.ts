export type PaginationRequest = {
    show?: number;
    page?: number;
    sort?: 'asc'|'desc';
    search?: string;
    user_id?: string;
    service_id?: string;
    role_id?: string;
}