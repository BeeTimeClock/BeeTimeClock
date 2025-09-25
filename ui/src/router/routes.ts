import {type RouteRecordRaw} from 'vue-router';

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
        name: 'SuspiciousTimestampsOverview',
        path: 'worktime/suspicious',
        component: () => import('pages/worktime/SuspiciousTimestampsOverviewPage.vue')
      },
      {
        name: 'MissingDaysOverview',
        path: 'worktime/missing',
        component: () => import('pages/worktime/MissingDaysOverviewPage.vue')
      },
      {
        name: 'AbsenceOverview',
        path: 'absence',
        component: () => import('pages/absence/AbsenceOverviewPage.vue'),
      },
      {
        name: 'ExternalWorkOverview',
        path: 'external_work',
        component: () => import('pages/externalWork/ExternalWorkOverviewPage.vue')
      },
      {
        name: 'ExternalWorkDetail',
        path: 'external_work/:externalWorkId',
        component: () => import('pages/externalWork/ExternalWorkDetailPage.vue')
      },
      {
        name: 'FuelOverview',
        path: 'fuel',
        component: () => import('pages/fuel/FuelOverviewPage.vue')
      },
      {
        name: 'OvertimeOverview',
        path: 'overtime',
        component: () => import('pages/overtime/OvertimeOverviewPage.vue')
      },
      {
        path: 'me',
        children: [
          {
            name: 'UserApikeyOverview',
            path: 'apikey',
            component: () => import('pages/user/UserApikeyPage.vue')
          },
          {
            name: 'UserSettings',
            path: 'settings',
            component: () => import('pages/user/UserSettingsPage.vue')
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
            name: 'AdministrationTeamOverview',
            path: 'team',
            component: () => import('pages/administration/AdministrationTeamOverviewPage.vue')
          },
          {
            name: 'AdministrationTeamDetail',
            path: 'team/:teamId',
            component: () => import('pages/administration/AdministrationTeamDetailPage.vue')
          },
          {
            path: 'settings',
            children: [
              {
                name: 'AdministrationSettingsCommon',
                path: 'common',
                component: () => import('pages/administration/settings/AdministrationSettingsCommonPage.vue')
              },
              {
                name: 'AdministrationSettingsAbsence',
                path: 'absence',
                component: () => import('pages/administration/settings/AdministrationSettingsAbsencePage.vue')
              },
              {
                name: 'AdministrationSettingsTimestamp',
                path: 'timestamp',
                component: () => import('pages/administration/settings/AdministrationSettingsTimestampPage.vue')
              },
              {
                name: 'AdministrationSettingsNotify',
                path: 'notify',
                component: () => import('pages/administration/settings/AdministrationSettingsNotifyPage.vue')
              },
              {
                name: 'AdministrationSettingsExternalWork',
                path: 'external_work',
                component: () => import('pages/administration/settings/AdministrationSettingsExternalWorkPage.vue')
              }
            ],
          },
          {
            path: 'debug',
            component: () => import('pages/DebugPage.vue')
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
