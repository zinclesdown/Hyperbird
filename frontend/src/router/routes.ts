import { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('layouts/StatusPanelLayout.vue'),
    children: [
      {
        path: '',
        component: () => import('pages/IndexPage.vue'),
      },
      {
        path: 'home',
        component: () => import('pages/HomePage.vue'),
      },
      {
        path: 'book_library',
        component: () => import('pages/BookLibrary.vue'),
      },
      {
        path: 'book_library/info',
        component: () => import('pages/BookInfo.vue'),
      },
      {
        path: 'book_library/pdf_viewer',
        component: () => import('pages/PdfViewer.vue'),
      },
      {
        path: 'book_library/manager',
        component: () => import('pages/BookLibraryManager.vue'),
      },
      // {
      //   path: 'library_browser',
      //   component: () => import('pages/LibraryBrowser.vue'),
      // },
      {
        path: 'about',
        component: () => import('pages/AboutPage.vue'),
      },
      // {
      //   path: 'file_browser',
      //   component: () => import('pages/FileBrowserPage.vue'),
      // },
      {
        path: 'debug_page',
        component: () => import('pages/DebugPage.vue'),
      },
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
