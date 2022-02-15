export interface HttpError extends Error {
  status?: string
  statusCode?: number
  response?: Response
}

const url = (path: string) => {
  if (!path.startsWith('/')) path = '/' + path
  const host = import.meta.env.VITE_API_HOST ? (import.meta.env.VITE_API_HOST as string) : ''
  const prefix = import.meta.env.VITE_API_PATH ? (import.meta.env.VITE_API_PATH as string) : ''
  return `${host}${prefix}${path}`
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
  url,
}
