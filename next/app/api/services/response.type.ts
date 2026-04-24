export type BaseResponse<T> = {
    success: boolean;
    message: string;
    data: T;
}

export type ErrorResponse = {
    success: boolean;
    message: string;
    error: string;
}

export type PaginationResponse<T> = BaseResponse<{
    rows: T[];
    limit: number;
    page:number;
    sort: string;
    total_rows: number;
    total_pages: number;
}>