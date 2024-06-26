import { defineStore } from 'pinia';

// 储存用例：
// export const useCounterStore = defineStore('counter', {
//   state: () => ({
//     counter: 0,
//   }),
//   getters: {
//     doubleCount: (state) => state.counter * 2,
//   },
//   actions: {
//     increment() {
//       this.counter++;
//     },
//   },
// });

// 使用用例：
// import { apiUrlStorage } from './../stores/api-urls';
// const urlStore = apiUrlStorage();
// const SOME_URL = urlStore.XXXXXXX;

export const apiUrlStorage = defineStore('apiUrls', {
  state: () => ({
    // // 遗产API, 从上一个工程里暴力移植过来的，没有意义，仅供参考
    // getFolder: 'http://0.0.0.0:8000/api/get_folder',
    // getFilePreview: 'http://0.0.0.0:8000/api/get_file_preview',
    // getMediaLibraries: 'http://0.0.0.0:8000/api/get_media_libraries',
    // getMediaLibraryContent: 'http://0.0.0.0:8000/api/get_media_library_content_by_id',
    // getMediaUnit: 'http://0.0.0.0:8000/api/get_media_unit_by_id',
    // getTmdbTvSeriesMetadataByUnit: 'http://0.0.0.0:8080/api/get_tmdb_tv_series_metadata_by_unit',
    // getTmdbTvEpisodeMetadataByUnit: 'http://0.0.0.0:8080/api/get_tmdb_tv_episode_metadata_by_unit',

    // 图书馆API相关
    // 接受参数：page:int [0, inf], page_size:int [1, inf]
    bookLibraryGetBooksShortInfo: 'http://127.0.0.1:8080/api/book_library/get_books_short_info',
    bookLibraryGetBooksInfo: 'http://127.0.0.1:8080/api/book_library/get_books_info',

    bookLibraryGetBookInfoById: 'http://127.0.0.1:8080/api/book_library/get_book_info_by_id',
    bookLibraryGetServedBookfileById: 'http://127.0.0.1:8080/api/book_library/serve_book_file_by_id',

    BookLibraryGetBookFirstPagePdf: 'http://127.0.0.1:8080/api/book_library/get_book_first_page_pdf',
  }),
});
