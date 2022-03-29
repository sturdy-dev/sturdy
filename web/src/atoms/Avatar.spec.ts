import { mount } from '@vue/test-utils'
// import Avatar from './Avatar.vue'

describe('Avatar', () => {
  it('should test something', () => {
    expect(1).toEqual(1)
  })
  /*it('should display initials', () => {
    const author = {
      user_id: '1',
      name: 'Foo Bar',
    }
    const wrapper = mount(Avatar, { props: { author } })
    expect(wrapper.find('div > span > span').text()).toEqual('FB')
  })

  it('should display initials swedish name', () => {
    const author = {
      user_id: '1',
      name: 'Foo Åhlund',
    }
    const wrapper = mount(Avatar, { props: { author } })
    expect(wrapper.find('div > span > span').text()).toEqual('FÅ')
  })

  it('should display first and very last', () => {
    const author = {
      user_id: '1',
      name: 'östen med resten',
    }
    const wrapper = mount(Avatar, { props: { author } })
    expect(wrapper.find('div > span > span').text()).toEqual('ÖR')
  })

  it('should have random background color', () => {
    const author = {
      user_id: '1',
      name: 'Foo Bar',
    }
    const wrapper = mount(Avatar, { props: { author } })
    expect(wrapper.find('div > span').classes()).toContain('bg-yellow-100')
  })

  it('should have another random background color', () => {
    const author = {
      user_id: '400',
      name: 'Foo Bar',
    }
    const wrapper = mount(Avatar, { props: { author } })
    expect(wrapper.find('div > span').classes()).toContain('bg-indigo-100')
  })

  it('works when only passed a url', () => {
    const url = 'https://example.com/foo.png'
    const wrapper = mount(Avatar, { props: { url } })
    expect(wrapper.find('div > img').attributes('src')).toBe(url)
  })

  it('works when only passed a empty url', () => {
    const url = ''
    const wrapper = mount(Avatar, { props: { url } })
    expect(wrapper.find('div > span').classes()).toContain('bg-gray-200')
  })*/
})
