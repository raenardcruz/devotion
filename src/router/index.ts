import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Rosary from '../views/Rosary.vue'
import DivineMercy from '../views/DivineMercy.vue'
import PrayersGallery from '../views/PrayersGallery.vue'

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            name: 'Home',
            component: Home
        },
        {
            path: '/rosary',
            name: 'Rosary',
            component: Rosary
        },
        {
            path: '/divine-mercy',
            name: 'DivineMercy',
            component: DivineMercy
        },
        {
            path: '/prayers',
            name: 'PrayersGallery',
            component: PrayersGallery
        }
    ]
})

export default router