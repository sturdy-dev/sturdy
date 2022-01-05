export const IsFocusChildOfElementWithClass = (e: FocusEvent, className: string): boolean => {
  let target = <Element>e.relatedTarget

  // If the blur is lost to the dropdown, don't de-select
  for (let i = 0; i < 10; i++) {
    if (!target) {
      break
    }
    if (target.classList.contains(className)) {
      return true
    }
    if (!target.parentElement) {
      return false
    }
    target = target.parentElement
  }

  return false
}
