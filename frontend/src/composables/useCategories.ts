import { ref } from 'vue'

import { listAllCategories } from '../services/storefront'
import type { Category } from '../models/storefront'

// Loads the list of categories for the New Listing and Edit Listing forms.
// Never throws: on failure it fills in categoriesError instead.
export function useCategories() {
    const categories = ref<Category[]>([])
    const categoriesError = ref('')

    async function loadCategories() {
        categoriesError.value = ''

        try {
            // The backend can send null instead of [] when there are no rows.
            categories.value = (await listAllCategories()) ?? []
        } catch (error) {
            console.error(error)
            categoriesError.value = 'Couldn’t load the categories. Check that the server is running.'
        }
    }

    return { categories, categoriesError, loadCategories }
}