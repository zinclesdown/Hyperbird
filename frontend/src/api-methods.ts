// import apiUrls from './apiUrls';

// export async function getTmdbTvSeriesMetadata(unit_id: number | undefined) {
//   const response = await axios.get(apiUrls.getTmdbTvSeriesMetadataByUnit, {
//     params: {
//       unit_id: unit_id,
//     },
//   });
//   return response.data;
// }

// // 获取某个Episode的 TMDB 信息。需要提供 unit_id, season, episode
// export async function getTmdbTvEpisodeMetadata(
//   unit_id: number | undefined,
//   season: number | undefined,
//   episode: number | undefined,
// ) {
//   const response = await axios.get(apiUrls.getTmdbTvEpisodeMetadataByUnit, {
//     params: {
//       unit_id: unit_id,
//       season: season,
//       episode: episode,
//     },
//   });
//   return response.data;
// }

// export default {
//   getTmdbTvSeriesMetadata,
//   getTmdbTvEpisodeMetadata,
// };

// 定义接口

import { ref } from 'vue';
import { apiUrlStorage } from './stores/api-urls';
import axios from 'axios';
const urlStore = apiUrlStorage(); // Pinia, 读取urlStore的API地址

export interface BookInfo {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string;
  bookid: string;
  bookname: string;
  bookimagepath: string;
}

export interface BookInfoResponse {
  book: BookInfo;
}

export interface BookShortInfo {
  book_id: string;
  book_name: string;
  book_image_path: string;
}

export interface BookShortInfoResponse {
  books: BookShortInfo[];
}

// 获取书籍信息。第一页为page=0，每页显示10本书
// 返回值：BookShortInfoResponse
export async function GetBooksShortInfo(page: number = 0, page_size: number = 10): Promise<BookShortInfoResponse> {
  try {
    const response = await axios.get<BookShortInfoResponse>(urlStore.bookLibraryGetBooksShortInfo, {
      params: {
        page: page,
        page_size: page_size,
      },
    });
    console.log('获取书籍信息：', response.data);
    return response.data;
  } catch (error) {
    console.error('在向服务器请求书籍信息时发生错误：', error);
    return { books: [] };
  }
}

// // 通过已有的book_id, 通过axios库，查询url获取书籍信息, 返回BookInfo
export async function GetBookInfoById(book_id: string): Promise<BookInfo> {
  const response = await axios.get<BookInfoResponse>(urlStore.bookLibraryGetBookInfoById, {
    params: {
      book_id: book_id,
    },
  });
  return response.data.book;
}

export async function GetBookPDFUrl(book_id: string) {
  const pdfFileUrl = ref<string>();
  if (book_id != null) {
    const _pdfFileURL = new URL(urlStore.bookLibraryGetServedBookfileById);
    _pdfFileURL.searchParams.append('book_id', book_id);
    pdfFileUrl.value = _pdfFileURL.toString();
  } else {
    console.error('book_id is null!');
  }

  console.log('欲访问PDF文件的URL为:', pdfFileUrl.value);
  return pdfFileUrl.value;
}

// 获取书籍的第一页PDF文件的URL
export async function GetFirstPagePDFUrl(book_id: string) {
  const pdfFileUrl = ref<string>();
  if (book_id != null) {
    const _pdfFileURL = new URL(urlStore.BookLibraryGetBookFirstPagePdf);
    _pdfFileURL.searchParams.append('book_id', book_id);
    pdfFileUrl.value = _pdfFileURL.toString();
  } else {
    console.error('book_id is null!');
  }

  console.log('欲访问PDF文件的URL为:', pdfFileUrl.value);
  return pdfFileUrl.value;
}
