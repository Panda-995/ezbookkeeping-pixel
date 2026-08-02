import type { ApiResponse } from './api.ts';

export function isSuccessfulApiResponse<T>(response: ApiResponse<T> | null | undefined): response is ApiResponse<T> {
    return !!response && response.success;
}
