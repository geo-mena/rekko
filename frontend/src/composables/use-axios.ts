import type { AxiosError } from 'axios'

import axios from 'axios'

import env from '@/utils/env'

export function useAxios() {
    const axiosInstance = axios.create({
        baseURL: env.VITE_SERVER_API_URL + env.VITE_SERVER_API_PREFIX,
        timeout: env.VITE_SERVER_API_TIMEOUT,
    })

    axiosInstance.interceptors.request.use((config) => {
        return config
    }, (error) => {
        return Promise.reject(error)
    })

    axiosInstance.interceptors.response.use((response) => {
        return response
    }, (error: AxiosError) => {
        return Promise.reject(error)
    })

    return {
        axiosInstance,
    }
}
