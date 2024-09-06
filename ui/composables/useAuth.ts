import { useStore } from "./useStore";
import type { LoginResponse } from "~/types"

export const useAuth = () => {
    const user = useStore<LoginResponse>('auth', {} as LoginResponse)
    return {
        user
    }
}