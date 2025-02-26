export interface SignInRequest {
  email: string
}

export interface RegisterReaderRequest {
  name: string
  email: string
  contact: string
  library_id: string
}

export interface RemoveBookRequest {
  isbn: string
}

export interface ApproveRequest {
  request_id: string
  user_id: string
}

export interface RejectRequest {
  request_id: string
  user_id: string
}

export interface SearchBookRequest {
  search_string: string
  field: 'title' | 'authors' | 'publisher'
}

export interface RemoveBookRequest {
  isbn: string
}

export interface RequestBookRequest {
  isbn: string
  email: string
}
