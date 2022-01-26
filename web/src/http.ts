export interface HttpError extends Error {
  status?: string
  statusCode?: number
  response?: Response
}

export default {
  checkStatus(response: Response) {
    if (response.status >= 200 && response.status < 300) {
      return response
    }

    const error: HttpError = new Error(`HTTP Error ${response.statusText}`)
    error.status = response.statusText
    error.statusCode = response.status
    error.response = response
    throw error
  },
  url(path: string) {
    return import.meta.env.VITE_API_HOST + path
  },
}
