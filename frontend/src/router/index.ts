import { createRouter, createWebHistory } from 'vue-router'

import { useSession } from '../composables/useSession'
import HomeView from '../views/HomeView.vue'

// Lets us put `meta: { guestOnly: true }` on a route with type checking.
declare module 'vue-router' {
    interface RouteMeta {
        guestOnly?: boolean
    }
}

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            name: 'home',
            component: HomeView,
        },
        {
            path: '/catalogue',
            name: 'catalogue',
            // Lazy-loaded: this page's code is only downloaded when someone visits it.
            component: () => import('../views/CatalogueView.vue'),
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
        // Later: /product/:id, /cart, /my-listings, /users ...
    ],
})

// Runs before every navigation. Someone who is already logged in has no reason to see
// the login or signup pages, so send them home instead.
router.beforeEach((to) => {
    const { isLoggedIn } = useSession()

    if (to.meta.guestOnly && isLoggedIn.value) {
        return { name: 'home' }
    }
})

export default router