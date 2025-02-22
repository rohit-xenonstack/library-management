export interface BookData {
  isbn: string;
  title: string;
  authors: string;
  publisher: string;
  version: string;
  avilableCopies?: number;
  totalCopies?: number;
  libraryID?: string;
}
