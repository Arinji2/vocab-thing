export function generateTitle(title?: string) {
  const defaultTitle = 'Vocab Thing'
  return title ? `${title} | ${defaultTitle}` : defaultTitle
}

export function generateDescription(description?: string) {
  const defaultDescription =
    'Save words and phrases you find on the internet, and use them in the future effortlessly.'
  return description
    ? `${description}. ${defaultDescription}`
    : defaultDescription
}
