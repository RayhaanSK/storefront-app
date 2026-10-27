import { computed, ref } from 'vue'

const STORAGE_KEY = 'storefront.loggedInUser'

export const ROLES = ['customer', 'seller', 'admin'] as const
export type Role = (typeof ROLES)[number]

// Only what the UI needs. The password hash is deliberately never stored.
export interface LoggedInUser {
    id: number
    username: string
    userType: string
    userTypeId: number
}

function readStoredUser(): LoggedInUser | null {
    try {
        const raw = localStorage.getItem(STORAGE_KEY)
        if (raw === null) return null

        const parsed: unknown = JSON.parse(raw)
        if (typeof parsed !== 'object' || parsed === null) return null

        const value = parsed as Record<string, unknown>
        if (
            typeof value.id !== 'number' ||
            typeof value.username !== 'string' ||
            typeof value.userType !== 'string' ||
            typeof value.userTypeId !== 'number'
        ) {
            return null
        }

        return {
            id: value.id,
            username: value.username,
            userType: value.userType,
            userTypeId: value.userTypeId,
        }
    } catch {
        // Storage unavailable or the stored text is not valid JSON: treat as logged out.
        return null
    }
}

// Module-level state: every component that calls useSession() shares this one object,
// so logging in or out updates the navbar and pages immediately.
const user = ref<LoggedInUser | null>(readStoredUser())

// Keep other open tabs in sync when this tab logs in or out.
window.addEventListener('storage', (event) => {
    if (event.key === STORAGE_KEY || event.key === null) {
        user.value = readStoredUser()
    }
})

const currentUser = computed(() => user.value)
const isLoggedIn = computed(() => user.value !== null)

const role = computed<Role | null>(() => {
    const userType = user.value?.userType.toLowerCase()
    return ROLES.find((candidate) => candidate === userType) ?? null
})

function login(newUser: LoggedInUser) {
    user.value = {
        id: newUser.id,
        username: newUser.username,
        userType: newUser.userType,
        userTypeId: newUser.userTypeId,
    }

    try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(user.value))
    } catch (error) {
        console.error('Could not save login to local storage', error)
    }
}

function logout() {
    user.value = null

    try {
        localStorage.removeItem(STORAGE_KEY)
    } catch (error) {
        console.error('Could not clear login from local storage', error)
    }
}

export function useSession() {
    return { currentUser, isLoggedIn, role, login, logout }
}