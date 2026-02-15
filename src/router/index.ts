import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Rosary from '../views/Rosary.vue'
import DivineMercy from '../views/DivineMercy.vue'
import PrayersGallery from '../views/PrayersGallery.vue'
import Catechism from '../views/Catechism.vue'

const router = createRouter({
    history: createWebHashHistory(),
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
        },
        {
            path: '/catechism',
            name: 'Catechism',
            component: Catechism
        }
    ]
})

export default router