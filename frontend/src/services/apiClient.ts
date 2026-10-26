const API_BASE_URL = 'http://localhost:8080'

export async function apiRequest<T>(
    path: string,
    options: RequestInit = {},
): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${path}`, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            ...options.headers,
        },
    })

    if (!response.ok) {
        let message = 'Request failed'

        try {
            const body = await response.json()

            if (typeof body.error === 'string') {
                message = body.error
            }
        } catch {
            // The response might not contain JSON.
        }

        throw new Error(message)
    }

    if (response.status === 204) {
        return undefined as T
    }

    return response.json() as Promise<T>
}