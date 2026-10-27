// Returns the SHA-256 hash of a password as a lowercase hex string.
// The API fields are named "passwordHash", so the browser sends this instead of the raw password.
// crypto.subtle only exists on secure origins (https:// or http://localhost).
export async function hashPassword(password: string): Promise<string> {
    const bytes = new TextEncoder().encode(password)
    const digest = await crypto.subtle.digest('SHA-256', bytes)

    return Array.from(new Uint8Array(digest))
        .map((byte) => byte.toString(16).padStart(2, '0'))
        .join('')
}