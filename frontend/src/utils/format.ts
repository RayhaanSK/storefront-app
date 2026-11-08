export function formatPrice(price: number): string {
    return `$${price.toFixed(2)}`
}

// Turns an ISO timestamp from the API into a date in the browser's own locale format.
export function formatDate(isoTimestamp: string): string {
    const date = new Date(isoTimestamp)
    return Number.isNaN(date.getTime()) ? '' : date.toLocaleDateString()
}