import mentionify from './mentionify'

describe('mentionify', () => {
  it('should not mofidy string without mentions', () => {
    const input = 'Hello world'
    const result = mentionify(input, '@', [], 'text-semibold')
    expect(result).toBe(input)
  })

  it('should highlight single-world mention', () => {
    const input = 'Hello @world-id'
    const result = mentionify(input, '@', [{ id: 'world-id', name: 'world' }], 'text-semibold')
    expect(result).toBe('Hello <span class="text-semibold">@world</span>')
  })

  it('should not use html syntax if style is not set', () => {
    const input = 'Hello @whole world id'
    const result = mentionify(input, '@', [{ id: 'whole world id', name: 'whole world' }])
    expect(result).toBe('Hello @whole world')
  })

  it('should highlight double-world mention', () => {
    const input = 'Hello @whole world id'
    const result = mentionify(
      input,
      '@',
      [{ id: 'whole world id', name: 'whole world' }],
      'text-semibold'
    )
    expect(result).toBe('Hello <span class="text-semibold">@whole world</span>')
  })

  it('should not modify string without users', () => {
    const input = 'Hello @world-id'
    const result = mentionify(input, '@', [], 'text-semibold')
    expect(result).toBe(input)
  })

  it('should not modify other mentions', () => {
    const input = 'Hello #world-id'
    const result = mentionify(input, '@', [{ id: 'world-id', name: 'world' }], 'text-semibold')
    expect(result).toBe(input)
  })

  it('should not detect mention in the middle of the word', () => {
    const input = 'Hello wor@ld-id'
    const result = mentionify(input, '@', [{ id: 'ld-id', name: 'ld' }], 'text-semibold')
    expect(result).toBe(input)
  })

  it('should not detect mention in the middle of the word', () => {
    const input = 'Hello wor@ld-id'
    const result = mentionify(input, '@', [{ id: 'ld-id', name: 'ld' }], 'text-semibold')
    expect(result).toBe(input)
  })

  it('should detect mention in the beginning of the word', () => {
    const input = 'Hello @world-id!'
    const result = mentionify(input, '@', [{ id: 'world-id', name: 'world' }], 'text-semibold')
    expect(result).toBe('Hello <span class="text-semibold">@world</span>!')
  })

  it('should detect mention in the beginning of the word', () => {
    const input = 'Hello @world-id!'
    const result = mentionify(input, '@', [{ id: 'world-id', name: 'world' }], 'text-semibold')
    expect(result).toBe('Hello <span class="text-semibold">@world</span>!')
  })

  it('should not detect mention inside a word', () => {
    const input = 'Hello @world-idish'
    const result = mentionify(input, '@', [{ id: 'world-id', name: 'world' }], 'text-semibold')
    expect(result).toBe('Hello @world-idish')
  })
})
