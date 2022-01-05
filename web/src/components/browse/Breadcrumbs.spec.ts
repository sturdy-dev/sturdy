import { Breadcrumbs, Crumb } from './Breadcrumbs'

describe('Breadcrumbs', () => {
  it('is empty', () => {
    const res = Breadcrumbs.paths('')
    expect(res).toEqual(Array<Crumb>())
  })

  it('to be set', () => {
    const res = Breadcrumbs.paths('this/is/a/path')

    const expected = Array<Crumb>(
      { fullPath: 'this', name: 'this', isCurrent: false },
      { fullPath: 'this/is', name: 'is', isCurrent: false },
      { fullPath: 'this/is/a', name: 'a', isCurrent: false },
      { fullPath: 'this/is/a/path', name: 'path', isCurrent: true }
    )

    expect(res).toEqual(expected)
  })
})
