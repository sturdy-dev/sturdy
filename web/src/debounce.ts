export default function debounce<A extends any[]>(
  func: (...args: A) => void,
  wait: number
): (...args: A) => void {
  let timeout: any

  return (...args) => {
    clearTimeout(timeout)
    timeout = setTimeout(() => func(...args), wait)
  }
}
