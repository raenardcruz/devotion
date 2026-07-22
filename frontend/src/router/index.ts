import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Rosary from '../views/Rosary.vue'
import DivineMercy from '../views/DivineMercy.vue'
import PrayersGallery from '../views/PrayersGallery.vue'
import Catechism from '../views/Catechism.vue'
import DailyMassReadings from '../views/DailyMassReadings.vue'
import Resources from '../views/Resources.vue'
import BibleStudy from '../views/BibleStudy.vue'

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: '/daily-readings',
            name: 'DailyMassReadings',
            component: DailyMassReadings
        },
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
        }
    ]
})

export default router