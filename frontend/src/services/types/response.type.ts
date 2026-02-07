export interface ApiResponse<T = undefined> {
    status: boolean
    message: string
    data: T
}

export interface PaginationMeta {
    current_page: number
    per_page: number
    total_items: number
    total_pages: number
    has_next: boolean
}

export interface ApiPaginatedResponse<T> extends ApiResponse<T> {
    meta: {
        pagination: PaginationMeta
    }
}

/** @deprecated Use ApiResponse instead */
export interface IResponse<T, E = Record<string, any>> {
    data: T
    extra: E
    code: number
    message: string
    success: boolean
}
