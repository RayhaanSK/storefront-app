import { createRouter, createWebHistory } from 'vue-router'

import { useSession, type Role } from '../composables/useSession'
import HomeView from '../views/HomeView.vue'

// Lets us put these on a route's `meta` with type checking.
declare module 'vue-router' {
    interface RouteMeta {
        guestOnly?: boolean // only for people who are logged out (login, signup)
        role?: Role // only for people logged in as this type of user
    }
}

const router = createRouter({
    history: createWebHistory(),
    routes: [
        { path: '/', name: 'home', component: HomeView },
        {
            path: '/catalogue',
            name: 'catalogue',
            // Lazy-loaded: this page's code is only downloaded when someone visits it.
            component: () => import('../views/CatalogueView.vue'),
        },
        {
            path: '/product/:id',
            name: 'product-details',
            component: () => import('../views/ProductDetailsView.vue'),
        },
        {
            path: '/login',
            name: 'login',
            component: () => import('../views/LoginView.vue'),
            meta: { guestOnly: true },
        },
        {
            path: '/signup',
            name: 'signup',
            component: () => import('../views/SignUpView.vue'),
            meta: { guestOnly: true },
        },

        // Customer
        {
            path: '/cart',
            name: 'cart',
            component: () => import('../views/CartView.vue'),
            meta: { role: 'customer' },
        },

        // Seller
        {
            path: '/my-listings',
            name: 'my-listings',
            component: () => import('../views/MyListingsView.vue'),
            meta: { role: 'seller' },
        },
        {
            path: '/my-listings/new',
            name: 'new-listing',
            component: () => import('../views/NewListingView.vue'),
            meta: { role: 'seller' },
        },
        {
            path: '/my-listings/:id/edit',
            name: 'edit-listing',
            component: () => import('../views/EditListingView.vue'),
            meta: { role: 'seller' },
        },

        // Admin
        {
            path: '/users',
            name: 'users',
            component: () => import('../views/UsersView.vue'),
            meta: { role: 'admin' },
        },
    ],
})

// Runs before every navigation.
router.beforeEach((to) => {
    const { isLoggedIn, role } = useSession()

    // Someone who is already logged in has no reason to see the login or signup pages.
    if (to.meta.guestOnly && isLoggedIn.value) {
        return { name: 'home' }
    }

    if (to.meta.role) {
        if (!isLoggedIn.value) return { name: 'login' }
        if (role.value !== to.meta.role) return { name: 'home' }
    }
})

export default router