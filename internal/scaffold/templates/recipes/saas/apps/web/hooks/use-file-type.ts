import { useMemo } from 'react';

export type FileCategory =
  | 'image'
  | 'video'
  | 'audio'
  | 'pdf'
  | 'spreadsheet'
  | 'document'
  | 'archive'
  | 'code'
  | 'other';

export function useFileType(fileNameOrType?: string | null) {
  const category = useMemo((): FileCategory => {
    if (!fileNameOrType) return 'other';

    const str = fileNameOrType.toLowerCase();

    // Check by MIME-like types or common extensions
    if (
      str.match(/\/(jpg|jpeg|png|gif|webp|svg|avif)$/) ||
      str.match(/\.(jpg|jpeg|png|gif|webp|svg|avif)$/)
    ) {
      return 'image';
    }
    if (
      str.match(/\/(mp4|webm|ogg|mov)$/) ||
      str.match(/\.(mp4|webm|ogg|mov)$/)
    ) {
      return 'video';
    }
    if (
      str.match(/\/(mp3|wav|flac|aac)$/) ||
      str.match(/\.(mp3|wav|flac|aac)$/)
    ) {
      return 'audio';
    }
    if (str.includes('pdf') || str.endsWith('.pdf')) {
      return 'pdf';
    }
    if (
      str.match(
        /\/(ms-excel|vnd\.openxmlformats-officedocument\.spreadsheetml)|csv/,
      ) ||
      str.match(/\.(xls|xlsx|csv)$/)
    ) {
      return 'spreadsheet';
    }
    if (
      str.match(
        /\/(msword|vnd\.openxmlformats-officedocument\.wordprocessingml|plain)/,
      ) ||
      str.match(/\.(doc|docx|txt|rtf)$/)
    ) {
      return 'document';
    }
    if (
      str.match(/\/(zip|x-7z-compressed|x-rar-compressed|x-tar|x-gzip)/) ||
      str.match(/\.(zip|7z|rar|tar|gz)$/)
    ) {
      return 'archive';
    }
    if (
      str.match(
        /\.(js|ts|tsx|jsx|html|css|json|py|go|rs|php|rb|cpp|c|cs|java|sh|md)$/,
      )
    ) {
      return 'code';
    }

    return 'other';
  }, [fileNameOrType]);

  return category;
}
