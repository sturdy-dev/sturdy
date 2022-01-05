export default {
  checkStatus(response: Response) {
    if (response.status >= 200 && response.status < 300) {
      return response
    }

    interface HttpError extends Error {
      status?: string
      response?: Response
    }

    const error: HttpError = new Error(`HTTP Error ${response.statusText}`)
    error.status = response.statusText
    error.response = response
    throw error
  },
  url(path: string) {
    return import.meta.env.VITE_API_HOST + path
  },
}
