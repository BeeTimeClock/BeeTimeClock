import {RouteRecordRaw} from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('pages/LoginPage.vue'),
  },
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    redirect: { name: 'Dashboard' },
    children: [
      {
        name: 'Dashboard',
        path: 'dashboard',
        component: () => import('pages/DashboardPage.vue'),
        alias: '/code',
      },
      {
        name: 'WorktimeOverview',
        path: 'worktime',
        component: () => import('pages/worktime/WorktimeOverviewPage.vue')
      },
      {
        name: 'AbsenceOverview',
        path: 'absence',
        component: () => import('pages/absence/AbsenceOverviewPage.vue'),
      },
      {
        name: 'FuelOverview',
        path: 'fuel',
        component: () => import('pages/fuel/FuelOverviewPage.vue')
      },
      {
        path: 'me',
        children: [
          {
            name: 'UserApikeyOverview',
            path: 'apikey',
            component: () => import('pages/user/UserApikeyPage.vue')
          }
        ]
      },
      {
        path: 'administration',
        children: [
          {
            name: 'AdministrationUserOverview',
            path: 'user',
            component: () => import('pages/administration/AdministrationUserOverviewPage.vue')
          },
          {
            name: 'AdministrationUserDetail',
            path: 'user/:userId',
            component: () => import('pages/administration/AdministrationUserDetailPage.vue')
          },
          {
            name: 'AdministrationSettings',
            path: 'settings',
            component: () => import('pages/administration/AdministrationSettingsPage.vue')
          }
        ]
      }
    ],
  },
  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
];

export default routes;
