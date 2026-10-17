import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Devotions from '../views/Devotions.vue'
import Rosary from '../views/Rosary.vue'
import SevenSorrows from '../views/SevenSorrows.vue'
import StMichaelRosary from '../views/StMichaelRosary.vue'
import DivineMercy from '../views/DivineMercy.vue'
import PrayersGallery from '../views/PrayersGallery.vue'
import Catechism from '../views/Catechism.vue'
import DailyMassReadings from '../views/DailyMassReadings.vue'
import Resources from '../views/Resources.vue'
import BibleStudy from '../views/BibleStudy.vue'
import MagisteriumChat from '../views/MagisteriumChat.vue'
import Admin from '../views/Admin.vue'

const router = createRouter({
    history: createWebHashHistory(),
    scrollBehavior(_to, _from, savedPosition) {
        if (savedPosition) {
            return savedPosition
        }
        return { top: 0, behavior: 'smooth' }
    },
    routes: [
        {
            path: '/',
            name: 'Home',
            component: Home
        },
        {
            path: '/devotions',
            name: 'Devotions',
            component: Devotions
        },
        {
            path: '/daily-readings',
            name: 'DailyMassReadings',
            component: DailyMassReadings
        },
        {
            path: '/rosary',
            name: 'Rosary',
            component: Rosary
        },
        {
            path: '/seven-sorrows',
            name: 'SevenSorrows',
            component: SevenSorrows
        },
        {
            path: '/st-michael-rosary',
            name: 'StMichaelRosary',
            component: StMichaelRosary
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
        },
        {
            path: '/resources',
            name: 'Resources',
            component: Resources
        },
        {
            path: '/bible-study',
            name: 'BibleStudy',
            component: BibleStudy
        },
        {
            path: '/magisterium-chat/:id?',
            name: 'MagisteriumChat',
            component: MagisteriumChat
        },
        {
            path: '/admin',
            name: 'Admin',
            component: Admin
        }
    ]
})

export default router