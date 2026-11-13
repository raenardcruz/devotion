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
import KindleReadings from '../views/KindleReadings.vue'
import Resources from '../views/Resources.vue'
import BibleStudy from '../views/BibleStudy.vue'
import Admin from '../views/Admin.vue'
import { useMagisteriumModal } from '../composables/useMagisteriumModal'
import { useLastRead } from '../composables/useLastRead'
import { useAccessControl } from '../composables/useAccessControl'

const router = createRouter({
    history: createWebHashHistory(),
    scrollBehavior(to, _from, savedPosition) {
        if (savedPosition) {
            return savedPosition
        }
        if (to.path === '/kindle') {
            return { top: 0 }
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
            component: Devotions,
            meta: { lastReadTitle: 'Devotions' }
        },
        {
            path: '/daily-readings',
            name: 'DailyMassReadings',
            component: DailyMassReadings,
            meta: { lastReadTitle: 'Daily Mass Readings' }
        },
        {
            path: '/kindle',
            name: 'KindleReadings',
            component: KindleReadings,
            meta: { lastReadTitle: 'Kindle Readings' }
        },
        {
            path: '/rosary',
            name: 'Rosary',
            component: Rosary,
            meta: { lastReadTitle: 'The Holy Rosary' }
        },
        {
            path: '/seven-sorrows',
            name: 'SevenSorrows',
            component: SevenSorrows,
            meta: { lastReadTitle: 'Seven Sorrows of Mary' }
        },
        {
            path: '/st-michael-rosary',
            name: 'StMichaelRosary',
            component: StMichaelRosary,
            meta: { lastReadTitle: 'Chaplet of St. Michael' }
        },
        {
            path: '/divine-mercy',
            name: 'DivineMercy',
            component: DivineMercy,
            meta: { lastReadTitle: 'Divine Mercy Chaplet' }
        },
        {
            path: '/prayers',
            name: 'PrayersGallery',
            component: PrayersGallery,
            meta: { lastReadTitle: 'Prayers Gallery' }
        },
        {
            path: '/catechism',
            name: 'Catechism',
            component: Catechism,
            meta: { lastReadTitle: 'Catechism of the Catholic Church' }
        },
        {
            path: '/resources',
            name: 'Resources',
            component: Resources,
            meta: { lastReadTitle: 'Early Church Resources' }
        },
        {
            path: '/bible-study',
            name: 'BibleStudy',
            component: BibleStudy,
            meta: { lastReadTitle: 'Bible Study' }
        },
        {
            path: '/magisterium-chat/:id?',
            redirect: (to) => {
                const { openMagisteriumModal } = useMagisteriumModal()
                const id = typeof to.params.id === 'string' && to.params.id ? to.params.id : undefined
                openMagisteriumModal(id)
                return '/'
            }
        },
        {
            path: '/admin',
            name: 'Admin',
            component: Admin
        }
    ]
})

router.beforeEach(async (to) => {
    const { acceptInviteCode } = useAccessControl()
    const hashCode = typeof to.query.access_code === 'string' ? to.query.access_code : ''
    const searchParams = typeof window !== 'undefined' ? new URLSearchParams(window.location.search) : null
    const searchCode = searchParams?.get('access_code') || ''
    const inviteCode = hashCode || searchCode

    if (inviteCode) {
        await acceptInviteCode(inviteCode)

        if (searchCode && typeof window !== 'undefined' && searchParams) {
            searchParams.delete('access_code')
            const newSearch = searchParams.toString()
            const cleanUrl = `${window.location.pathname}${newSearch ? `?${newSearch}` : ''}${window.location.hash}`
            window.history.replaceState({}, '', cleanUrl)
        }

        if (hashCode) {
            const cleanedQuery = { ...to.query }
            delete cleanedQuery.access_code
            return { path: to.path, query: cleanedQuery, hash: to.hash, replace: true }
        }
    }
})

router.afterEach((to) => {
    const title = to.meta?.lastReadTitle
    if (typeof title === 'string' && title.length > 0) {
        const { saveLastRead } = useLastRead()
        saveLastRead({
            title,
            path: to.fullPath,
        })
    }
})

export default router